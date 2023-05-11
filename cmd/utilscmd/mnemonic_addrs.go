package utilscmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/hd"

	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
)

func newMnemonicAddrsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mnemonic-addrs [mnemonic] [hrp]",
		Short: "Show public addresses for a BIP39 mnemonic",
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

			fmt.Printf("=== C-Chain [%s] ===\n", hrp)

			hdkeys, err := hd.DeriveHDKeys(args[0], hd.EthDerivationPath, 10)
			if err != nil {
				return fmt.Errorf("error deriving keys: %s", err)
			}

			for _, k := range hdkeys {
				fmt.Printf("%s %s\n", k.EthAddr(), k.Path)
			}

			fmt.Printf("=== X-Chain [%s] ===\n", hrp)

			hdkeys, err = hd.DeriveHDKeys(args[0], hd.AvaDerivationPath, 10)
			if err != nil {
				return fmt.Errorf("error deriving keys: %s", err)
			}

			for _, k := range hdkeys {
				fmt.Printf("%s %s\n", k.AvaAddr("X", hrp), k.Path)
			}

			fmt.Printf("=== P-Chain [%s] ===\n", hrp)

			hdkeys, err = hd.DeriveHDKeys(args[0], hd.AvaDerivationPath, 10)
			if err != nil {
				return fmt.Errorf("error deriving keys: %s", err)
			}

			for _, k := range hdkeys {
				fmt.Printf("%s %s\n", k.AvaAddr("P", hrp), k.Path)
			}

			return nil
		},
	}
	return cmd
}
