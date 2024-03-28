package walletcmd

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/wallet/chain/p"
	"github.com/ava-labs/avalanchego/wallet/chain/x"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance addr",
		Short: "Balance of AVAX",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			uri := viper.GetString("node-url")
			addrStr := args[0]
			balanceP, err := fetchBalanceP(uri, addrStr)
			cobra.CheckErr(err)
			balanceX, err := fetchBalanceX(uri, addrStr)
			cobra.CheckErr(err)
			fmt.Printf("Address: %s\n", addrStr)
			fmt.Printf("Balance P: %d\n", balanceP)
			fmt.Printf("Balance X: %d\n", balanceX)
		},
	}
	return cmd
}

func fetchBalanceP(uri string, addrStr string) (uint64, error) {
	addr, err := address.ParseToID(addrStr)
	if err != nil {
		return 0, nil
	}

	addresses := set.Of(addr)

	ctx := context.Background()

	state, err := primary.FetchState(ctx, uri, addresses)
	if err != nil {
		return 0, nil
	}

	pUTXOs := common.NewChainUTXOs(constants.PlatformChainID, state.UTXOs)
	pBackend := p.NewBackend(state.PCTX, pUTXOs, nil)
	pBuilder := p.NewBuilder(addresses, pBackend)

	currentBalances, err := pBuilder.GetBalance()
	if err != nil {
		return 0, err
	}

	avaxID := state.PCTX.AVAXAssetID()
	avaxBalance := currentBalances[avaxID]
	return avaxBalance, nil
}

func fetchBalanceX(uri string, addrStr string) (uint64, error) {
	addr, err := address.ParseToID(addrStr)
	if err != nil {
		return 0, nil
	}

	addresses := set.Of(addr)

	ctx := context.Background()

	state, err := primary.FetchState(ctx, uri, addresses)
	if err != nil {
		return 0, nil
	}

	xUTXOs := common.NewChainUTXOs(constants.PlatformChainID, state.UTXOs)
	xBackend := x.NewBackend(state.XCTX, xUTXOs)
	xBuilder := x.NewBuilder(addresses, xBackend)

	currentBalances, err := xBuilder.GetFTBalance()
	if err != nil {
		return 0, err
	}

	avaxID := state.XCTX.AVAXAssetID()
	avaxBalance := currentBalances[avaxID]
	return avaxBalance, nil
}
