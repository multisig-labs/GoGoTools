package utilscmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/hd"

	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
)

func newMnemonicKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mnemonic-keys [mnemonic] [hrp]",
		Short: "Show keys and addresses for a BIP39 mnemonic",
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

			fmt.Println("=== C-Chain ===")

			hdkeys, err := hd.DeriveHDKeys(args[0], hd.EthDerivationPath, 10)
			if err != nil {
				return fmt.Errorf("error deriving keys: %s", err)
			}

			for _, k := range hdkeys {
				fmt.Printf("%s %s %s\n", k.EthAddr(), k.EthPrivKey(), k.Path)
			}

			fmt.Printf("\n=== P-Chain [%s] ===\n", hrp)

			hdkeys, err = hd.DeriveHDKeys(args[0], hd.AvaDerivationPath, 10)
			if err != nil {
				return fmt.Errorf("error deriving keys: %s", err)
			}

			for _, k := range hdkeys {
				fmt.Printf("%s %s %s\n", k.AvaAddr("P", hrp), k.AvaPrivKey(), k.Path)
			}

			fmt.Printf("\n=== P-Chain [%s] (using eth derivation path)===\n", hrp)

			hdkeys, err = hd.DeriveHDKeys(args[0], hd.EthDerivationPath, 10)
			if err != nil {
				return fmt.Errorf("error deriving keys: %s", err)
			}

			for _, k := range hdkeys {
				fmt.Printf("%s %s %s\n", k.AvaAddr("P", hrp), k.EthPrivKey(), k.Path)
			}

			return nil
		},
	}
	return cmd
}
