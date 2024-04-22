package utilscmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/multisig-labs/gogotools/pkg/hd"

	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
)

func newMnemonicKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mnemonic-key [mnemonic] [hrp] [chain] [index]",
		Short: "Show address and priv key for a BIP39 mnemonic index",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if ok := bip39.IsMnemonicValid(args[0]); !ok {
				return fmt.Errorf("invaid mnemonic")
			}

			hrp := "avax"
			if len(args) > 1 {
				hrp = args[1]
			}

			chain := "P"
			if len(args) > 2 {
				chain = args[2]
			}

			var dpath accounts.DerivationPath
			if chain == "C" {
				dpath = hd.EthDerivationPath
			} else {
				dpath = hd.AvaDerivationPath
			}

			idx := 0
			if len(args) > 2 {
				idx, _ = strconv.Atoi(args[3])
			}

			hdkeys, err := hd.DeriveHDKeys(args[0], dpath, idx+1)
			if err != nil {
				return fmt.Errorf("error deriving keys: %s", err)
			}

			k := hdkeys[idx]

			var address string
			var privkey string

			if chain == "C" {
				address = k.EthAddr()
				privkey = k.EthPrivKey()
			} else {
				address = k.AvaAddr(chain, hrp)
				privkey = k.AvaPrivKey()
			}

			keyPair := struct {
				Address string `json:"addr"`
				PrivKey string `json:"pk"`
			}{
				Address: address,
				PrivKey: privkey,
			}

			jsonOutput, err := json.MarshalIndent(keyPair, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(jsonOutput))

			return nil
		},
	}
	return cmd
}
