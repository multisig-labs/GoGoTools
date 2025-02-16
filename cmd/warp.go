package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	subnetEvmWarp "github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	warpmessages "github.com/ava-labs/subnet-evm/warp/messages"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jxskiss/mcli"
)

func getWarpMsgCmd() {
	args := struct {
		TxID string `cli:"txid, Transaction ID"`

		URLFlags
	}{}
	mcli.MustParse(&args)

	txid := common.HexToHash(args.TxID)

	c, err := ethclient.Dial(args.EthUrl)
	checkErr(err)
	receipt, err := c.TransactionReceipt(context.Background(), txid)
	checkErr(err)
	uwm, err := warpMessageFromLogs(receipt.Logs)
	checkErr(err)
	wm := &warp.Message{UnsignedMessage: *uwm}
	fmt.Printf("\n%+v\n", wm)
}

func constructUptimeMsgCmd() {
	args := struct {
		Network       string `cli:"network, Network (mainnet, fuji, etc)"`
		BlockchainID  string `cli:"blockchainID, Blockchain ID"`
		ValidationID  string `cli:"validationID, Validation ID"`
		UptimeSeconds string `cli:"uptimeSeconds, Uptime Seconds"`
	}{}
	mcli.MustParse(&args)

	networkID, err := constants.NetworkID(args.Network)
	checkErr(err)
	blockchainID, err := ids.FromString(args.BlockchainID)
	checkErr(err)
	validationID, err := ids.FromString(args.ValidationID)
	checkErr(err)
	uptimeSeconds, err := strconv.ParseUint(args.UptimeSeconds, 10, 64)
	checkErr(err)
	uptimePayload, err := warpmessages.NewValidatorUptime(validationID, uptimeSeconds)
	checkErr(err)
	addressedCall, err := payload.NewAddressedCall(nil, uptimePayload.Bytes())
	checkErr(err)
	uptimeProofUnsignedMessage, err := warp.NewUnsignedMessage(
		networkID,
		blockchainID,
		addressedCall.Bytes(),
	)
	checkErr(err)
	fmt.Printf("\n%s\n", hexutil.Encode(uptimeProofUnsignedMessage.Bytes()))
}

func parseWarpMsgCmd() {
	args := struct {
		WarpMsg string `cli:"warpMsg, Warp Message"`
	}{}
	mcli.MustParse(&args)

	m, err := parseWarpMessage(args.WarpMsg)
	if err != nil && strings.Contains(err.Error(), "insufficient length") {
		um, err := parseUnsignedWarpMessage(args.WarpMsg)
		checkErr(err)
		m = &warp.Message{UnsignedMessage: *um}
	} else {
		checkErr(err)
	}

	payloadType, payload, err := parsePayload(m.Payload)
	checkErr(err)

	fmt.Printf("\n%+v\n\nPayload (%s): %s\n", m, payloadType, payload)
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

func warpMessageFromLogs(logs []*types.Log) (*warp.UnsignedMessage, error) {
	return getEventFromLogs(logs, parseSendWarpMessage)
}

func parseSendWarpMessage(log types.Log) (*warp.UnsignedMessage, error) {
	return subnetEvmWarp.UnpackSendWarpEventDataToMessage(log.Data)
}

// Returns json of the decoded payload
func parsePayload(msg []byte) (string, []byte, error) {
	addressedCall, err := payload.ParseAddressedCall(msg)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse addressed call: %w", err)
	}

	payloadIntf, err := message.Parse(addressedCall.Payload)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	var out []byte

	switch payload := payloadIntf.(type) {
	case *message.RegisterL1Validator:
		out, err = json.Marshal(payload)
		if err != nil {
			return "", nil, err
		}
	case *message.L1ValidatorRegistration:
		out, err = json.Marshal(payload)
		if err != nil {
			return "", nil, err
		}
	default:
		return "", nil, fmt.Errorf("unknown type: %T", payload)
	}

	return fmt.Sprintf("%T", payloadIntf), out, nil
}
