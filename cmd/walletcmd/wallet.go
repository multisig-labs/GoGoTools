package walletcmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/multisig-labs/gogotools/pkg/application"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var app *application.GoGoTools
var pkStr string
var keyFactory = new(crypto.FactorySECP256K1R)

var (
	ErrInvalidType = errors.New("invalid type")
	ErrCantSpend   = errors.New("can't spend")
)

func NewCmd(injectedApp *application.GoGoTools) *cobra.Command {
	app = injectedApp

	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "Wallet",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	cmd.PersistentFlags().StringVar(&pkStr, "pk", "PrivateKey-ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN", "Private key")
	viper.BindPFlag("pk", cmd.PersistentFlags().Lookup("pk"))

	cmd.AddCommand(newCreateSubnetCmd())
	cmd.AddCommand(newCreateChainCmd())

	return cmd
}

func decodePrivateKey(enc string) (*crypto.PrivateKeySECP256K1R, error) {
	rawPk := strings.Replace(enc, "PrivateKey-", "", 1)
	skBytes, err := cb58.Decode(rawPk)
	if err != nil {
		return nil, fmt.Errorf("unable to decode private key: %w", err)
	}
	rpk, err := keyFactory.ToPrivateKey(skBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to decode private key: %w", err)
	}
	privKey, ok := rpk.(*crypto.PrivateKeySECP256K1R)
	if !ok {
		return nil, ErrInvalidType
	}
	return privKey, nil
}
