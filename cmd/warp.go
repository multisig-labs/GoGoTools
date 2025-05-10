package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	subnetevmwarp "github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	subnetevmmessages "github.com/ava-labs/subnet-evm/warp/messages"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jxskiss/mcli"
	"github.com/multisig-labs/gogotools/pkg/sigagg"
	"github.com/multisig-labs/gogotools/pkg/utils"
)

func getWarpMsgCmd() {
	args := struct {
		TxID string `cli:"#R, txid, Transaction ID"`
		Hex  bool   `cli:"--hex, Output as hex"`
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

	payloadType, payload, err := parsePayload(uwm.Payload)
	checkErr(err)

	if args.Hex {
		fmt.Println(utils.BytesToHex(uwm.Bytes()))
	} else {
		fmt.Printf("\n%+v\n\nPayload (%s): %s\n", uwm, payloadType, payload)
	}
}

func constructL1ValidatorRegistrationMsg(networkID uint32, validationID ids.ID, registered bool) (*warp.UnsignedMessage, error) {
	addressedCallPayload, err := message.NewL1ValidatorRegistration(validationID, registered)
	if err != nil {
		return nil, err
	}
	addressedCall, err := payload.NewAddressedCall(
		nil,
		addressedCallPayload.Bytes(),
	)
	if err != nil {
		return nil, err
	}
	unsignedMessage, err := warp.NewUnsignedMessage(
		networkID,
		ids.Empty,
		addressedCall.Bytes(),
	)
	if err != nil {
		return nil, err
	}
	return unsignedMessage, nil
}

func constructL1ValidatorRegistrationMsgCmd() {
	args := struct {
		Network      string `cli:"#R, -n,   --network, Network (mainnet, fuji, etc)"`
		ValidationID string `cli:"#R, -v,   --validation, Validation ID"`
		Registered   bool   `cli:"#R, -r,   --registered, Is validator registered"`
	}{}
	mcli.MustParse(&args)

	networkID, err := constants.NetworkID(args.Network)
	checkErr(err)
	validationID, err := ids.FromString(args.ValidationID)
	checkErr(err)

	unsignedMsg, err := constructL1ValidatorRegistrationMsg(networkID, validationID, args.Registered)
	checkErr(err)

	fmt.Printf("0x%x\n", unsignedMsg.Bytes())
}

func constructL1WeightMsgCmd() {
	args := struct {
		Network      string `cli:"#R, -n,   --network, Network (mainnet, fuji, etc)"`
		ValidationID string `cli:"#R, -v,   --validation, Validation ID"`
		Weight       uint64 `cli:"#R, -r,   --weight"`
		Nonce        uint64 `cli:"#R, --nonce"`
		BlockchainID string `cli:"--blockchain, Source Blockchain ID" default:"11111111111111111111111111111111LpoYY"`
		Address      string `cli:"--address, Address for AddressedCall"`
	}{}
	mcli.MustParse(&args)

	networkID, err := constants.NetworkID(args.Network)
	checkErr(err)
	validationID, err := ids.FromString(args.ValidationID)
	checkErr(err)
	blockchainID, err := ids.FromString(args.BlockchainID)
	checkErr(err)

	addressedCallPayload, err := message.NewL1ValidatorWeight(
		validationID,
		args.Nonce,
		args.Weight,
	)
	checkErr(err)

	var addr []byte
	if args.Address != "" {
		addr = utils.HexToBytes(args.Address)
	} else {
		addr = nil
	}
	addressedCall, err := payload.NewAddressedCall(
		addr,
		addressedCallPayload.Bytes(),
	)
	checkErr(err)
	unsignedMessage, err := warp.NewUnsignedMessage(
		networkID,
		blockchainID,
		addressedCall.Bytes(),
	)
	checkErr(err)

	fmt.Printf("0x%x\n", unsignedMessage.Bytes())
}

func constructUptimeMsgCmd() {
	args := struct {
		Network       string `cli:"#R, -n,   --network, Network (mainnet, fuji, etc)"`
		BlockchainID  string `cli:"#R, -b,   --blockchain, Blockchain ID"`
		ValidationID  string `cli:"#R, -v,   --validation, Validation ID"`
		UptimeSeconds string `cli:"#R, -t,   --uptime, Uptime Seconds"`
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
	uptimePayload, err := subnetevmmessages.NewValidatorUptime(validationID, uptimeSeconds)
	checkErr(err)
	addressedCall, err := payload.NewAddressedCall(nil, uptimePayload.Bytes())
	checkErr(err)
	uptimeProofUnsignedMessage, err := warp.NewUnsignedMessage(
		networkID,
		blockchainID,
		addressedCall.Bytes(),
	)
	checkErr(err)

	fmt.Println(utils.BytesToHex(uptimeProofUnsignedMessage.Bytes()))
}

func parseWarpMsgCmd() {
	args := struct {
		WarpMsg string `cli:"warpMsg, Warp Message"`
	}{}
	mcli.MustParse(&args)

	m, err := parseWarpMessage(args.WarpMsg)
	if err != nil {
		um, err := parseUnsignedWarpMessage(args.WarpMsg)
		checkErr(err)
		m = &warp.Message{UnsignedMessage: *um}
	} else {
		checkErr(err)
	}

	addressedCall, err := payload.ParseAddressedCall(m.Payload)
	checkErr(err)

	payloadType, payload, err := parsePayload(m.Payload)
	checkErr(err)

	fmt.Printf("\n%+v\n\nAddressedCall: %+v\n\nPayload (%s): %s\nPayload Hex: 0x%x\n", m, addressedCall, payloadType, payload, addressedCall.Payload)
}

func aggregateSignaturesCmd() {
	args := struct {
		SubnetID string `cli:"--subnet, SubnetID" default:"NOTSUPPLIED"`
		Msg      string `cli:"#R, msg, Warp Message"`
		URL      string `cli:"#R, --url, Glacier URL" default:"https://glacier-api.avax.network/v1/signatureAggregator/mainnet/aggregateSignatures"`
		Hex      bool   `cli:"--hex, Output as hex"`
	}{}
	mcli.MustParse(&args)

	var subnetIDPtr *ids.ID
	var err error
	if args.SubnetID != "NOTSUPPLIED" {
		subnetID, err := ids.FromString(args.SubnetID)
		checkErr(err)
		subnetIDPtr = &subnetID
	}

	msg, err := parseUnsignedWarpMessage(args.Msg)
	checkErr(err)

	c, err := sigagg.NewClient(args.URL)
	checkErr(err)

	msgSigned, err := c.AggregateSignatures(msg, subnetIDPtr, nil)
	checkErr(err)

	if args.Hex {
		fmt.Println(utils.BytesToHex(msgSigned.Bytes()))
	} else {
		fmt.Printf("\n%+v\n\nHex: %s\n", msgSigned, utils.BytesToHex(msgSigned.Bytes()))
	}
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
	return subnetevmwarp.UnpackSendWarpEventDataToMessage(log.Data)
}

// Returns json of the decoded payload
func parsePayload(msg []byte) (string, []byte, error) {
	addressedCall, err := payload.ParseAddressedCall(msg)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse addressed call: %w", err)
	}

	payloadIntf, err := message.Parse(addressedCall.Payload)
	if err != nil {
		return parseSubnetEvmPayload(addressedCall.Payload)
	}

	var out []byte

	switch payload := payloadIntf.(type) {
	case *message.RegisterL1Validator:
		out, err = utils.RegisterL1ValidatorToJSON(payload)
		if err != nil {
			return "", nil, err
		}
	case *message.L1ValidatorRegistration:
		out, err = json.Marshal(payload)
		if err != nil {
			return "", nil, err
		}
	case *message.L1ValidatorWeight:
		out, err = json.Marshal(payload)
		if err != nil {
			return "", nil, err
		}
	default:
		return "", nil, fmt.Errorf("unknown type: %T", payload)
	}

	return fmt.Sprintf("%T", payloadIntf), out, nil
}

func parseSubnetEvmPayload(payload []byte) (string, []byte, error) {
	payloadIntf, err := subnetevmmessages.Parse(payload)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	var out []byte
	switch payload := payloadIntf.(type) {
	case *subnetevmmessages.ValidatorUptime:
		out, err = json.Marshal(payload)
		if err != nil {
			return "", nil, err
		}

	default:
		return "", nil, fmt.Errorf("unknown type: %T", payload)
	}

	return fmt.Sprintf("%T", payloadIntf), out, nil
}
