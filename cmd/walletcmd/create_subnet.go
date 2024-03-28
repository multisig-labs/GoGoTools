package walletcmd

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCreateSubnetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-subnet [owner]",
		Short: "Issue a CreateSubnet tx (optionally owned by owner) and return the txID (which is the subnetID)",
		Long:  ``,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := decodePrivateKey(viper.GetString("pk"))
			owner := key.Address()
			if len(args) > 0 {
				owner, err = address.ParseToID(args[0])
				cobra.CheckErr(err)
			}

			cobra.CheckErr(err)
			txID, err := createSubnet(key, owner)
			cobra.CheckErr(err)
			fmt.Println(txID)
			return nil
		},
	}
	return cmd
}

func createSubnet(key *secp256k1.PrivateKey, subnetOwner ids.ShortID) (ids.ID, error) {
	uri := viper.GetString("node-url")
	kc := secp256k1fx.NewKeychain(key)
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
