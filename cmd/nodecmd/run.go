package nodecmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	gocmd "github.com/go-cmd/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run work-dir",
		Short: "Run avalanchego from a previously prepared directory",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			if viper.GetBool("clear-logs") {
				logsPath := filepath.Join(args[0], "data", "logs")
				err := os.RemoveAll(logsPath)
				cobra.CheckErr(err)
			}
			c, a := nodeCmd(args[0])
			return runNodeAndWait(c, a)
		},
	}
	cmd.Flags().String("port", "9650", "Port that the node will listen on for API commands")
	cmd.Flags().Bool("clear-logs", false, "Delete data/logs/* before starting node")

	return cmd
}

func nodeCmd(workDir string) (string, []string) {
	// TODO combine this config with prepare
	avaBin := filepath.Join(workDir, "bin", "avalanchego")
	dataPath := filepath.Join(workDir, "data")
	// configsPath := filepath.Join(workDir, "configs")
	chainConfigsPath := filepath.Join(workDir, "configs", "chains")
	vmAliasesConfig := filepath.Join(workDir, "configs", "vms", "aliases.json")
	chainAliasesConfig := filepath.Join(workDir, "configs", "chains", "aliases.json")
	nodeConfig := filepath.Join(workDir, "configs", "node-config.json")
	pluginsPath := filepath.Join(workDir, "bin", "plugins")

	// TODO Not sure why we have to also specify --chain-config-dir etc, it should just be by default a child of --data-dir ?
	args := []string{
		"--http-host=0.0.0.0", // allow connections from anywhere
		fmt.Sprintf("--http-port=%s", viper.GetString("port")),
		"--log-display-level=off", // only log to files
		"--public-ip=127.0.0.1",   // this disables NAT
		"--bootstrap-ids=",        // dont try to connect to anyone else
		"--bootstrap-ips=",
		fmt.Sprintf("--data-dir=%s", dataPath),
		fmt.Sprintf("--config-file=%s", nodeConfig),
		fmt.Sprintf("--chain-config-dir=%s", chainConfigsPath),
		fmt.Sprintf("--plugin-dir=%s", pluginsPath),
		fmt.Sprintf("--vm-aliases-file=%s", vmAliasesConfig),
		fmt.Sprintf("--chain-aliases-file=%s", chainAliasesConfig),
	}
	return avaBin, args
}

// TODO make a properly shell-escaped cmd to print out if user wants to copy paste it somewhere
// func displayCmd(cmd string args []string) string {
// }

func runNodeAndWait(cmd string, args []string) error {
	var envCmd *gocmd.Cmd

	// If the node starts sucessfully, then we want to throw away stdout as to not clutter the terminal
	// but if there is an error starting the node then we want to show the user any error messages.
	// TODO Figure out a better approach here.
	if viper.GetBool("verbose") {
		envCmd = gocmd.NewCmdOptions(gocmd.Options{Buffered: true}, cmd, args...)
	} else {
		envCmd = gocmd.NewCmdOptions(gocmd.Options{Buffered: false, Streaming: false}, cmd, args...)
	}

	statusChan := envCmd.Start()

	// Respond to Ctl-C
	go func() {
		cSigTerm := make(chan os.Signal, 1)
		signal.Notify(cSigTerm, os.Interrupt, syscall.SIGTERM)
		<-cSigTerm
		fmt.Println("Sigterm recvd, shutting down...")
		envCmd.Stop()
	}()

	fmt.Printf("Avalanche node listening on http://0.0.0.0:%s\n\n", viper.GetString("port"))

	// Wait for node to exit
	finalStatus := <-statusChan

	if finalStatus.Exit > 0 {
		fmt.Println(strings.Join(finalStatus.Stdout, "\n"))
		return fmt.Errorf("program exited with code: %d (use --verbose flag for more info)", finalStatus.Exit)
	}
	fmt.Println(strings.Join(finalStatus.Stdout, "\n"))
	return finalStatus.Error
}
