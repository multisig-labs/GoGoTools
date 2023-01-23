package nodecmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/multisig-labs/gogotools/pkg/process"
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
			c, a := nodeCmd(args[0])
			return runNodeAndWait(c, a)
		},
	}
	cmd.Flags().String("port", "9650", "Port that the node will listen on for API commands")

	return cmd
}

func nodeCmd(workDir string) (string, []string) {
	avaBin := filepath.Join(workDir, "bin", "avalanchego")
	dataPath := filepath.Join(workDir, "data")
	configsPath := filepath.Join(workDir, "configs")
	vmAliasesConfig := filepath.Join(workDir, "configs", "vms", "aliases.json")
	chainAliasesConfig := filepath.Join(workDir, "configs", "chains", "aliases.json")
	nodeConfig := filepath.Join(workDir, "configs", "node-config.json")
	pluginsPath := filepath.Join(workDir, "bin", "plugins")

	// TODO Not sure why we have to also specify --chain-config-dir etc, it should just be by default a child of --data-dir ?
	args := []string{
		"--http-host=0.0.0.0", // allow connections from anywhere
		fmt.Sprintf("--http-port=%s", viper.GetString("port")),
		"--public-ip=127.0.0.1", // this disables NAT
		"--bootstrap-ids=",      // dont try to connect to anyone else
		"--bootstrap-ips=",
		fmt.Sprintf("--data-dir=%s", dataPath),
		fmt.Sprintf("--config-file=%s", nodeConfig),
		fmt.Sprintf("--chain-config-dir=%s", configsPath),
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

	done := make(chan error)

	p := process.NewProcess(cmd, args, "")
	go func() {
		err := p.Start()
		if err != nil {
			app.Log.Fatalf("%v", err)
			done <- err
			return
		}
		fmt.Printf("Avalanche node listening on http://0.0.0.0:%s [PID: %d]\n", viper.GetString("port"), p.Process.Pid)
		done <- p.Wait()
	}()

	go func() {
		cSigTerm := make(chan os.Signal, 1)
		signal.Notify(cSigTerm, os.Interrupt, syscall.SIGTERM)
		<-cSigTerm
		fmt.Println("Sigterm recvd, shutting down...")
		p.Kill()
	}()

	err := <-done
	fmt.Println("Process completed.")
	// TODO fix extraneous error on interrupt
	return err
}
