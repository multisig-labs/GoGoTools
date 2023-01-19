package nodecmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/multisig-labs/gogotools/pkg/configs"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

func newPrepareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare [work-dir]",
		Short: "",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			if viper.GetString("ava-bin") == "" {
				return fmt.Errorf("must supply --ava-bin flag or AVA_BIN env")
			}
			if err := prepareWorkDir(args[0], viper.GetString("ava-bin"), viper.GetString("vm-bin"), viper.GetString("vm-name")); err != nil {
				return err
			}
			fmt.Printf("Success! run 'ggt node run %s' to start the node", args[0])
			return nil
		},
	}
	cmd.Flags().String("ava-bin", "", "Location of avalanchego binary")
	cmd.Flags().String("vm-bin", "", "Location of subnetevm binary")
	cmd.Flags().String("vm-name", "subnetevm", "Name of vm")
	cmd.Flags().String("node-config", "", "Location of node config file")
	return cmd
}

func prepareWorkDir(workDir string, avaBin string, vmBin string, vmName string) error {
	if _, err := os.Stat(workDir); err == nil {
		return fmt.Errorf("%s exists, aborting", workDir)
	}

	binPath := filepath.Join(workDir, "bin")
	pluginsPath := filepath.Join(workDir, "bin", "plugins")
	dataPath := filepath.Join(workDir, "data")
	configsPath := filepath.Join(workDir, "configs")
	configsVmsPath := filepath.Join(workDir, "configs", "vms")
	configsChainsPath := filepath.Join(workDir, "configs", "chains")

	dirList := []string{binPath, pluginsPath, dataPath, configsPath, configsVmsPath, configsChainsPath}
	for i := 0; i < len(dirList); i++ {
		err := os.MkdirAll(dirList[i], os.ModePerm)
		if err != nil {
			return err
		}
	}

	fn := filepath.Join(binPath, "avalanchego")
	if err := utils.LinkFile(avaBin, fn); err != nil {
		return fmt.Errorf("failed linking file '%s' to '%s': %w", avaBin, fn, err)
	}

	fn = filepath.Join(configsPath, "node-config.json")
	ioutil.WriteFile(fn, []byte(configs.NodeConfig), 0644)

	fn = filepath.Join(configsChainsPath, "aliases.json")
	ioutil.WriteFile(fn, []byte("{}"), 0644)

	if vmBin != "" {
		paddedBytes := [32]byte{}
		copy(paddedBytes[:], []byte(vmName))
		vmID, err := ids.ToID(paddedBytes[:])
		if err != nil {
			return err
		}

		fn = filepath.Join(pluginsPath, vmID.String())
		if err := utils.LinkFile(vmBin, fn); err != nil {
			return fmt.Errorf("failed linking file '%s' to '%s': %w", vmBin, fn, err)
		}

		vmAliases, _ := sjson.Set("{}", vmID.String(), []string{vmName})
		fn = filepath.Join(configsVmsPath, "aliases.json")
		ioutil.WriteFile(fn, []byte(vmAliases), 0644)
	}

	return nil
}