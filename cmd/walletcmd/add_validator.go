package walletcmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO This doesnt seem to work on a --staking-disabled node, so maybe not worth including?

func newAddValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-validator nodeID duration",
		Short: "Issue a AddValidator tx and return the txID",
		Long:  ``,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)
			d, err := strconv.ParseUint(args[1], 10, 64)
			cobra.CheckErr(err)
			txID, err := addValidator(args[0], d, key)
			cobra.CheckErr(err)
			fmt.Println(txID)
			return nil
		},
	}
	return cmd
}

func addValidator(nodeID string, duration uint64, key *secp256k1.PrivateKey) (ids.ID, error) {
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

	nodeShortID, err := ids.NodeIDFromString(nodeID)
	if err != nil {
		return ids.ID{}, fmt.Errorf("error decoding nodeID %s: %w", nodeID, err)
	}

	startTime := time.Now().Add(5 * time.Second)
	endTime := startTime.Add(time.Duration(duration) * time.Second)

	vdr := txs.Validator{
		NodeID: nodeShortID,
		Start:  uint64(startTime.Unix()),
		End:    uint64(endTime.Unix()),
		Wght:   2 * units.KiloAvax,
	}

	tx, err := wallet.P().IssueAddValidatorTx(&vdr, owner, 20000)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to issue AddValidatorTx: %w", err)
	}

	return tx.TxID, nil
}
