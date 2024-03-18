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
	cmd.AddCommand(newAddrVariantsCmd())
	cmd.AddCommand(newMnemonicCmd())
	cmd.AddCommand(newMnemonicKeysCmd())
	cmd.AddCommand(newMnemonicAddrsCmd())
	cmd.AddCommand(newPortFwdCmd())
	cmd.AddCommand(newCB58DecodeCmd())
	cmd.AddCommand(newCB58DecodeSigCmd())
	// cmd.AddCommand(newInspectBlockCmd())
	cmd.AddCommand(newRandomCmd())

	return cmd
}
