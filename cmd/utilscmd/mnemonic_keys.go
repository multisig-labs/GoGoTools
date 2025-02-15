package utilscmd

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/hd"

	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
)

func newMnemonicKeysCmd() *cobra.Command {
	var numKeys int

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

			fmt.Println("=== Ethereum Derivation Path ===")
			hdkeys, err := hd.DeriveHDKeys(args[0], hd.EthDerivationPath, numKeys)
			if err != nil {
				return fmt.Errorf("error deriving keys: %s", err)
			}

			fmt.Printf("%-16s %42s %45s %64s %61s\n", "Path", "EVM Addr", "Ava Addr", "EVM Private Key", "Ava Private Key")
			for _, k := range hdkeys {
				fmt.Printf("%-16s %42s %45s %64s %61s\n",
					k.Path,
					k.EthAddr(),
					k.AvaAddr("P", hrp),
					k.EthPrivKey(),
					k.AvaPrivKey(),
				)
			}

			fmt.Println("=== Avalanche Derivation Path ===")
			hdkeys, err = hd.DeriveHDKeys(args[0], hd.AvaDerivationPath, numKeys)
			if err != nil {
				return fmt.Errorf("error deriving keys: %s", err)
			}

			for _, k := range hdkeys {
				fmt.Printf("%s %s %s %s %s\n",
					k.Path,
					k.EthAddr(),
					k.AvaAddr("P", hrp),
					k.EthPrivKey(),
					k.AvaPrivKey(),
				)
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(&numKeys, "num-keys", "n", 10, "number of keys to generate")
	return cmd
}
