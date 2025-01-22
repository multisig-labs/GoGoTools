package walletcmd

import (
	"context"
	"log"
	"time"

	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newSetL1WeightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-l1-weight warpMsg",
		Short: "Issue a SetL1Weight tx and return the txID.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			warpMsgBytes, err := hexutil.Decode(args[0])
			cobra.CheckErr(err)

			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)
			uri := viper.GetString("node-url")
			kc := secp256k1fx.NewKeychain(key)

			setL1Weight(warpMsgBytes, kc, uri)

			return nil
		},
	}

	return cmd
}

func setL1Weight(warpMsgBytes []byte, kc *secp256k1fx.Keychain, uri string) {
	ctx := context.Background()
	// MakePWallet fetches the available UTXOs owned by [kc] on the P-chain that
	// [uri] is hosting.
	walletSyncStartTime := time.Now()
	wallet, err := primary.MakePWallet(
		ctx,
		uri,
		kc,
		primary.WalletConfig{},
	)
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	setWeightStartTime := time.Now()
	setWeightTx, err := wallet.IssueSetL1ValidatorWeightTx(
		warpMsgBytes,
	)
	if err != nil {
		log.Fatalf("failed to issue set L1 validator weight transaction: %s\n", err)
	}
	log.Printf("issued set weight txID %s in %s\n", setWeightTx.ID(), time.Since(setWeightStartTime))
}
