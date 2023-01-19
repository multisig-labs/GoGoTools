package nodecmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/multisig-labs/gogotools/pkg/process"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [work-dir]",
		Short: "",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, a := nodeCmd(args[0])
			return runNodeAndWait(c, a)
		},
	}
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

	args := []string{
		fmt.Sprintf("--data-dir=%s", dataPath),
		fmt.Sprintf("--config-file=%s", nodeConfig),
		fmt.Sprintf("--chain-config-dir=%s", configsPath),
		fmt.Sprintf("--plugin-dir=%s", pluginsPath),
		fmt.Sprintf("--vm-aliases-file=%s", vmAliasesConfig),
		fmt.Sprintf("--chain-aliases-file=%s", chainAliasesConfig),
	}
	return avaBin, args
}

// TODO make an escaped cmd to print out for copy paste
// func displayCmd(cmd string args []string) string {
// }

func runNodeAndWait(cmd string, args []string) error {

	done := make(chan error)

	p := process.NewProcess(cmd, args, "")
	go func() {
		err := p.Start()
		if err != nil {
			done <- err
			return
		}
		fmt.Printf("Avalanche node started with PID: %d\n", p.Process.Pid)
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
