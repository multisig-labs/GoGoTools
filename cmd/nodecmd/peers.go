package nodecmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func newPeersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "peers",
		Short: "Get peers for a running node",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			result, err := getPeers()
			cobra.CheckErr(err)
			fmt.Println(result.String())
		},
	}
	return cmd
}

// It's not you, Types, it's me. I think we need a break for a bit.
func getPeers() (*gjson.Result, error) {
	uri := viper.GetString("node-url")
	urlInfo := fmt.Sprintf("%s/ext/info", uri)

	peers, err := utils.FetchRPCGJSON(urlInfo, "info.peers", "")
	if err != nil {
		return nil, err
	}

	out := "{}"
	out, _ = sjson.SetRaw(out, "peers", peers.Get("result.peers").String())
	out, _ = sjson.Set(out, "numPeers", peers.Get("result.numPeers").Int())

	result := gjson.Parse(out)
	return &result, nil
}
