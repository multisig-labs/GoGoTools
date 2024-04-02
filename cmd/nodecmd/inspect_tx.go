package nodecmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func newInspectTXCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect-tx [txid]",
		Short: "Inspect a TX id",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			result, err := inspectTX(args[0])
			cobra.CheckErr(err)
			fmt.Println(result.String())
		},
	}
	return cmd
}

// It's not you, Types, it's me. I think we need a break for a bit.
func inspectTX(id string) (*gjson.Result, error) {
	uri := viper.GetString("node-url")
	urlP := fmt.Sprintf("%s/ext/bc/P", uri)

	tx, err := utils.FetchRPCGJSON(urlP, "platform.getTx", fmt.Sprintf(`{"txID":"%s","encoding":"json"}`, id))
	if err != nil {
		return nil, err
	}

	return tx, nil
}
