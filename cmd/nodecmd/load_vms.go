package nodecmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func newLoadVMsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load-vms",
		Short: "Dynamically loads any virtual machines installed on the node as plugins.",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			result, err := loadVMs()
			cobra.CheckErr(err)
			fmt.Println(result.Get("result").String())
		},
	}
	return cmd
}

func loadVMs() (*gjson.Result, error) {
	uri := viper.GetString("node-url")
	urlAdmin := fmt.Sprintf("%s/ext/admin", uri)

	return utils.FetchRPCGJSON(urlAdmin, "admin.loadVMs", "")
}
