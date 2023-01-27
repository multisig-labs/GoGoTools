package nodecmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	gocmd "github.com/go-cmd/cmd"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run work-dir",
		Short: "Run avalanchego from a previously prepared directory",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
			if viper.GetBool("clear-logs") {
				logsPath := filepath.Join(args[0], "data", "logs")
				err := os.RemoveAll(logsPath)
				cobra.CheckErr(err)
			}
			workDir := args[0]
			c, a := nodeCmd(workDir)
			runNodeAndWait(workDir, c, a)
		},
	}
	cmd.Flags().String("port", "9650", "Port that the node will listen on for API commands")
	cmd.Flags().Bool("clear-logs", false, "Delete data/logs/* before starting node")
	cmd.Flags().Bool("watch", false, "(Experimental!) Watch data/bin and restart on any file changes")

	return cmd
}

func nodeCmd(workDir string) (string, []string) {
	// TODO combine this config with prepare
	avaBin := filepath.Join(workDir, "bin", "avalanchego")
	dataPath := filepath.Join(workDir, "data")
	chainConfigsPath := filepath.Join(workDir, "configs", "chains")
	// subnetConfigsPath := filepath.Join(workDir, "configs", "subnets")
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
		// fmt.Sprintf("--subnet-config-dir=%s", subnetConfigsPath),
		fmt.Sprintf("--plugin-dir=%s", pluginsPath),
		fmt.Sprintf("--vm-aliases-file=%s", vmAliasesConfig),
		fmt.Sprintf("--chain-aliases-file=%s", chainAliasesConfig),
		// fmt.Sprintf("--track-subnets=p4jUwqZsA2LuSftroCd3zb4ytH8W99oXKuKVZdsty7eQ3rXD6"),
	}
	return avaBin, args
}

// TODO make a properly shell-escaped cmd to print out if user wants to copy paste it somewhere
// func displayCmd(cmd string args []string) string {
// }

func runNodeAndWait(workDir string, cmd string, args []string) error {
	for {
		if err := runNode(workDir, cmd, args); err != nil {
			return err
		}
	}
}

// TODO This whole thing is wonky.
// Goals:
//   - Ability to start a node and wait for Ctl-C to (reliably) quit
//   - Ability to have USR1 gracefully stop and restart the node
//   - Listen for changes to vm files and gracefully stop and restart the node
//   - Good UX and error reporting so user knows whats happening at all times
//
// Issues:
//   - If we have gocmd discard stdout, then if node fails to start we have nothing to show user
//   - If we keep stdout, it eats major memory unless we have something that clears it out occassionally
//   - Is it worth using gocmd? Is there something better? Roll our own?
func runNode(workDir string, cmd string, args []string) error {
	var envCmd *gocmd.Cmd
	var finalStatus gocmd.Status
	var shouldRestart bool

	// If the node starts sucessfully, then we want to throw away stdout as to not clutter the terminal
	// but if there is an error starting the node then we want to show the user any error messages.
	// TODO Figure out a better approach here.
	if viper.GetBool("verbose") {
		envCmd = gocmd.NewCmdOptions(gocmd.Options{Buffered: true}, cmd, args...)
	} else {
		envCmd = gocmd.NewCmdOptions(gocmd.Options{Buffered: false, Streaming: false}, cmd, args...)
	}

	// Ctl-C wil stop the node
	cSigTerm := make(chan os.Signal, 1)
	signal.Notify(cSigTerm, os.Interrupt, syscall.SIGTERM)

	// USR1 can be sent to restart the node, use https://github.com/watchexec/watchexec
	cSigUser1 := make(chan os.Signal, 1)
	signal.Notify(cSigUser1, syscall.SIGUSR1)

	// Start the node
	statusChan := envCmd.Start()
	doneChan := envCmd.Done()

	fmt.Printf("Avalanche node listening on http://0.0.0.0:%s\n", viper.GetString("port"))
	fmt.Printf("(Send USR1 to PID %d to restart the node)\n\n", os.Getpid())
	fmt.Printf("In another terminal, run this command to create a subnetEVM\n")
	fmt.Printf("  ggt wallet create-chain MyNodeName MyChainName subnetevm\n")

	// TODO this doesnt quite work yet
	if viper.GetBool("watch") {
		fmt.Println("Watching bin/ will restart on changes")

		w := watcher.New()
		w.IgnoreHiddenFiles(true)
		if err := w.AddRecursive(fmt.Sprintf("%s/bin", workDir)); err != nil {
			cobra.CheckErr(err)
		}

		go func() {
			for {
				select {
				case event := <-w.Event:
					fmt.Printf("%+v\n", event)
					cSigUser1 <- syscall.SIGUSR1
				case err := <-w.Error:
					fmt.Println(err)
				case <-w.Closed:
					return
				}
			}
		}()

		go func() {
			if err := w.Start(time.Second * 2); err != nil {
				cobra.CheckErr(err)
			}
		}()

		// for path, f := range w.WatchedFiles() {
		// 	fmt.Printf("%s: %s\n", path, f.Name())
		// }
	}

	// Write a .pid so other commands can restart us with a USR1 if necessary
	utils.WriteFileBytes(".pid", []byte(fmt.Sprintf("%d", os.Getpid())))

	for {
		select {
		case <-cSigUser1:
			shouldRestart = true
			envCmd.Stop()
		case <-cSigTerm:
			envCmd.Stop()
		case finalStatus = <-statusChan:
			fmt.Println(strings.Join(finalStatus.Stdout, "\n"))
		case <-doneChan:
			os.Remove(".pid")
			if shouldRestart {
				fmt.Println("Restarting node...")
				time.Sleep(time.Second * 3)
				return nil
			} else {
				if finalStatus.Exit > 0 {
					fmt.Printf("program exited with code: %d (use --verbose flag for more info)", finalStatus.Exit)
				}
				// Normal exit, but we have to return err so we break out of loop and dont restart
				return fmt.Errorf("program exited")
			}
		}
	}
}
