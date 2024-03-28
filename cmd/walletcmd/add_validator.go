package walletcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/platformvm/reward"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newAddValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-validator nodeID blsKey blsSig duration",
		Short: "Issue a AddValidator tx and return the txID",
		Long:  ``,
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)
			d, err := strconv.ParseUint(args[3], 10, 64)
			cobra.CheckErr(err)
			txID, err := addValidator(args[0], args[1], args[2], d, key)
			cobra.CheckErr(err)
			fmt.Println(txID)
			return nil
		},
	}
	return cmd
}

func addValidator(nodeID string, blsKey string, blsSig string, duration uint64, key *secp256k1.PrivateKey) (ids.ID, error) {
	ctx := context.Background()
	uri := viper.GetString("node-url")
	kc := secp256k1fx.NewKeychain(key)
	startTime := time.Now().Add(5 * time.Second)
	endTime := startTime.Add(time.Duration(duration) * time.Second)
	validatorRewardAddr := key.Address()
	delegatorRewardAddr := key.Address()
	delegationFee := uint32(reward.PercentDenominator / 2) // 50%

	nodeShortID, err := ids.NodeIDFromString(nodeID)
	if err != nil {
		return ids.ID{}, fmt.Errorf("error decoding nodeID %s: %w", nodeID, err)
	}

	nodePOP, err := makePoP(blsKey, blsSig)
	if err != nil {
		return ids.ID{}, fmt.Errorf("error decoding nodeID %s: %w", nodeID, err)
	}

	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:          uri,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		return ids.Empty, fmt.Errorf("failed to initialize wallet: %w", err)
	}
	pWallet := wallet.P()
	avaxAssetID := pWallet.AVAXAssetID()

	tx, err := pWallet.IssueAddPermissionlessValidatorTx(
		&txs.SubnetValidator{Validator: txs.Validator{
			NodeID: nodeShortID,
			Start:  uint64(startTime.Unix()),
			End:    uint64(endTime.Unix()),
			Wght:   2_000 * units.Avax, // TODO fix this
		}},
		nodePOP,
		avaxAssetID,
		&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{validatorRewardAddr},
		},
		&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{delegatorRewardAddr},
		},
		delegationFee,
	)
	if err != nil {
		return ids.Empty, err
	}

	return tx.TxID, nil
}

func makePoP(blsPubkey string, blsSig string) (*signer.ProofOfPossession, error) {
	pop := &signer.ProofOfPossession{}
	popjs := fmt.Sprintf(`{"publicKey":"%s","proofOfPossession":"%s"}`, blsPubkey, blsSig)
	if err := json.Unmarshal([]byte(popjs), pop); err != nil {
		return nil, err
	}

	if err := pop.Verify(); err != nil {
		return nil, err
	}

	return pop, nil
}
