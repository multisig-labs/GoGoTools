package castcmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/application"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var app *application.GoGoTools

func NewCmd(injectedApp *application.GoGoTools) *cobra.Command {
	app = injectedApp

	cmd := &cobra.Command{
		Use:   "cast",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmd.PersistentFlags().String("accounts", "accounts.json", "JSON of actors")
	viper.BindPFlag("accounts", cmd.PersistentFlags().Lookup("accounts"))

	cmd.PersistentFlags().String("contracts", "contracts.json", "JSON of contract addresses")
	viper.BindPFlag("contracts", cmd.PersistentFlags().Lookup("contracts"))

	cmd.AddCommand(newBalancesCmd())
	cmd.AddCommand(newCallCmd())
	cmd.AddCommand(newSendCmd())
	cmd.AddCommand(newSendEthCmd())
	return cmd
}
