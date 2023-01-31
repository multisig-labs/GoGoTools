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
	"github.com/multisig-labs/gogotools/pkg/constants"
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
			workDir := args[0]
			viper.BindPFlags(cmd.Flags())

			exitIfRunning()

			// Truncate instead of delete so log tailing is not affected
			if viper.GetBool("clear-logs") {
				logsPath := filepath.Join(workDir, "data", "logs")
				logFiles, err := utils.FilePathWalk(logsPath, "log")
				cobra.CheckErr(err)
				for _, f := range logFiles {
					err := utils.Truncate(f, constants.DefaultPerms755)
					cobra.CheckErr(err)
				}
				cobra.CheckErr(err)
			}

			startCmd := filepath.Join(workDir, constants.BashScriptFilename)
			runNodeAndWait(workDir, startCmd)
		},
	}
	cmd.Flags().Bool("clear-logs", false, "Delete logs/* before starting node")
	cmd.Flags().Bool("watch", false, "(Experimental!) Watch data/bin and restart on any file changes")

	return cmd
}

func exitIfRunning() {
	if utils.FileExists(".pid") {
		app.Log.Fatalf(".pid file exists, is another node already running? Delete .pid and try again.")
	}
}

func runNodeAndWait(workDir string, cmd string) error {
	for {
		if err := runNode(workDir, cmd); err != nil {
			return err
		}
	}
}

// TODO This whole thing is wonky. Maybe do the job stuff in start.sh?
// Goals:
//   - Ability to start a node and wait for Ctl-C to (reliably) quit
//   - Ability to have USR1 gracefully stop and restart the node
//   - Listen for changes to vm files and gracefully stop and restart the node
//   - Good UX and error reporting so user knows whats happening at all times
//   - Use start.sh so user can see exactly how avalanchego is being started
//
// Issues:
//   - If we have gocmd discard stdout, then if node fails to start we have nothing to show user
//   - Maybe just tail main.log to show any errors if node doesnt start?
//   - If we keep stdout, it eats major memory unless we have something that clears it out occassionally
//   - Is it worth using gocmd? Is there something better? Roll our own?
func runNode(workDir string, cmd string) error {
	var envCmd *gocmd.Cmd
	var finalStatus gocmd.Status
	var shouldRestart bool

	// if viper.GetBool("verbose") {
	envCmd = gocmd.NewCmdOptions(gocmd.Options{}, cmd)
	// } else {
	// envCmd = gocmd.NewCmdOptions(gocmd.Options{Buffered: false, Streaming: false}, cmd)
	// }

	// Ctl-C wil stop the node
	cSigTerm := make(chan os.Signal, 1)
	signal.Notify(cSigTerm, os.Interrupt, syscall.SIGTERM)

	// USR1 can be sent to restart the node, use https://github.com/watchexec/watchexec
	cSigUser1 := make(chan os.Signal, 1)
	signal.Notify(cSigUser1, syscall.SIGUSR1)

	// Start the node
	app.Log.Info("Starting node...")
	statusChan := envCmd.Start()
	doneChan := envCmd.Done()

	app.Log.Infof("Avalanche node listening on http://0.0.0.0:9650")
	app.Log.Infof("(Send USR1 to PID %d to restart the node)", os.Getpid())
	app.Log.Infof("(If you have problems you can always run '%s/start.sh' directly)", workDir)
	// TODO dont show the below if they have already created a subnet
	app.Log.Infof("In another terminal, run this command to create a subnetEVM")
	app.Log.Infof("  ggt wallet create-chain %s MyChainName subnetevm", workDir)

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
	defer func() { os.Remove(".pid") }()

	for {
		select {
		case <-cSigUser1:
			shouldRestart = true
			app.Log.Debug("Recvd SIGUSR1")
			err := envCmd.Stop()
			if err != nil {
				app.Log.Errorf("error stopping node (USR1): %s", err)
			}
		case <-cSigTerm:
			app.Log.Debug("Recvd SIGTERM")
			err := envCmd.Stop()
			if err != nil {
				app.Log.Errorf("error stopping node (SIGTERM): %s", err)
			}
		case finalStatus = <-statusChan:
			app.Log.Debug("Recvd statuschan")
			app.Log.Info(strings.Join(finalStatus.Stdout, "\n"))
		case <-doneChan:
			app.Log.Debug("Recvd donechan")
			app.Log.Debugf("%+v", finalStatus)
			if shouldRestart {
				app.Log.Info("Restarting node...")
				time.Sleep(time.Second * 5)
				return nil
			} else {
				if finalStatus.Exit > 0 {
					app.Log.Infof("program exited with code: %d (check logs for more info)", finalStatus.Exit)
				}
				// Normal exit, but we have to return err so we break out of loop and dont restart
				return fmt.Errorf("program exited")
			}
		}
	}
}
