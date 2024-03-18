package walletcmd

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCreateSubnetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-subnet work-dir",
		Short: "Issue a CreateSubnet tx and return the txID (which is the subnetID)",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// if exists := utils.DirExists(args[0]); !exists {
			// 	return fmt.Errorf("node directory does not exist: %s", args[0])
			// }
			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)
			txID, err := createSubnet(key)
			cobra.CheckErr(err)
			fmt.Println(txID)
			return nil
		},
	}
	return cmd
}

func createSubnet(key *secp256k1.PrivateKey) (ids.ID, error) {
	uri := viper.GetString("node-url")
	kc := secp256k1fx.NewKeychain(key)
	subnetOwner := key.Address()
	ctx := context.Background()

	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:          uri,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to initialize wallet: %w", err)
	}

	owner := &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs: []ids.ShortID{
			subnetOwner,
		},
	}

	createSubnetTx, err := wallet.P().IssueCreateSubnetTx(owner)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to issue CreateSubnetTx: %w", err)
	}

	return createSubnetTx.TxID, nil
}
