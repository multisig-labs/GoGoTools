package nodecmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/application"
	"github.com/spf13/cobra"
)

var app *application.GoGoTools

func NewCmd(injectedApp *application.GoGoTools) *cobra.Command {
	app = injectedApp

	cmd := &cobra.Command{
		Use:   "node",
		Short: "Create and run avalanchego configured as a single node with staking disabled",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmd.AddCommand(newCreateUserCmd())
	cmd.AddCommand(newInfoCmd())
	cmd.AddCommand(newLoadVMsCmd())
	cmd.AddCommand(newLogLevelCmd())
	cmd.AddCommand(newPrepareCmd())
	cmd.AddCommand(newRunCmd())
	return cmd
}
