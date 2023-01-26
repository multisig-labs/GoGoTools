package nodecmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Get the health info for a node",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getHealth()
		},
	}
	return cmd
}

// It's not you, Types, it's me. I think we need a break for a bit.
func getHealth() error {
	uri := viper.GetString("node-url")
	urlHealth := fmt.Sprintf("%s/ext/health", uri)

	getHealth, err := utils.FetchRPCGJSON(urlHealth, "health.health", "")
	cobra.CheckErr(err)

	fmt.Println(getHealth.String())
	return nil
}
