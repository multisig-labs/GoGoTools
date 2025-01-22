package walletcmd

import (
	"fmt"
	"strconv"

	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
)

func newSignWarpMsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign-warp-msg unsignedWarpMsg blsPK index",
		Short: "Sign a warp message.",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			unsignedWarpMsgBytes, err := hexutil.Decode(args[0])
			cobra.CheckErr(err)
			unsignedWarpMsg, err := warp.ParseUnsignedMessage(unsignedWarpMsgBytes)
			cobra.CheckErr(err)
			blsSKBytes, err := hexutil.Decode(args[1])
			cobra.CheckErr(err)
			sk, err := bls.SecretKeyFromBytes(blsSKBytes)
			cobra.CheckErr(err)
			index, err := strconv.Atoi(args[2])
			cobra.CheckErr(err)

			signedWarpMsg, err := signWarpMsg(unsignedWarpMsg, sk, index)
			cobra.CheckErr(err)

			fmt.Println(hexutil.Encode(signedWarpMsg.Bytes()))
			return nil
		},
	}

	return cmd
}

func signWarpMsg(unsignedWarpMsg *warp.UnsignedMessage, sk *bls.LocalSigner, index int) (*warp.Message, error) {
	// signedWarpMsg, err := warp.NewMessage(
	// 	unsignedWarpMsg,
	// 	&warp.BitSetSignature{
	// 		Signers: set.NewBits(index).Bytes(),
	// 		Signature: ([bls.SignatureLen]byte)(
	// 			bls.SignatureToBytes(
	// 				sk.Sign(unsignedWarpMsg.Bytes()),
	// 			),
	// 		),
	// 	},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	sig := ([bls.SignatureLen]byte)(bls.SignatureToBytes(sk.Sign(unsignedWarpMsg.Bytes())))
	sigMap := make(map[int][bls.SignatureLen]byte, index)
	for i := 0; i <= index; i++ {
		sigMap[i] = sig
	}
	aggSig, bitset, err := aggregateSignatures(sigMap)

	signedWarpMsg, err := warp.NewMessage(
		unsignedWarpMsg,
		&warp.BitSetSignature{
			Signers:   bitset.Bytes(),
			Signature: ([bls.SignatureLen]byte)(bls.SignatureToBytes(aggSig)),
		},
	)
	if err != nil {
		return nil, err
	}

	return signedWarpMsg, nil
}

func aggregateSignatures(
	signatureMap map[int][bls.SignatureLen]byte,
) (*bls.Signature, set.Bits, error) {
	// Aggregate the signatures
	signatures := make([]*bls.Signature, 0, len(signatureMap))
	vdrBitSet := set.NewBits()

	for i, sigBytes := range signatureMap {
		sig, err := bls.SignatureFromBytes(sigBytes[:])
		if err != nil {
			msg := "Failed to unmarshal signature"
			return nil, set.Bits{}, fmt.Errorf("%s: %w", msg, err)
		}
		signatures = append(signatures, sig)
		vdrBitSet.Add(i)
	}

	aggSig, err := bls.AggregateSignatures(signatures)
	if err != nil {
		msg := "Failed to aggregate signatures"
		return nil, set.Bits{}, fmt.Errorf("%s: %w", msg, err)
	}
	return aggSig, vdrBitSet, nil
}
