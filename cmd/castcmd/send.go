package castcmd

import (
	"fmt"
	"os"
	"strings"

	gocmd "github.com/go-cmd/cmd"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newSendCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "send from contract fnSig [args]",
		Short: "Sign and pub a tx from a user in the accounts.json file to contract",
		Long:  `Use --verbose flag to see the full 'cast' command that gets run`,
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			accounts, err := utils.LoadJSON(viper.GetString("accounts"))
			cobra.CheckErr(err)

			contracts, err := utils.LoadJSON(viper.GetString("contracts"))
			cobra.CheckErr(err)

			fromAddr := accounts.Get(args[0]).Get("addr").String()
			fromPk := accounts.Get(args[0]).Get("pk").String()
			contractAddr := contracts.Get(args[1]).String()
			fnSig := args[2]

			// If any of the args have a user name, resolve to an addr
			args = utils.ResolveAccountAddrs(accounts, args)

			allArgs := []string{"send", "--json", "--from", fromAddr, "--private-key", fromPk, contractAddr, fnSig}
			allArgs = append(allArgs, args[3:]...)
			envCmd := gocmd.NewCmd("cast", allArgs...)

			if viper.GetBool("verbose") {
				fmt.Fprintf(os.Stderr, "%s %s\n\n", envCmd.Name, strings.Join(envCmd.Args, " "))
			}

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
