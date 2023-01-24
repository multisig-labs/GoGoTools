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
		Short: "Create default files in the current dir (accounts.json, contracts.json, genesis.json)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			var fn string

			fn = "accounts.json"
			if utils.FileExists(fn) {
				fmt.Printf("File exists, skipping %s\n", fn)
			} else {
				fmt.Printf("Creating %s\n", fn)
				utils.WriteFileBytes(fn, []byte(configs.Accounts))
			}

			fn = "contracts.json"
			if utils.FileExists(fn) {
				fmt.Printf("File exists, skipping %s\n", fn)
			} else {
				fmt.Printf("Creating %s\n", fn)
				utils.WriteFileBytes(fn, []byte(configs.Contracts))
			}

			fn = "genesis.json"
			if utils.FileExists(fn) {
				fmt.Printf("File exists, skipping %s\n", fn)
			} else {
				fmt.Printf("Creating %s\n", fn)
				utils.WriteFileBytes(fn, []byte(configs.GenesisSubnetEVM))
			}
		},
	}

	return cmd
}
