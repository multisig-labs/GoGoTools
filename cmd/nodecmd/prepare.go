package nodecmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/multisig-labs/gogotools/pkg/configs"
	"github.com/multisig-labs/gogotools/pkg/constants"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

func newPrepareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare work-dir",
		Short: "Create a new self-contained directory in [work-dir] for a node",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			if viper.GetString("ava-bin") == "" {
				return fmt.Errorf("must supply --ava-bin flag or AVA_BIN env")
			}
			if exists := utils.FileExists(viper.GetString("ava-bin")); !exists {
				return fmt.Errorf("ava-bin file does not exist: %s", viper.GetString("ava-bin"))
			}

			if viper.GetString("vm-bin") == "" {
				app.Log.Warnln("WARNING: --vm-bin or VM_BIN not supplied, creating avalanchego node without any subnet vms")
			} else {
				if exists := utils.FileExists(viper.GetString("vm-bin")); !exists {
					return fmt.Errorf("vm-bin file does not exist: %s", viper.GetString("vm-bin"))
				}
			}

			if err := prepareWorkDir(args[0], viper.GetString("ava-bin"), viper.GetString("vm-bin"), viper.GetString("vm-name")); err != nil {
				return err
			}

			app.Log.Infof("Success! run 'ggt node run %s' to start the node", args[0])
			return nil
		},
	}
	cmd.Flags().String("ava-bin", "", "Location of avalanchego binary (also AVA_BIN)")
	cmd.Flags().String("vm-bin", "", "(optional) Location of subnetevm binary (also VM_BIN)")
	cmd.Flags().String("vm-name", "subnetevm", "(optional) Name of vm (also VM_NAME)")
	return cmd
}

func prepareWorkDir(workDir string, avaBin string, vmBin string, vmName string) error {
	if _, err := os.Stat(workDir); err == nil {
		return fmt.Errorf("%s exists, aborting", workDir)
	}
	err := mkDirs(workDir)
	cobra.CheckErr(err)

	dirStruct := utils.NewDirectoryLayout(workDir)
	fileLocations := utils.NewFileLocations(workDir)

	bash, err := prepareBashScript(workDir)
	cobra.CheckErr(err)
	app.Log.Infof("Creating %s", filepath.Join(workDir, constants.BashScriptFilename))
	err = ioutil.WriteFile(filepath.Join(workDir, constants.BashScriptFilename), []byte(bash), constants.DefaultPerms755)
	cobra.CheckErr(err)

	app.Log.Infof("Linking %s to %s", avaBin, fileLocations.AvaBinFile)
	if err := utils.LinkFile(avaBin, fileLocations.AvaBinFile); err != nil {
		return fmt.Errorf("failed linking file: %w", err)
	}

	// Copy configs from cur dir where `ggt init` put some defaults.
	app.Log.Infof("Copying %s to %s", constants.NodeConfigFilename, fileLocations.ConfigFile)
	err = utils.CopyFile(constants.NodeConfigFilename, fileLocations.ConfigFile)
	cobra.CheckErr(err)
	app.Log.Infof("Copying %s to %s", constants.CChainConfigFilename, fileLocations.CChainConfigFile)
	err = utils.CopyFile(constants.CChainConfigFilename, fileLocations.CChainConfigFile)
	cobra.CheckErr(err)
	app.Log.Infof("Copying %s to %s", constants.XChainConfigFilename, fileLocations.XChainConfigFile)
	err = utils.CopyFile(constants.XChainConfigFilename, fileLocations.XChainConfigFile)
	cobra.CheckErr(err)

	// Always write a vm aliases file even if empty to make avalanchego happy
	vmAliases := "{}"
	if vmBin != "" {
		// TODO make this a fn
		paddedBytes := [32]byte{}
		copy(paddedBytes[:], []byte(vmName))
		vmID, err := ids.ToID(paddedBytes[:])
		if err != nil {
			return err
		}

		fn := filepath.Join(dirStruct.PluginDir, vmID.String())
		app.Log.Infof("Linking %s to %s", vmBin, fn)
		if err := utils.LinkFile(vmBin, fn); err != nil {
			return fmt.Errorf("failed linking file %w", err)
		}

		vmAliases, _ = sjson.Set("{}", vmID.String(), []string{vmName})
	}
	app.Log.Infof("Creating %s", fileLocations.VMAliasesFile)
	err = ioutil.WriteFile(fileLocations.VMAliasesFile, []byte(vmAliases), 0644)
	cobra.CheckErr(err)

	app.Log.Infof("Creating %s", fileLocations.ChainAliasesFile)
	err = ioutil.WriteFile(fileLocations.ChainAliasesFile, []byte("{}"), 0644)
	cobra.CheckErr(err)

	return err
}

type bashCmdParams struct {
	utils.DirectoryLayout
	utils.FileLocations
}

func prepareBashScript(workDir string) (string, error) {
	layout := utils.NewDirectoryLayout("")
	files := utils.NewFileLocations("")
	params := bashCmdParams{layout, files}
	bash := configs.StartBash
	buf := &bytes.Buffer{}
	t, err := template.New("").Parse(bash)
	if err != nil {
		return "", nil
	}
	err = t.Execute(buf, params)
	return buf.String(), err
}

func mkDirs(workDir string) error {
	dirStruct := utils.NewDirectoryLayout(workDir)
	values := reflect.ValueOf(dirStruct)
	for i := 0; i < values.NumField(); i++ {
		err := os.MkdirAll(values.Field(i).String(), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
