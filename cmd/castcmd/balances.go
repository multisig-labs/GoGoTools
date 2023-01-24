package castcmd

import (
	"fmt"
	"strings"

	gocmd "github.com/go-cmd/cmd"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func newBalancesCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "balances",
		Short: "",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			accounts, err := utils.LoadJSON(viper.GetString("accounts"))
			cobra.CheckErr(err)

			balances := "{}"

			accounts.ForEach(func(key gjson.Result, value gjson.Result) bool {
				envCmd := gocmd.NewCmd("cast", "balance", value.Get("addr").String())
				status := <-envCmd.Start()
				if len(status.Stderr) > 0 {
					err = fmt.Errorf(strings.Join(status.Stderr, "\n"))
					return false
				}
				result := status.Stdout[0]
				ether := utils.ToDecimal(result, 18)
				balances, _ = sjson.Set(balances, key.String(), ether)
				return true
			})

			if err != nil {
				return err
			}
			fmt.Println(balances)
			return nil
		},
	}

	return cmd
}
