package utilscmd

import (
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
)

func newPrivkeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privkey key",
		Short: "Show address of a private key",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var pkBytes []byte
			var err error

			if strings.HasPrefix(args[0], "0x") {
				// Decode hex input
				pkBytes, err = hexutil.Decode(args[0])
				cobra.CheckErr(err)
			} else {
				// Existing CB58 decoding
				pk, _ := strings.CutPrefix(args[0], "PrivateKey-")
				pkBytes, err = cb58.Decode(pk)
				cobra.CheckErr(err)
			}

			pkcb, _ := cb58.Encode(pkBytes)
			fmt.Printf("PrivKey Hex: %#x\n", pkBytes)
			fmt.Printf("PrivKey CB58: PrivateKey-%s\n", pkcb)

			secpk, err := secp256k1.ToPrivateKey(pkBytes)
			cobra.CheckErr(err)

			addr, err := address.Format("P", "avax", secpk.PublicKey().Address().Bytes())
			cobra.CheckErr(err)

			addrs, err := addrVariants(addr)
			cobra.CheckErr(err)

			fmt.Println(strings.Join(addrs, "\n"))
		},
	}
	return cmd
}
