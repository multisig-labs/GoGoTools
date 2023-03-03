package walletcmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/multisig-labs/gogotools/pkg/application"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var app *application.GoGoTools
var pkStr string
var keyFactory = new(secp256k1.Factory)

var (
	ErrInvalidType = errors.New("invalid type")
	ErrCantSpend   = errors.New("can't spend")
)

func NewCmd(injectedApp *application.GoGoTools) *cobra.Command {
	app = injectedApp

	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "Issue various txs to create subnets and blockchains",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	// Default key is the one used by Avalanche 'local' and 'custom' (ANR) networks
	// PrivateKey-ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN => P-local18jma8ppw3nhx5r4ap8clazz0dps7rv5u00z96u
	// PrivateKey-ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN => P-custom18jma8ppw3nhx5r4ap8clazz0dps7rv5u9xde7p
	// 56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027 => 0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC
	cmd.PersistentFlags().StringVar(&pkStr, "pk", "PrivateKey-ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN", "Private key")
	viper.BindPFlag("pk", cmd.PersistentFlags().Lookup("pk"))

	cmd.AddCommand(newCreateSubnetCmd())
	cmd.AddCommand(newCreateChainCmd())
	cmd.AddCommand(newAddValidatorCmd())

	return cmd
}

func decodePrivateKey(enc string) (*secp256k1.PrivateKey, error) {
	rawPk := strings.Replace(enc, "PrivateKey-", "", 1)
	skBytes, err := cb58.Decode(rawPk)
	if err != nil {
		return nil, fmt.Errorf("unable to decode private key: %w", err)
	}
	privKey, err := keyFactory.ToPrivateKey(skBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to decode private key: %w", err)
	}
	return privKey, nil
}
