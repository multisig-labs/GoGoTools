package utilscmd

import (
	"github.com/multisig-labs/gogotools/pkg/configs"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init [avago-version] [subnet-evm-version]",
		Short: "Create default files in the current dir",
		Long: `Create default config files in the current dir, and
also attempt to download avalanchego and subnet-evm binaries from
GitHub.

Example:  ggt init v1.9.7 v0.4.8
`,
		Run: func(cmd *cobra.Command, args []string) {
			files := make(map[string]string)
			files["accounts.json"] = configs.Accounts
			files["cchain-config.json"] = configs.CChainConfig
			files["contracts.json"] = configs.Contracts
			files["node-config.json"] = configs.NodeConfig
			files["README.md"] = configs.Readme
			files["subnetevm-genesis.json"] = configs.SubnetEVMGenesis
			files["subnetevm-config.json"] = configs.SubnetEVMConfig
			files["xchain-config.json"] = configs.XChainConfig

			for fn, content := range files {
				if utils.FileExists(fn) {
					app.Log.Infof("File exists, skipping %s", fn)
				} else {
					app.Log.Infof("Creating %s", fn)
					utils.WriteFileBytes(fn, []byte(content))
				}
			}

			if len(args) > 0 && args[0] != "" {
				url, destFile, err := utils.DownloadAvalanchego(".", args[0])
				if err != nil {
					app.Log.Warnf("Error downloading %s: %s", url, err)
				} else {
					app.Log.Infof("Downloaded %s to %s", url, destFile)
				}
			}

			if len(args) > 1 && args[1] != "" {
				url, destFile, err := utils.DownloadSubnetevm(".", args[1])
				if err != nil {
					app.Log.Warnf("Error downloading %s: %s", url, err)
				} else {
					app.Log.Infof("Downloaded %s to %s", url, destFile)
				}
			}

		},
	}

	return cmd
}
