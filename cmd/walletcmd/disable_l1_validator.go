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

func newDisableL1ValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-l1-validator validationID",
		Short: "Issue a DisableL1Validator tx and return the txID.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			validationID, err := ids.FromString(args[0])
			cobra.CheckErr(err)

			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)
			uri := viper.GetString("node-url")
			kc := secp256k1fx.NewKeychain(key)

			disableL1Validator(validationID, kc, uri)

			return nil
		},
	}

	return cmd
}

func disableL1Validator(validationID ids.ID, kc *secp256k1fx.Keychain, uri string) {
	ctx := context.Background()

	// MakePWallet fetches the available UTXOs owned by [kc] on the P-chain that
	// [uri] is hosting and registers [validationID].
	walletSyncStartTime := time.Now()
	wallet, err := primary.MakePWallet(
		ctx,
		uri,
		kc,
		primary.WalletConfig{
			ValidationIDs: []ids.ID{validationID},
		},
	)
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	disableL1ValidatorStartTime := time.Now()
	disableL1ValidatorTx, err := wallet.IssueDisableL1ValidatorTx(
		validationID,
	)
	if err != nil {
		log.Fatalf("failed to issue disable L1 validator transaction: %s\n", err)
	}
	log.Printf("disabled %s with %s in %s\n",
		validationID,
		disableL1ValidatorTx.ID(),
		time.Since(disableL1ValidatorStartTime),
	)
}
