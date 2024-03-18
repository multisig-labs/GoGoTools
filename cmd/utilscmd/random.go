package utilscmd

import (
	"crypto/rand"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/spf13/cobra"
)

func newRandomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "random",
		Short: "Generate random ids of various types",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmd.AddCommand(newRandomNodeIDCmd())

	return cmd
}

func newRandomNodeIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodeid",
		Short: "Create random NodeID",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			r := make([]byte, 20)
			_, err := rand.Read(r)
			cobra.CheckErr(err)
			nodeid := ids.NodeID(r)
			fmt.Println(nodeid)
		},
	}
	return cmd
}
