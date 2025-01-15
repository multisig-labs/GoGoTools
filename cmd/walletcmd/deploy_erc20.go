package walletcmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ava-labs/coreth/accounts/abi/bind"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/ethclient"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/multisig-labs/gogotools/pkg/contracts/erc20"
	"github.com/multisig-labs/gogotools/pkg/hd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tyler-smith/go-bip39"
)

// TODO hacking this in for now. Figure out better way to handle private keys / mnemonics for X/P and C chains.
func newDeployERC20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy-erc20 [name] [symbol]",
		Short: "Deploy an OpenZeppelin ERC20",
		Long:  `MNEMONIC must be set in env, will be deployer / owner of the token`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenName := args[0]
			tokenSym := args[1]

			baseURL := viper.GetString("node-url")
			client, err := ethclient.Dial(baseURL + "/ext/bc/C/rpc")
			cobra.CheckErr(err)

			mnemonic := os.Getenv("MNEMONIC")
			if ok := bip39.IsMnemonicValid(mnemonic); !ok {
				return fmt.Errorf("invaid mnemonic")
			}

			hdkeys, err := hd.DeriveHDKeys(mnemonic, hd.EthDerivationPath, 1)
			if err != nil {
				return fmt.Errorf("error deriving keys: %s", err)
			}

			k := hdkeys[0]
			ownerAddr := ethcrypto.PubkeyToAddress(k.PK.PublicKey)

			evmChainID, err := client.ChainID(context.Background())
			cobra.CheckErr(err)
			auth, err := bind.NewKeyedTransactorWithChainID(k.PK, evmChainID)
			cobra.CheckErr(err)
			erc20Addr, tx, _, err := erc20.DeployErc20(auth, client, tokenName, tokenSym, ownerAddr)
			cobra.CheckErr(err)
			waitForEVMTx(client, tx)

			cobra.CheckErr(err)
			fmt.Println(erc20Addr)
			return nil
		},
	}
	return cmd
}

func waitForEVMTx(client ethclient.Client, tx *types.Transaction) {
	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()
	maxAttempts := 5
	attempt := 0

	for {
		select {
		case <-timeout.C:
			panic("Transaction confirmation timed out!")
		default:
			attempt++
			if attempt > maxAttempts {
				panic("Transaction confirmation max attempts!")
			}
			receipt, _ := client.TransactionReceipt(context.Background(), tx.Hash())
			if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
				// tx confirmed
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}
