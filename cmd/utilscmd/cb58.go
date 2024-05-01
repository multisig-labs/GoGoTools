package utilscmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/cb58"

	"github.com/spf13/cobra"
)

func newCB58DecodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cb58decode [value]",
		Short: "Decode a CB58 encoded string",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := cb58.Decode(args[0])
			cobra.CheckErr(err)
			fmt.Printf("%x\n", b)
			// _, _, b2, err := address.Parse("P-avax1gfpj30csekhwmf4mqkncelus5zl2ztqzvv7aww")
			// cobra.CheckErr(err)
			// fmt.Printf("%x\n", b2)

			// factory := secp256k1.Factory{}
			// pk, err := factory.ToPrivateKey(b)
			// cobra.CheckErr(err)
			// fmt.Printf("PrivKey: %x\n", pk)

			// a := pk.PublicKey().Bytes()
			// fmt.Printf("Serialized compressed pub key bytes: %x\n", a)

			// addr, err := address.FormatBech32("P-avax1", id.Bytes())
			// cobra.CheckErr(err)
			// fmt.Printf("Addr: %s\n", addr)

			// fmt.Printf("%x\n", b)
			return nil
		},
	}
	return cmd
}

func newCB58DecodeNodeIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cb58decodenodeid [value]",
		Short: "Decode a CB58 encoded NodeID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, _ := strings.CutPrefix(args[0], "NodeID-")
			b, err := cb58.Decode(s)
			cobra.CheckErr(err)
			fmt.Printf("%x\n", b)
			return nil
		},
	}
	return cmd

}

func newCB58DecodeSigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cb58decodesig [value]",
		Short: "Decode a CB58 encoded signature into r,s,v",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := cb58.Decode(args[0])
			cobra.CheckErr(err)
			sig := struct {
				R string `json:"r"`
				S string `json:"s"`
				V string `json:"v"`
			}{}
			sig.R = fmt.Sprintf("0x%x", b[0:32])
			sig.S = fmt.Sprintf("0x%x", b[32:64])
			sig.V = fmt.Sprintf("0x%x", b[64:])
			j, _ := json.Marshal(sig)
			fmt.Printf("%s\n", j)
			return nil
		},
	}
	return cmd
}
