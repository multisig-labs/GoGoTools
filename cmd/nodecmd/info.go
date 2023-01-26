package nodecmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get all info for a running node in a single JSON blob",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getInfo()
		},
	}
	return cmd
}

// It's not you, Types, it's me. I think we need a break for a bit.
func getInfo() error {
	uri := viper.GetString("node-url")
	urlInfo := fmt.Sprintf("%s/ext/info", uri)
	urlP := fmt.Sprintf("%s/ext/bc/P", uri)
	urlAdmin := fmt.Sprintf("%s/ext/admin", uri)

	getNetworkName, err := utils.FetchRPCGJSON(urlInfo, "info.getNetworkName", "")
	cobra.CheckErr(err)

	getNetworkID, err := utils.FetchRPCGJSON(urlInfo, "info.getNetworkID", "")
	cobra.CheckErr(err)

	getNodeID, err := utils.FetchRPCGJSON(urlInfo, "info.getNodeID", "")
	cobra.CheckErr(err)

	getNodeVersion, err := utils.FetchRPCGJSON(urlInfo, "info.getNodeVersion", "")
	cobra.CheckErr(err)

	getVMs, err := utils.FetchRPCGJSON(urlInfo, "info.getVMs", "")
	cobra.CheckErr(err)

	getUptime, err := utils.FetchRPCGJSON(urlInfo, "info.uptime", "")
	cobra.CheckErr(err)

	getBlockchains, err := utils.FetchRPCGJSON(urlP, "platform.getBlockchains", "")
	cobra.CheckErr(err)

	getSubnets, err := utils.FetchRPCGJSON(urlP, "platform.getSubnets", "")
	cobra.CheckErr(err)

	aliases := `{"blockchainAliases":"AdminAPI disabled on node"}`
	getBlockchains.Get("result.blockchains").ForEach(func(key, value gjson.Result) bool {
		// println(value.Get("id").String())
		blockchainID := value.Get("id").String()
		blockchainAliases, err := utils.FetchRPCGJSON(urlAdmin, "admin.getChainAliases", fmt.Sprintf(`{"chain":"%s"}`, blockchainID))
		if err != nil {
			// Maybe Admin API is disabled on this node, skip it.
			return false
		}
		// If subnet didnt start for some reason, this will be blank
		s := blockchainAliases.Get("result.aliases").String()
		if s == "" {
			s = `["ERROR starting blockchain"]`
		}
		aliases, _ = sjson.SetRaw(aliases, fmt.Sprintf("blockchainAliases.%s", blockchainID), s)
		return true
	})

	rpcs := "{}"
	getBlockchains.Get("result.blockchains").ForEach(func(key, value gjson.Result) bool {
		if value.Get("subnetID").String() != "11111111111111111111111111111111LpoYY" {
			blockchainID := value.Get("id").String()
			name := value.Get("name").String()
			url := viper.GetString("node-url")
			rpcs, _ = sjson.Set(rpcs, name, fmt.Sprintf("%s/ext/bc/%s/rpc", url, blockchainID))
		}
		return true
	})

	out := "{}"
	out, _ = sjson.Set(out, "nodeID", getNodeID.Get("result.nodeID").String())
	out, _ = sjson.Set(out, "networkID", getNetworkID.Get("result.networkID").Int())
	out, _ = sjson.Set(out, "networkName", getNetworkName.Get("result.networkName").String())
	out, _ = sjson.SetRaw(out, "uptime", getUptime.Get("result").String())
	out, _ = sjson.SetRaw(out, "getNodeVersion", getNodeVersion.Get("result").String())
	out, _ = sjson.SetRaw(out, "getVMs", getVMs.Get("result").String())
	out, _ = sjson.SetRaw(out, "subnets", getSubnets.Get("result.subnets").String())
	out, _ = sjson.SetRaw(out, "blockchains", getBlockchains.Get("result.blockchains").String())
	out, _ = sjson.SetRaw(out, "aliases", aliases)
	out, _ = sjson.SetRaw(out, "rpcs", rpcs)

	fmt.Println(gjson.Parse(out).String())
	return nil
}
