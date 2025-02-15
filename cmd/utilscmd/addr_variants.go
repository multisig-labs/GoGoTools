package utilscmd

import (
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/spf13/cobra"
)

func newAddrVariantsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addr-variants addr",
		Short: "Show address variants",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addrs, err := addrVariants(args[0])
			if err != nil {
				return err
			}
			fmt.Println(strings.Join(addrs, "\n"))

			return nil
		},
	}
	return cmd
}

func addrVariants(addr string) ([]string, error) {
	hrps := []string{"avax", "fuji", "local", "custom"}

	id, err := address.ParseToID(addr)
	if err != nil {
		return nil, err
	}

	out := []string{fmt.Sprintf("Raw Bytes of %s: %0xs", strings.TrimPrefix(addr, "P-"), id.Hex())}

	for _, hrp := range hrps {
		a, err := address.Format("P", hrp, id.Bytes())
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}

	return out, nil
}
