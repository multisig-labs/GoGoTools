package nodecmd

import (
	"context"

	"github.com/ava-labs/avalanchego/api/admin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newLogLevelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log-level level [name]",
		Short: "Set the log level for a running node",
		Long:  `Set the log level for a running node (DEBUG, INFO, ERROR) and optionally for a specific chain (X, P, C)`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var name string
			if len(args) > 1 {
				name = args[1]
			}
			return setLoggerLevel(args[0], name)
		},
	}
	return cmd
}

func setLoggerLevel(level string, name string) error {
	uri := viper.GetString("node-url")

	ctx := context.Background()

	c := admin.NewClient(uri)
	return c.SetLoggerLevel(ctx, name, level, level)
}
