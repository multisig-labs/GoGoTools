package nodecmd

import (
	"context"
	"strings"

	"github.com/ava-labs/avalanchego/api/admin"
	"github.com/ava-labs/coreth/plugin/evm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func newLogLevelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log-level level chain-name",
		Short: "Set the log level for a running node and blockchain (X,P,C)",
		Long:  `Set the log level for a running node (DEBUG, INFO, ERROR) and for a specific chain (X, P, C, BlkchainID)`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setLoggerLevel(strings.ToLower(args[0]), strings.ToUpper(args[1]))
		},
	}
	return cmd
}

func setLoggerLevel(level string, chainName string) error {
	ctx := context.Background()
	uri := viper.GetString("node-url")

	if chainName == "X" || chainName == "P" {
		c := admin.NewClient(uri)
		_, err := c.SetLoggerLevel(ctx, chainName, level, level)
		return err
	}

	// Convert string to slog.Level
	var logLevel slog.Level
	err := logLevel.UnmarshalText([]byte(strings.ToUpper(level)))
	if err != nil {
		return err
	}

	c := evm.NewClient(uri, chainName)
	return c.SetLogLevel(ctx, logLevel)
}
