package nodecmd

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

func newExplorerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explorer chain-name",
		Short: "Launch a browser to a blockchain explorer pointed at chain-name",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			s, err := getInfo()
			cobra.CheckErr(err)
			rpc := gjson.Parse(s).Get(fmt.Sprintf("rpcs.%s", args[0])).String()
			if rpc == "" {
				app.Log.Fatalf("Unable to find chain-name %s in the 'rpcs' key of 'ggt node info'", args[0])
			}
			url := fmt.Sprintf("http://expedition.fly.dev?rpcUrl=%s", rpc)
			app.Log.Infof("Opening %s", url)
			browser.OpenURL(url)
		},
	}
	return cmd
}
