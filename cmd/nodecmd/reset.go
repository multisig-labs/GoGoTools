package nodecmd

import (
	"os"
	"path/filepath"

	"github.com/multisig-labs/gogotools/pkg/constants"
	"github.com/spf13/cobra"
)

func newResetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset work-dir",
		Short: "Nuke the data directory of the node",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			workDir := args[0]
			err := os.RemoveAll(filepath.Join(workDir, "data"))
			if err != nil {
				app.Log.Fatal("unable to delete data directory")
			}
			err = os.Mkdir(filepath.Join(workDir, "data"), constants.DefaultPerms755)
			if err != nil {
				app.Log.Fatal("unable to recreate data directory")
			}
			// TODO use .pid file to bounce node?
		},
	}
	return cmd
}
