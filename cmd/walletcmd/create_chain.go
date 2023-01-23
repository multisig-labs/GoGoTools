package walletcmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCreateChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-chain name vm genesisFile [subnetID]",
		Short: "Issue a CreateBlockchain tx and return the txID. Also creates a new Subnet if subnetID is not specified.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)

			name := args[0]

			paddedBytes := [32]byte{}
			copy(paddedBytes[:], []byte(args[1]))
			vmID, err := ids.ToID(paddedBytes[:])
			cobra.CheckErr(err)

			genesisBytes, err := os.ReadFile(args[2])
			cobra.CheckErr(err)

			var subnetID ids.ID
			if len(args) < 4 || args[3] == "" {
				app.Log.Debug("No SubnetID supplied, creating a new Subnet")
				subnetID, err = createSubnet(key)
				cobra.CheckErr(err)
			} else {
				subnetID, err = ids.FromString(args[3])
				cobra.CheckErr(err)
			}

			txID, err := createChain(key, subnetID, name, vmID, genesisBytes)
			cobra.CheckErr(err)
			fmt.Printf("%s/ext/bc/%s/rpc\n", primary.LocalAPIURI, txID)
			return nil
		},
	}
	return cmd
}

func createChain(key *crypto.PrivateKeySECP256K1R, subnetID ids.ID, name string, vmID ids.ID, genesisBytes []byte) (ids.ID, error) {
	uri := viper.GetString("node-url")
	kc := secp256k1fx.NewKeychain(key)
	ctx := context.Background()

	wallet, err := primary.NewWalletWithTxs(ctx, uri, kc, subnetID)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to initialize wallet: %w", err)
	}

	createChainTxID, err := wallet.P().IssueCreateChainTx(
		subnetID,
		genesisBytes,
		vmID,
		nil,
		name,
	)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to issue CreateBlockchainTx: %w", err)
	}

	app.Log.Info("created new blockchain ", createChainTxID)
	app.Log.Info("NOTE: Check the data/logs/main.log file, as the blockchain may not start if anything is wrong with the VM binary or paths")
	return createChainTxID, nil
}
