package utilscmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/application"
	"github.com/spf13/cobra"
)

var app *application.GoGoTools

func NewCmd(injectedApp *application.GoGoTools) *cobra.Command {
	app = injectedApp

	cmd := &cobra.Command{
		Use:          "utils",
		Short:        "Misc utilities",
		Long:         ``,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newMsgDigestCmd())
	cmd.AddCommand(newVMIDCmd())
	cmd.AddCommand(newVMNameCmd())

	return cmd
}
