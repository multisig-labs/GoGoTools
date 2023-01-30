package nodecmd

import (
	"github.com/multisig-labs/gogotools/pkg/application"
	"github.com/spf13/cobra"
)

var app *application.GoGoTools

func NewCmd(injectedApp *application.GoGoTools) *cobra.Command {
	app = injectedApp

	cmd := &cobra.Command{
		Use:   "node",
		Short: "Create and run a single-node avalanchego",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newCreateUserCmd())
	cmd.AddCommand(newHealthCmd())
	cmd.AddCommand(newExplorerCmd())
	cmd.AddCommand(newInfoCmd())
	cmd.AddCommand(newLoadVMsCmd())
	cmd.AddCommand(newLogLevelCmd())
	cmd.AddCommand(newPrepareCmd())
	cmd.AddCommand(newRunCmd())
	return cmd
}
