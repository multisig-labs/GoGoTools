package walletcmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO This doesnt seem to work on a --staking-disabled node, so maybe not worth including?

func newAddValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-validator work-dir",
		Short: "Issue a AddValidator tx and return the txID",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if exists := utils.DirExists(args[0]); !exists {
				return fmt.Errorf("node directory does not exist: %s", args[0])
			}
			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)
			txID, err := addValidator(key)
			cobra.CheckErr(err)
			fmt.Println(txID)
			return nil
		},
	}
	return cmd
}

func addValidator(key *secp256k1.PrivateKey) (ids.ID, error) {
	uri := viper.GetString("node-url")
	kc := secp256k1fx.NewKeychain(key)
	subnetOwner := key.Address()
	ctx := context.Background()

	wallet, err := primary.NewWalletFromURI(ctx, uri, kc)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to initialize wallet: %w", err)
	}

	owner := &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs: []ids.ShortID{
			subnetOwner,
		},
	}

	nodeID := "NodeID-5FKRwqFyQnZGoN7FTc3t3TUHyuUGPhBxJ"
	nodeShortID, err := ids.NodeIDFromString(nodeID)
	if err != nil {
		return ids.ID{}, fmt.Errorf("error decoding nodeID %s: %w", nodeID, err)
	}

	startTime := time.Now().Add(10 * time.Second)
	endTime := startTime.Add(24 * time.Hour)

	vdr := txs.Validator{
		NodeID: nodeShortID,
		Start:  uint64(startTime.Unix()),
		End:    uint64(endTime.Unix()),
		Wght:   2 * units.KiloAvax,
	}

	txID, err := wallet.P().IssueAddValidatorTx(&vdr, owner, 20000)
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to issue AddValidatorTx: %w", err)
	}

	return txID, nil
}
