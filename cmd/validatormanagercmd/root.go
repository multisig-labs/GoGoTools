package validatormanagercmd

import (
	"github.com/multisig-labs/gogotools/pkg/application"
	"github.com/spf13/cobra"
)

var app *application.GoGoTools

func NewCmd(injectedApp *application.GoGoTools) *cobra.Command {
	app = injectedApp

	cmd := &cobra.Command{
		Use:   "vmgr",
		Short: "Tools interacting with a Validator Manager contract",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newInfoCmd())
	return cmd
}
