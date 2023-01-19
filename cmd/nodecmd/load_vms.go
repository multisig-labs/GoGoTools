package nodecmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newLoadVMsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load-vms",
		Short: "Dynamically loads any virtual machines installed on the node as plugins.",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadVMs()
		},
	}
	return cmd
}

func loadVMs() error {
	// TODO uri in env
	uri := viper.GetString("node-url")
	urlAdmin := fmt.Sprintf("%s/ext/admin", uri)

	loadVms, err := utils.FetchRPCGJSON(urlAdmin, "admin.loadVMs", "")
	cobra.CheckErr(err)
	fmt.Println(loadVms.Get("result").String())
	return nil
}
