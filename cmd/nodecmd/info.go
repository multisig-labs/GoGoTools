package nodecmd

import (
	"fmt"
	"net/url"

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
		Run: func(cmd *cobra.Command, args []string) {
			result, err := getInfo()
			cobra.CheckErr(err)
			fmt.Println(result.String())
		},
	}
	return cmd
}

// It's not you, Types, it's me. I think we need a break for a bit.
func getInfo() (*gjson.Result, error) {
	nodeURL := viper.GetString("node-url")
	parsedURL, err := url.Parse(nodeURL)
	if err != nil {
		return nil, fmt.Errorf("invalid node-url: %w", err)
	}

	// Store the query string if present
	queryString := parsedURL.RawQuery
	if queryString != "" {
		queryString = "?" + queryString
	}

	// Create base URL without query string
	parsedURL.RawQuery = ""
	baseURL := parsedURL.String()

	urlInfo := fmt.Sprintf("%s/ext/info%s", baseURL, queryString)
	urlP := fmt.Sprintf("%s/ext/bc/P%s", baseURL, queryString)
	urlAdmin := fmt.Sprintf("%s/ext/admin%s", baseURL, queryString)

	getNetworkName, err := utils.FetchRPCGJSON(urlInfo, "info.getNetworkName", "")
	if err != nil {
		return nil, err
	}

	getNetworkID, err := utils.FetchRPCGJSON(urlInfo, "info.getNetworkID", "")
	if err != nil {
		return nil, err
	}

	getNodeID, err := utils.FetchRPCGJSON(urlInfo, "info.getNodeID", "")
	if err != nil {
		return nil, err
	}

	getNodeVersion, err := utils.FetchRPCGJSON(urlInfo, "info.getNodeVersion", "")
	if err != nil {
		return nil, err
	}

	getVMs, err := utils.FetchRPCGJSON(urlInfo, "info.getVMs", "")
	if err != nil {
		return nil, err
	}

	getUptime, err := utils.FetchRPCGJSON(urlInfo, "info.uptime", "")
	if err != nil {
		return nil, err
	}

	getBlockchains, err := utils.FetchRPCGJSON(urlP, "platform.getBlockchains", "")
	if err != nil {
		return nil, err
	}

	getSubnets, err := utils.FetchRPCGJSON(urlP, "platform.getSubnets", "")
	if err != nil {
		return nil, err
	}

	stakingAssetIDs := `{}`
	getSubnets.Get("result.subnets").ForEach(func(key, value gjson.Result) bool {
		subnetID := value.Get("id").String()
		stakingAsset, err := utils.FetchRPCGJSON(urlP, "platform.getStakingAssetID", fmt.Sprintf(`{"subnetID":"%s"}`, subnetID))
		if err != nil {
			app.Log.Infof("error retrieving stakingAssetID for subnetID: %s", subnetID)
			return false
		}
		id := stakingAsset.Get("result.assetID").String()
		if id != "" {
			stakingAssetIDs, _ = sjson.Set(stakingAssetIDs, subnetID, id)
		}
		return true
	})

	aliases := `{"blockchainAliases":"AdminAPI disabled on node"}`
	getBlockchains.Get("result.blockchains").ForEach(func(key, value gjson.Result) bool {
		blockchainID := value.Get("id").String()
		blockchainAliases, err := utils.FetchRPCGJSON(urlAdmin, "admin.getChainAliases", fmt.Sprintf(`{"chain":"%s"}`, blockchainID))
		if err != nil {
			// Maybe Admin API is disabled on this node, skip it.
			return false
		}
		// If subnet didnt start for some reason, this will be blank
		s := blockchainAliases.Get("result.aliases").String()
		if s == "" {
			s = `["blockchain not started, check logs"]`
		}

		aliases, _ = sjson.SetRaw(aliases, fmt.Sprintf("blockchainAliases.%s", blockchainID), s)
		return true
	})

	url := viper.GetString("node-url")
	rpcs := fmt.Sprintf(`{"C":"%s/ext/bc/C/rpc"}`, url)
	getBlockchains.Get("result.blockchains").ForEach(func(key, value gjson.Result) bool {
		if value.Get("subnetID").String() != "11111111111111111111111111111111LpoYY" {
			blockchainID := value.Get("id").String()
			name := value.Get("name").String()
			rpcs, _ = sjson.Set(rpcs, name, fmt.Sprintf("%s/ext/bc/%s/rpc", url, blockchainID))
		}
		return true
	})

	out := "{}"
	out, _ = sjson.Set(out, "nodeID", getNodeID.Get("result.nodeID").String())
	out, _ = sjson.Set(out, "networkID", getNetworkID.Get("result.networkID").Int())
	out, _ = sjson.Set(out, "networkName", getNetworkName.Get("result.networkName").String())
	if getUptime.Get("result").String() != "" {
		out, _ = sjson.SetRaw(out, "uptime", getUptime.Get("result").String())
	}
	out, _ = sjson.SetRaw(out, "getNodeVersion", getNodeVersion.Get("result").String())
	out, _ = sjson.SetRaw(out, "getVMs", getVMs.Get("result").String())
	out, _ = sjson.SetRaw(out, "subnets", getSubnets.Get("result.subnets").String())
	out, _ = sjson.SetRaw(out, "stakingAssetIDs", stakingAssetIDs)
	out, _ = sjson.SetRaw(out, "blockchains", getBlockchains.Get("result.blockchains").String())
	out, _ = sjson.SetRaw(out, "aliases", aliases)
	out, _ = sjson.SetRaw(out, "rpcs", rpcs)

	result := gjson.Parse(out)
	return &result, nil
}
