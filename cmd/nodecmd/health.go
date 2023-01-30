package nodecmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func newHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Get the health info for a node",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			result, err := getHealth()
			cobra.CheckErr(err)
			fmt.Println(result.String())
		},
	}
	return cmd
}

func getHealth() (*gjson.Result, error) {
	uri := viper.GetString("node-url")
	urlHealth := fmt.Sprintf("%s/ext/health", uri)

	return utils.FetchRPCGJSON(urlHealth, "health.health", "")
}
