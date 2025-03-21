package main

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
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
		Address string `cli:"#R, address, P-Chain or C-Chain address" env:"ETH_FROM"`
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
		_, addrBytes, err := address.ParseBech32(strings.TrimPrefix(args.Address, "P-"))
		checkErr(err)
		addr, err := ids.ToShortID(addrBytes)
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

func crossChainTransferCmd() {
	args := struct {
		Amount string `cli:"#R, amount, Amount to transfer in MilliAvax, 1000 = 1 AVAX" env:"AMOUNT"`
		PK     string `cli:"#R, pk, Transfer from C-Chain to P-Chain address for a given private key" env:"PRIVATE_KEY"`
		URLFlags
	}{}
	mcli.MustParse(&args)

	avaKey, _, err := utils.ParsePrivateKey(args.PK)
	checkErr(err)

	kc := secp256k1fx.NewKeychain(avaKey)
	ctx := context.Background()
	exportWallet, err := primary.MakeWallet(ctx,
		args.AvaUrl,
		kc,
		kc,
		primary.WalletConfig{
			SubnetIDs:     []ids.ID{},
			ValidationIDs: []ids.ID{},
		},
	)
	checkErr(err)

	cChainID := exportWallet.C().Builder().Context().BlockchainID

	owner := secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs: []ids.ShortID{
			avaKey.Address(),
		},
	}

	PlatformChainID := ids.Empty
	argsAmount, err := strconv.ParseUint(args.Amount, 10, 64)
	checkErr(err)

	fmt.Printf("starting export... \n")
	amount := argsAmount * units.MilliAvax
	_, err = exportWallet.C().IssueExportTx(PlatformChainID,
		[]*secp256k1fx.TransferOutput{{
			Amt:          amount,
			OutputOwners: owner,
		}},
	)
	checkErr(err)
	fmt.Printf("finished export. \n")

	newWallet, err := primary.MakeWallet(ctx,
		args.AvaUrl,
		kc,
		kc,
		primary.WalletConfig{
			SubnetIDs:     []ids.ID{},
			ValidationIDs: []ids.ID{},
		},
	)
	checkErr(err)

	fmt.Printf("staring import... \n")
	_, err = newWallet.P().IssueImportTx(cChainID, &owner)
	checkErr(err)
	fmt.Printf("finished import. \n")
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
