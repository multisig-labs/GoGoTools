package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/wallet/chain/p/builder"
	"github.com/ava-labs/avalanchego/wallet/chain/p/wallet"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary/common"
	"github.com/ava-labs/subnet-evm/ethclient"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/multisig-labs/gogotools/pkg/utils"

	"github.com/jxskiss/mcli"
)

func balanceAddressCmd() {
	args := struct {
		Address string `cli:"#R, address, Address"`
		Raw     bool   `cli:"--raw, raw, Output raw balance with no conversions"`
		URLFlags
	}{}
	mcli.MustParse(&args)

	if strings.HasPrefix(args.Address, "0x") {
		balance, err := getEthBalance(context.Background(), args.EthUrl, ethcommon.HexToAddress(args.Address))
		checkErr(err)

		if args.Raw {
			fmt.Println(balance.String())
		} else {
			fbalance := new(big.Float).SetInt(balance)
			eth := new(big.Float).SetInt(big.NewInt(1e18))
			amt := new(big.Float).Quo(fbalance, eth)
			fmt.Printf("%.18f ETH\n", amt)
		}
	} else {
		addr, err := address.ParseToID(args.Address)
		checkErr(err)

		balance, err := getPBalance(context.Background(), args.AvaUrl, addr)
		checkErr(err)

		if args.Raw {
			fmt.Println(balance)
		} else {
			fmt.Printf("%.9f AVAX\n", float64(balance)/1e9)
		}
	}
}

func balancePKCmd() {
	args := struct {
		PK string `cli:"#R, pk, Show P-Chain and C-Chain balances for a private key" env:"PRIVATE_KEY"`
		URLFlags
	}{}
	mcli.MustParse(&args)

	avaAddrStr, ethAddrStr, err := utils.ParsePrivateKeyToAddresses(args.PK, "avax")
	checkErr(err)

	ethBalance, err := getEthBalance(context.Background(), args.EthUrl, ethcommon.HexToAddress(ethAddrStr))
	checkErr(err)

	avaAddr, err := address.ParseToID(avaAddrStr)
	checkErr(err)

	avaBalance, err := getPBalance(context.Background(), args.AvaUrl, avaAddr)
	checkErr(err)

	fbalance := new(big.Float).SetInt(ethBalance)
	eth := new(big.Float).SetInt(big.NewInt(1e18))
	amt := new(big.Float).Quo(fbalance, eth)
	fmt.Printf("%.9f ETH   %s\n", amt, ethAddrStr)
	fmt.Printf("%.9f AVAX  %s\n", float64(avaBalance)/1e9, avaAddrStr)
}

func getEthBalance(ctx context.Context, url string, address ethcommon.Address) (*big.Int, error) {
	c, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	balance, err := c.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func getPBalance(ctx context.Context, url string, address ids.ShortID) (uint64, error) {
	addresses := set.Of(address)

	state, err := primary.FetchState(ctx, url, addresses)
	if err != nil {
		return 0, err
	}

	pUTXOs := common.NewChainUTXOs(constants.PlatformChainID, state.UTXOs)
	pBackend := wallet.NewBackend(state.PCTX, pUTXOs, nil)
	pBuilder := builder.New(addresses, state.PCTX, pBackend)

	currentBalances, err := pBuilder.GetBalance()
	if err != nil {
		return 0, err
	}

	avaxID := state.PCTX.AVAXAssetID
	avaxBalance := currentBalances[avaxID]

	return avaxBalance, nil
}
