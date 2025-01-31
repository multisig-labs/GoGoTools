package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/multisig-labs/gogotools/cmd/castcmd"
	"github.com/multisig-labs/gogotools/cmd/nodecmd"
	"github.com/multisig-labs/gogotools/cmd/subnetcmd"
	"github.com/multisig-labs/gogotools/cmd/utilscmd"
	"github.com/multisig-labs/gogotools/cmd/validatormanagercmd"
	"github.com/multisig-labs/gogotools/cmd/walletcmd"
	"github.com/multisig-labs/gogotools/cmd/warpcmd"
	"github.com/multisig-labs/gogotools/pkg/application"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	app     *application.GoGoTools
	cfgFile string
	verbose bool
)

func NewRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "ggt",
		Short: "GoGoTools, a utility belt for Avalanche developers",
		Long: `GoGoTools, a utility belt for Avalanche developers

To get started, run these commands in an empty directory:

  $ ggt utils init v1.9.7 v0.4.8
  $ ggt node prepare MyNodeV1 \
    --ava-bin=avalanchego-v1.9.7 \
    --vm-name=subnetevm \
    --vm-bin=subnet-evm-v0.4.8
  $ ggt node run MyNodeV1
  
  # Now in another terminal
  $ ggt wallet create-chain MyNodeV1 MyChain subnetevm genesis.json
  $ ggt node info

  # Set your ethereum rpc to your new subnet/blockchain
  $ export ETH_RPC_URL=$(ggt node info | jq -r '.rpcs.MyChain')

  # Run some commands against it
  $ cast chain-id
  $ ggt cast balances
  $ ggt cast send-eth owner alice 1ether
  $ ggt cast send owner NativeMinter "mintNativeCoin(address,uint256)" alice 1ether

  See the repo for more info. 
  https://github.com/multisig-labs/GoGoTools
  `,
		PersistentPreRunE: initApp,
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SilenceUsage = true // So cobra doesn't print usage when a command fails.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/ggt.json)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "output more verbose logs")
	rootCmd.PersistentFlags().String("node-url", "http://localhost:9650", "Avalanche node URL")

	err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	cobra.CheckErr(err)
	err = viper.BindPFlag("node-url", rootCmd.PersistentFlags().Lookup("node-url"))
	cobra.CheckErr(err)

	rootCmd.AddCommand(castcmd.NewCmd(app))
	rootCmd.AddCommand(nodecmd.NewCmd(app))
	rootCmd.AddCommand(subnetcmd.NewCmd(app))
	rootCmd.AddCommand(utilscmd.NewCmd(app))
	rootCmd.AddCommand(warpcmd.NewCmd(app))
	rootCmd.AddCommand(walletcmd.NewCmd(app))
	rootCmd.AddCommand(validatormanagercmd.NewCmd(app))
	rootCmd.AddCommand(versionCmd)
	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for default config.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(fmt.Sprintf("%s/.config", home))
		viper.SetConfigType("json")
		viper.SetConfigName("ggt")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		app.Log.Debugf("Using config file %s", viper.ConfigFileUsed())
	}
}

func initApp(_ *cobra.Command, _ []string) error {
	initConfig()
	if viper.GetBool("verbose") {
		app.Verbose()
	}
	return nil
}

// TODO figure out how to error properly
func Execute() {
	app = application.New()
	rootCmd := NewRootCmd()
	_ = rootCmd.Execute()
	// cobra.CheckErr(err)
}
