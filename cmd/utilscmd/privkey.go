package utilscmd

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

func newPrivkeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privkey key",
		Short: "Show various address encodings of a private key",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var pkBytes []byte
			var err error

			if strings.HasPrefix(args[0], "PrivateKey-") {
				// CB58 decoding
				pk, _ := strings.CutPrefix(args[0], "PrivateKey-")
				pkBytes, err = cb58.Decode(pk)
				cobra.CheckErr(err)
			} else {
				// Hex input decoding
				pkBytes, err = hex.DecodeString(strings.TrimPrefix(args[0], "0x"))
				cobra.CheckErr(err)
			}

			pkcb, _ := cb58.Encode(pkBytes)
			fmt.Printf("%-14s %#x\n", "PrivKey Hex:", pkBytes)
			fmt.Printf("%-14s PrivateKey-%s\n", "PrivKey CB58:", pkcb)

			secpk, err := secp256k1.ToPrivateKey(pkBytes)
			cobra.CheckErr(err)

			ethAddr := ethcrypto.PubkeyToAddress(secpk.ToECDSA().PublicKey).String()
			fmt.Printf("%-14s %s\n", "Eth addr:", ethAddr)

			addr, err := address.Format("P", "avax", secpk.PublicKey().Address().Bytes())
			cobra.CheckErr(err)

			addrs, err := addrVariants(addr)
			cobra.CheckErr(err)

			fmt.Println(strings.Join(addrs, "\n"))
		},
	}
	return cmd
}
