package warpcmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	warppayload "github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	subnetEvmWarp "github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	warpmessages "github.com/ava-labs/subnet-evm/warp/messages"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/multisig-labs/gogotools/pkg/application"
	"github.com/spf13/cobra"
)

var app *application.GoGoTools

func NewCmd(injectedApp *application.GoGoTools) *cobra.Command {
	app = injectedApp

	cmd := &cobra.Command{
		Use:          "warp",
		Short:        "Warp utilities",
		Long:         ``,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmd.AddCommand(newParseWarpMsgCmd())
	cmd.AddCommand(newGetWarpMsgCmd())
	cmd.AddCommand(newConstructUptimeMsgCmd())
	return cmd
}

func newGetWarpMsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [url] [txid]",
		Short: "Get a warp message from rpc at [url] and parse msg from tx logs",
		Long:  ``,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			txid := common.HexToHash(args[1])
			c, err := ethclient.Dial(args[0])
			cobra.CheckErr(err)
			receipt, err := c.TransactionReceipt(context.Background(), txid)
			cobra.CheckErr(err)
			uwm, err := warpMessageFromLogs(receipt.Logs)
			cobra.CheckErr(err)
			wm := &warp.Message{UnsignedMessage: *uwm}
			fmt.Printf("\n%+v\n", wm)
		},
	}
	return cmd
}

func newConstructUptimeMsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uptime [network] [blockchainID] [validationID] [uptimeSeconds]",
		Short: "Construct a warp Uptime message",
		Long:  ``,
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			networkID, err := constants.NetworkID(args[0])
			cobra.CheckErr(err)
			blockchainID, err := ids.FromString(args[1])
			cobra.CheckErr(err)
			validationID, err := ids.FromString(args[2])
			cobra.CheckErr(err)
			uptime, err := strconv.ParseUint(args[3], 10, 64)
			uptimePayload, err := warpmessages.NewValidatorUptime(validationID, uptime)
			cobra.CheckErr(err)
			addressedCall, err := warppayload.NewAddressedCall(nil, uptimePayload.Bytes())
			cobra.CheckErr(err)
			uptimeProofUnsignedMessage, err := warp.NewUnsignedMessage(
				networkID,
				blockchainID,
				addressedCall.Bytes(),
			)
			fmt.Printf("%s\n", hexutil.Encode(uptimeProofUnsignedMessage.Bytes()))
		},
	}
	return cmd
}

func newParseWarpMsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse [warpMsg]",
		Short: "Parse a warp message",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			m, err := parseWarpMessage(args[0])
			if err != nil && strings.Contains(err.Error(), "insufficient length") {
				um, err := parseUnsignedWarpMessage(args[0])
				cobra.CheckErr(err)
				m = &warp.Message{UnsignedMessage: *um}
			} else {
				cobra.CheckErr(err)
			}

			fmt.Printf("\n%+v\n", m)
		},
	}
	return cmd
}

func parseWarpMessage(msgHex string) (*warp.Message, error) {
	msgBytes, err := hexutil.Decode(msgHex)
	if err != nil {
		return nil, fmt.Errorf("error decoding warp message: %w", err)
	}
	msg, err := warp.ParseMessage(msgBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing signed message: %w", err)
	}
	return msg, nil
}

func parseUnsignedWarpMessage(msgHex string) (*warp.UnsignedMessage, error) {
	msgBytes, err := hexutil.Decode(msgHex)
	if err != nil {
		return nil, fmt.Errorf("error decoding warp message: %w", err)
	}
	msg, err := warp.ParseUnsignedMessage(msgBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing unsigned message: %w", err)
	}
	return msg, nil
}

// Returns the first log in 'logs' that is successfully parsed by 'parser'
func getEventFromLogs[T any](logs []*types.Log, parser func(log types.Log) (T, error)) (T, error) {
	cumErrMsg := ""
	for i, log := range logs {
		event, err := parser(*log)
		if err == nil {
			return event, nil
		}
		if cumErrMsg != "" {
			cumErrMsg += "; "
		}
		cumErrMsg += fmt.Sprintf("log %d -> %s", i, err.Error())
	}
	return *new(T), fmt.Errorf("failed to find %T event in receipt logs: [%s]", *new(T), cumErrMsg)
}

func warpMessageFromLogs(logs []*types.Log) (*avalancheWarp.UnsignedMessage, error) {
	return getEventFromLogs(logs, parseSendWarpMessage)
}

func parseSendWarpMessage(log types.Log) (*avalancheWarp.UnsignedMessage, error) {
	return subnetEvmWarp.UnpackSendWarpEventDataToMessage(log.Data)
}
