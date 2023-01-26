package utilscmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/configs"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create default config files in the current dir",
		Long:  `iklkjh`,
		Run: func(cmd *cobra.Command, args []string) {
			files := make(map[string]string)
			files["accounts.json"] = configs.Accounts
			files["contracts.json"] = configs.Contracts
			files["subnetevm-genesis.json"] = configs.SubnetEVMGenesis
			files["subnetevm-config.json"] = configs.SubnetEVMConfig
			files["node-config.json"] = configs.NodeConfig
			files["README.md"] = configs.Readme

			for fn, content := range files {
				if utils.FileExists(fn) {
					fmt.Printf("File exists, skipping %s\n", fn)
				} else {
					fmt.Printf("Creating %s\n", fn)
					utils.WriteFileBytes(fn, []byte(content))
				}
			}
		},
	}

	return cmd
}
