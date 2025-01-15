package utilscmd

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
)

func newInspectTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect-tx-p [tx]",
		Short: "Output JSON for a P-chain tx encoded in either hex or base64",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			txStr := args[0]
			var txb []byte
			var err error
			if strings.HasPrefix(txStr, "0x") {
				txb, err = hexutil.Decode(txStr)
				cobra.CheckErr(err)
			} else {
				txb, err = b64.StdEncoding.DecodeString(txStr)
				cobra.CheckErr(err)
			}

			tx := &txs.Tx{}
			_, err = txs.Codec.Unmarshal(txb, tx)
			cobra.CheckErr(err)
			// fmt.Printf("%+v", tx)
			js, err := json.Marshal(tx)
			cobra.CheckErr(err)
			fmt.Printf("%s\n", js)
		},
	}
	return cmd
}
