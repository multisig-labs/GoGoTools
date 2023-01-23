package nodecmd

import (
	"context"

	"github.com/ava-labs/avalanchego/api"
	"github.com/ava-labs/avalanchego/api/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCreateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-user username password pk",
		Short: "Create a user in the node's keystore and import pk",
		Long:  ``,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createUser(args[0], args[1], args[2])
		},
	}
	return cmd
}

// TODO import pk? Is that useful?

func createUser(username string, password string, pk string) error {
	uri := viper.GetString("node-url")

	ctx := context.Background()

	up := api.UserPass{
		Username: username,
		Password: password,
	}

	c := keystore.NewClient(uri)
	return c.CreateUser(ctx, up)
}
