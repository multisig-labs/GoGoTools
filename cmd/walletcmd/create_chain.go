package walletcmd

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCreateChainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-chain [subnetID] [name] [vm] [genesis]",
		Short: "Issue a CreateBlockchain tx and return the txID",
		Long:  ``,
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)

			subnetID, err := ids.FromString(args[0])
			cobra.CheckErr(err)

			name := args[1]

			paddedBytes := [32]byte{}
			copy(paddedBytes[:], []byte(args[2]))
			vmID, err := ids.ToID(paddedBytes[:])
			cobra.CheckErr(err)

			genesisBytes, err := hex.DecodeString(args[3])
			cobra.CheckErr(err)

			txID, err := createChain(key, subnetID, name, vmID, genesisBytes)
			cobra.CheckErr(err)
			fmt.Println(txID)
			return nil
		},
	}
	return cmd
}

func createChain(key *crypto.PrivateKeySECP256K1R, subnetID ids.ID, name string, vmID ids.ID, genesisBytes []byte) (ids.ID, error) {
	uri := primary.LocalAPIURI
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

	app.Log.Info("created new blockchain", createChainTxID)
	return createChainTxID, nil
}
