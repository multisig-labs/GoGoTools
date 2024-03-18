package nodecmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func newValidatorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators subnetID",
		Short: "Get current validators for a subnet (leave empty for Primary Network)",
		Long:  ``,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			subnetID := "11111111111111111111111111111111LpoYY"
			if len(args) > 0 {
				subnetID = args[0]
			}
			result, err := getValidators(subnetID)
			cobra.CheckErr(err)
			fmt.Println(result.String())
		},
	}
	return cmd
}

func getValidators(subnetID string) (*gjson.Result, error) {
	uri := viper.GetString("node-url")
	urlP := fmt.Sprintf("%s/ext/bc/P", uri)

	currValdrs, err := utils.FetchRPCGJSON(urlP, "platform.getCurrentValidators", fmt.Sprintf(`{"subnetID":"%s"}`, subnetID))
	return currValdrs, err
}
