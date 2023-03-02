package utilscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tyler-smith/go-bip39"
)

func newMnemonicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mnemonic",
		Short: "Generate a BIP39 mnemonic",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			entropy, _ := bip39.NewEntropy(256)
			phrase, _ := bip39.NewMnemonic(entropy)

			fmt.Println(phrase)
			return nil
		},
	}
	return cmd
}
