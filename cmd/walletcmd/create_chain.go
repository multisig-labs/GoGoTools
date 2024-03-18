package walletcmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func newCreateChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-chain work-dir name vm [subnetID]",
		Short: "Issue a CreateBlockchain tx and return the txID. Creates a new Subnet if subnetID is not specified.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			workDir := args[0]
			name := args[1]
			vm := args[2]

			var subnetID ids.ID
			if len(args) > 3 {
				subnetID, err = ids.FromString(args[3])
				cobra.CheckErr(err)
			}

			// Dont allow duplicate chain names, for simplicity
			uri := viper.GetString("node-url")
			urlP := fmt.Sprintf("%s/ext/bc/P", uri)
			getBlockchains, err := utils.FetchRPCGJSON(urlP, "platform.getBlockchains", "")
			cobra.CheckErr(err)
			for _, obj := range getBlockchains.Get("result.blockchains").Array() {
				if obj.Get("name").String() == name {
					return fmt.Errorf("blockchain %s already exists, aborting", name)
				}
			}

			_ = viper.BindPFlags(cmd.Flags())

			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)

			// Construct vm ids.ID
			paddedBytes := [32]byte{}
			copy(paddedBytes[:], []byte(vm))
			vmID, err := ids.ToID(paddedBytes[:])
			cobra.CheckErr(err)

			genesisBytes, err := os.ReadFile(viper.GetString("genesis-file"))
			cobra.CheckErr(err)

			if subnetID == ids.Empty {
				app.Log.Info("No SubnetID supplied, creating...")
				subnetID, err = createSubnet(key)
				cobra.CheckErr(err)
				app.Log.Infof("SubnetID %s created", subnetID)
			}

			txID, err := createChain(key, subnetID, name, vmID, genesisBytes)
			cobra.CheckErr(err)
			app.Log.Infof("Chain created with txID: %s", txID)

			// Copy the chain config to the right place
			chainConfigDir := filepath.Join(workDir, "configs", "chains", txID.String())
			err = os.MkdirAll(chainConfigDir, os.ModePerm)
			cobra.CheckErr(err)
			err = utils.CopyFile(viper.GetString("config-file"), filepath.Join(chainConfigDir, "config.json"))
			cobra.CheckErr(err)

			// TODO Creating the chain alias would be nice IF I could get it to work for the RPC url too
			// like instead of http://localhost:9650/ext/bc/ByeHH...yL9/rpc I want http://localhost:9650/ext/bc/MyChainAlias/rpc
			// but I cant get it to work, not even sure if it is capable of working like that.
			//
			// Create an alias in aliases.json
			fileLocations := utils.NewFileLocations(workDir)
			aliasesContent, err := os.ReadFile(fileLocations.ChainAliasesFile)
			cobra.CheckErr(err)
			var aliasesJson string
			aliasesJson = gjson.Parse(string(aliasesContent)).String()
			if aliasesJson == "" {
				app.Log.Warnf("Chain alias not created, unable to parse %s", fileLocations.ChainAliasesFile)
			} else {
				aliasesJson, _ = sjson.Set(aliasesJson, txID.String(), []string{name})
				_ = utils.WriteFileBytes(fileLocations.ChainAliasesFile, []byte(aliasesJson))
			}

			app.Log.Infof("created new blockchain %s with ID: %s", name, txID)
			app.Log.Info("NOTE: Check the data/logs/main.log file, as the blockchain may not start if anything is wrong with the VM binary or paths")
			app.Log.Info("")
			app.Log.Infof("RPC: %s/ext/bc/%s/rpc\n", primary.LocalAPIURI, txID)
			app.Log.Info("")
			app.Log.Info("run 'ggt node info' to see more")

			// Chain config doesnt get picked up until a restart happens.
			// Update: Not sure this is true
			// if exists := utils.FileExists(".pid"); !exists {
			// 	app.Log.Info("Can't find .pid file in current directory, unable to restart node. Stop and restart it to pick up changes.")
			// } else {
			// 	pidContents, err := os.ReadFile(".pid")
			// 	cobra.CheckErr(err)
			// 	pid, err := strconv.Atoi(strings.TrimSpace(string(pidContents)))
			// 	cobra.CheckErr(err)
			// 	err = syscall.Kill(pid, syscall.SIGUSR1)
			// 	cobra.CheckErr(err)
			// 	app.Log.Infof("Sent USR1 to pid %d to restart node", pid)
			// }
			return nil
		},
	}
	cmd.Flags().String("genesis-file", "subnetevm-genesis.json", "Full path to genesis file (Defaults to subnetEVM)")
	cmd.Flags().String("config-file", "subnetevm-config.json", "Full path to chain config file (Defaults to subnetEVM)")
	return cmd
}

func createChain(key *secp256k1.PrivateKey, subnetID ids.ID, name string, vmID ids.ID, genesisBytes []byte) (ids.ID, error) {
	uri := viper.GetString("node-url")
	kc := secp256k1fx.NewKeychain(key)
	ctx := context.Background()

	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:              uri,
		AVAXKeychain:     kc,
		EthKeychain:      kc,
		PChainTxsToFetch: set.Of(subnetID),
	})
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to initialize wallet: %w", err)
	}

	createChainTx, err := wallet.P().IssueCreateChainTx(
		subnetID,
		genesisBytes,
		vmID,
		nil,
		name,
	)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to issue CreateBlockchainTx: %w", err)
	}

	return createChainTx.TxID, nil
}
