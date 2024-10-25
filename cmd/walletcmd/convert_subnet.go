package walletcmd

import (
	"context"
	"fmt"
	"log"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newConvertSubnetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-subnet subnetID mgrChainID mgrAddress",
		Short: "Issue a ConvertSubnet tx and return the txID.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			key, err := decodePrivateKey(viper.GetString("pk"))
			cobra.CheckErr(err)
			subnetID, err := ids.FromString(args[0])
			cobra.CheckErr(err)
			mgrChainID, err := ids.FromString(args[1])
			cobra.CheckErr(err)

			if !common.IsHexAddress(args[2]) {
				return fmt.Errorf("invalid manager address format: %s", args[2])
			}
			mgrAddress := common.HexToAddress(args[2])

			// For now, grab nodeID for the current node
			ctx := context.Background()
			uri := viper.GetString("node-url")
			infoClient := info.NewClient(uri)
			nodeID, nodePoP, err := infoClient.GetNodeID(ctx)
			cobra.CheckErr(err)

			txID, err := convertSubnet(subnetID, mgrChainID, mgrAddress, key, nodeID, nodePoP)
			cobra.CheckErr(err)
			fmt.Println(txID)
			return nil
		},
	}

	return cmd
}

func convertSubnet(subnetID ids.ID, mgrChainID ids.ID, mgrAddress common.Address, key *secp256k1.PrivateKey, nodeID ids.NodeID, nodePoP *signer.ProofOfPossession) (*txs.Tx, error) {
	weight := uint64(1000000)
	uri := viper.GetString("node-url")
	kc := secp256k1fx.NewKeychain(key)

	validationID := subnetID.Append(0)
	conversionID, err := message.SubnetConversionID(message.SubnetConversionData{
		SubnetID:       subnetID,
		ManagerChainID: mgrChainID,
		ManagerAddress: mgrAddress.Bytes(),
		Validators: []message.SubnetConversionValidatorData{
			{
				NodeID:       nodeID.Bytes(),
				BLSPublicKey: nodePoP.PublicKey,
				Weight:       weight,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:          uri,
		AVAXKeychain: kc,
		EthKeychain:  kc,
		SubnetIDs:    []ids.ID{subnetID},
	})
	if err != nil {
		return nil, err
	}

	pWallet := wallet.P()

	owner := message.PChainOwner{
		Threshold: 1,
		Addresses: kc.Addresses().List(),
	}

	tx, err := pWallet.IssueConvertSubnetTx(
		subnetID,
		mgrChainID,
		mgrAddress.Bytes(),
		[]*txs.ConvertSubnetValidator{
			{
				NodeID:                nodeID.Bytes(),
				Weight:                weight,
				Balance:               units.Avax,
				Signer:                *nodePoP,
				RemainingBalanceOwner: owner,
				DeactivationOwner:     owner,
			},
		},
	)

	log.Printf("converted subnet %s with transactionID %s, validationID %s, and conversionID %s\n",
		subnetID,
		tx.ID(),
		validationID,
		conversionID,
	)

	return tx, nil
}
