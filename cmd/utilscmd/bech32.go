package utilscmd

import (
	"fmt"

	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/spf13/cobra"
)

func newBech32Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bech32decode [value]",
		Short: "Decode a Bech32 encoded string",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, _, addrBytes, err := address.Parse(args[0])
			cobra.CheckErr(err)
			fmt.Printf("0x%x\n", addrBytes)
		},
	}
	return cmd
}
