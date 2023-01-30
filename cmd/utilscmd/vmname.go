package utilscmd

import (
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/spf13/cobra"
)

func newVMNameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vmname [ID]",
		Short: "Given a vmID, try to decode the ASCII name",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := ids.FromString(args[0])
			if err != nil {
				return err
			}
			out := strings.Builder{}
			for _, v := range id {
				out.Write([]byte{v})
			}
			fmt.Println(out.String())
			return nil
		},
	}
	return cmd
}
