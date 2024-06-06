package nodecmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRpcCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rpc [endpoint] [method] [params]",
		Short: "Ergonomic access to node RPC APIs",
		Long:  ``,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			uri := viper.GetString("node-url")
			url := fmt.Sprintf("%s%s", uri, args[0])
			p := ""
			if len(args) > 2 {
				p = args[2]
			}
			result, err := utils.FetchRPCGJSON(url, args[1], p)
			cobra.CheckErr(err)
			fmt.Println(result.String())
		},
	}
	return cmd
}
