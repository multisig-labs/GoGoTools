package utilscmd

import (
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/spf13/cobra"
)

func newPrivkeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privkey key",
		Short: "Show address of a private key",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pk, _ := strings.CutPrefix(args[0], "PrivateKey-")
			pkBytes, err := cb58.Decode(pk)
			cobra.CheckErr(err)
			fmt.Printf("PrivKey Bytes: %#x\n", pkBytes)

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
