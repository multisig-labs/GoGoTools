package utilscmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/vms/platformvm/block"

	"github.com/spf13/cobra"
)

// TODO Need to fix this to work with avalanchego 1.11.x

func newInspectBlockCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspectblk [block]",
		Short: "Decode a hex encoded block",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := genMainnetCtx()
			var blkBytes []byte
			var err error
			if strings.HasPrefix(args[0], "0x") {
				blkBytes, err = formatting.Decode(formatting.Hex, args[0])
			} else {
				blkBytes, err = base64.StdEncoding.DecodeString(args[0])
			}
			cobra.CheckErr(err)

			blk, js, err := decodeBlock(ctx, blkBytes)
			cobra.CheckErr(err)
			fmt.Printf("%s\n", js)
			fmt.Printf("%v\n", blk)
			fmt.Printf("%x\n", blk.Txs()[0].Unsigned.Bytes())
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

func decodeBlock(ctx *snow.Context, b []byte) (block.Block, string, error) {
	decoded := decodeProposerBlock(b)

	blk, js, err := decodeInnerBlock(ctx, decoded)
	if err != nil {
		return blk, "", err
	}
	return blk, string(js), nil
}

// Tries to decode as proposal block (post-Banff) if it fails just return the original bytes
func decodeProposerBlock(b []byte) []byte {
	innerBlk, err := block.Parse(b)
	if err != nil {
		return b
	}
	return innerBlk.Block()
}

func decodeInnerBlock(ctx *snow.Context, b []byte) (block.Block, string, error) {
	res, err := block.Parse(block.GenesisCodec, b)
	if err != nil {
		return res, "", fmt.Errorf("blocks.Parse error: %w", err)
	}

	res.InitCtx(ctx)
	j, err := json.Marshal(res)
	if err != nil {
		return res, "", fmt.Errorf("json.Marshal error: %w", err)
	}
	return res, string(j), nil
}

// Simple context so that Marshal works
func genMainnetCtx() *snow.Context {
	pChainID, _ := ids.FromString("11111111111111111111111111111111LpoYY")
	xChainID, _ := ids.FromString("2oYMBNV4eNHyqk2fjjV5nVQLDbtmNJzq5s3qs3Lo6ftnC6FByM")
	cChainID, _ := ids.FromString("2q9e4r6Mu3U68nU1fYjgbR6JvwrRx36CohpAX5UQxse55x1Q5")
	avaxAssetID, _ := ids.FromString("FvwEAhmxKfeiG8SnEvq42hc6whRyY3EFYAvebMqDNDGCgxN5Z")
	lookup := ids.NewAliaser()
	lookup.Alias(xChainID, "X")
	lookup.Alias(cChainID, "C")
	lookup.Alias(pChainID, "P")
	c := &snow.Context{
		NetworkID:   1,
		SubnetID:    [32]byte{},
		ChainID:     [32]byte{},
		NodeID:      [20]byte{},
		XChainID:    xChainID,
		CChainID:    cChainID,
		AVAXAssetID: avaxAssetID,
		Lock:        sync.RWMutex{},
		BCLookup:    lookup,
	}
	return c
}
