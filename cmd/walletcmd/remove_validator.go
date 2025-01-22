package walletcmd

import (
	"context"
	"log"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRemoveValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-validator subnetID nodeID",
		Short: "Issue a RemoveSubnetValidator tx and return the txID.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)
			subnetID, err := ids.FromString(args[0])
			cobra.CheckErr(err)
			nodeID, err := ids.NodeIDFromString(args[1])
			cobra.CheckErr(err)

			uri := viper.GetString("node-url")
			kc := secp256k1fx.NewKeychain(key)

			removeValidator(subnetID, nodeID, kc, uri)

			return nil
		},
	}

	return cmd
}

func removeValidator(subnetID ids.ID, nodeID ids.NodeID, kc *secp256k1fx.Keychain, uri string) {
	ctx := context.Background()

	// MakePWallet fetches the available UTXOs owned by [kc] on the P-chain that
	// [uri] is hosting and registers [subnetID].
	walletSyncStartTime := time.Now()
	wallet, err := primary.MakePWallet(
		ctx,
		uri,
		kc,
		primary.WalletConfig{
			SubnetIDs: []ids.ID{subnetID},
		},
	)
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	removeValidatorStartTime := time.Now()
	removeValidatorTx, err := wallet.IssueRemoveSubnetValidatorTx(
		nodeID,
		subnetID,
	)
	if err != nil {
		log.Fatalf("failed to issue remove subnet validator transaction: %s\n", err)
	}
	log.Printf("removed subnet validator %s from %s with %s in %s\n", nodeID, subnetID, removeValidatorTx.ID(), time.Since(removeValidatorStartTime))
}
