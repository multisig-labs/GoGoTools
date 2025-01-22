package walletcmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
)

func newSignL1WeightMsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign-l1-weight-msg network blockchainID validationID weight nonce blsPK index",
		Short: "Sign a L1 weight message.",
		Args:  cobra.MinimumNArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			networkID, err := constants.NetworkID(args[0])
			cobra.CheckErr(err)
			blockchainID, err := ids.FromString(args[1])
			cobra.CheckErr(err)
			validationID, err := ids.FromString(args[2])
			cobra.CheckErr(err)
			weight, err := strconv.Atoi(args[3])
			cobra.CheckErr(err)
			nonce, err := strconv.Atoi(args[4])
			cobra.CheckErr(err)
			blsSKBytes, err := hexutil.Decode(args[5])
			cobra.CheckErr(err)
			sk, err := bls.SecretKeyFromBytes(blsSKBytes)
			cobra.CheckErr(err)
			index, err := strconv.Atoi(args[6])
			cobra.CheckErr(err)
			signedWarpMsg, err := signL1WeightMsg(networkID, blockchainID, validationID, uint64(weight), uint64(nonce), sk, index)
			cobra.CheckErr(err)
			fmt.Println(hexutil.Encode(signedWarpMsg.Bytes()))
			return nil
		},
	}

	return cmd
}

func signL1WeightMsg(networkID uint32, blockchainID, validationID ids.ID, weight, nonce uint64, sk *bls.LocalSigner, index int) (*warp.Message, error) {
	addressedCallPayload, err := message.NewL1ValidatorWeight(
		validationID,
		nonce,
		weight,
	)
	if err != nil {
		log.Fatalf("failed to create L1ValidatorWeight message: %s\n", err)
	}
	// addressedCallPayloadJSON, err := json.MarshalIndent(addressedCallPayload, "", "\t")
	// if err != nil {
	// 	log.Fatalf("failed to marshal L1ValidatorWeight message: %s\n", err)
	// }
	// log.Println(string(addressedCallPayloadJSON))

	addressedCall, err := payload.NewAddressedCall(
		[]byte{},
		addressedCallPayload.Bytes(),
	)
	if err != nil {
		log.Fatalf("failed to create AddressedCall message: %s\n", err)
	}

	unsignedWarp, err := warp.NewUnsignedMessage(
		networkID,
		blockchainID,
		addressedCall.Bytes(),
	)
	if err != nil {
		log.Fatalf("failed to create unsigned Warp message: %s\n", err)
	}
	fmt.Printf("Unsigned: %s\n", hexutil.Encode(unsignedWarp.Bytes()))

	signedWarpMsg, err := signWarpMsg(unsignedWarp, sk, index)
	if err != nil {
		log.Fatalf("failed to sign Warp message: %s\n", err)
	}

	return signedWarpMsg, nil
}
