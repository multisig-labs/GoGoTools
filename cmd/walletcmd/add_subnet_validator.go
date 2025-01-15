package walletcmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newAddSubnetValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-subnet-validator subnetID nodeID duration weight",
		Short: "Issue a AddSubnetValidatorTx tx and return the txID",
		Long:  ``,
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)
			subnetID, err := ids.FromString(args[0])
			cobra.CheckErr(err)
			nodeShortID, err := ids.NodeIDFromString(args[1])
			cobra.CheckErr(err)
			d, err := strconv.ParseUint(args[2], 10, 64)
			cobra.CheckErr(err)
			weight, err := strconv.ParseUint(args[3], 10, 64)
			cobra.CheckErr(err)

			txID, err := addSubnetValidator(subnetID, nodeShortID, d, weight, key)
			cobra.CheckErr(err)
			fmt.Println(txID)
			return nil
		},
	}
	return cmd
}

func addSubnetValidator(subnetID ids.ID, nodeID ids.NodeID, duration uint64, weight uint64, key *secp256k1.PrivateKey) (ids.ID, error) {
	ctx := context.Background()
	uri := viper.GetString("node-url")
	kc := secp256k1fx.NewKeychain(key)
	startTime := time.Now().Add(5 * time.Second)
	endTime := startTime.Add(time.Duration(duration) * time.Second)

	wallet, err := primary.MakeWallet(ctx, uri, kc, kc, primary.WalletConfig{
		SubnetIDs: []ids.ID{subnetID},
	})
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to initialize wallet: %w", err)
	}

	pWallet := wallet.P()

	tx, err := pWallet.IssueAddSubnetValidatorTx(&txs.SubnetValidator{
		Validator: txs.Validator{
			NodeID: nodeID,
			Start:  uint64(startTime.Unix()),
			End:    uint64(endTime.Unix()),
			Wght:   weight,
		},
		Subnet: subnetID,
	})
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to issue add subnet validator transaction: %s\n", err)
	}

	return tx.TxID, nil
}
