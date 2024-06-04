package castcmd

import (
	"fmt"
	"strings"

	gocmd "github.com/go-cmd/cmd"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newSendEthCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "send-eth from to amount",
		Short: "",
		Long:  ``,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			accounts, err := utils.LoadJSON(viper.GetString("accounts"))
			cobra.CheckErr(err)

			fromAddr := accounts.Get(args[0]).Get("addr").String()
			fromPk := accounts.Get(args[0]).Get("pk").String()

			toAddr := accounts.Get(args[1]).Get("addr").String()
			if toAddr == "" {
				toAddr = args[1]
			}

			envCmd := gocmd.NewCmd("cast", "send", "--json", "--from", fromAddr, "--private-key", fromPk, "--value", args[2], "--gas-price", viper.GetString("gas-price"), toAddr)
			status := <-envCmd.Start()
			if len(status.Stderr) > 0 {
				return fmt.Errorf(strings.Join(status.Stderr, "\n"))
			}

			fmt.Println(strings.Join(status.Stdout, "\n"))

			return nil
		},
	}

	return cmd
}
