package validatormanagercmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ava-labs/coreth/accounts/abi/bind"
	"github.com/ava-labs/coreth/ethclient"
	"github.com/ava-labs/coreth/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/multisig-labs/gogotools/pkg/glacier"
	"github.com/multisig-labs/gogotools/pkg/uptime"
	"github.com/spf13/cobra"
)

// Multicall contract ABI and address
const MulticallABI = `[{"inputs":[{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call[]","name":"calls","type":"tuple[]"}],"name":"aggregate","outputs":[{"internalType":"uint256","name":"blockNumber","type":"uint256"},{"internalType":"bytes[]","name":"returnData","type":"bytes[]"}],"stateMutability":"payable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bool","name":"allowFailure","type":"bool"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call3[]","name":"calls","type":"tuple[]"}],"name":"aggregate3","outputs":[{"components":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"bytes","name":"returnData","type":"bytes"}],"internalType":"struct Multicall3.Result[]","name":"returnData","type":"tuple[]"}],"stateMutability":"payable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bool","name":"allowFailure","type":"bool"},{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call3Value[]","name":"calls","type":"tuple[]"}],"name":"aggregate3Value","outputs":[{"components":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"bytes","name":"returnData","type":"bytes"}],"internalType":"struct Multicall3.Result[]","name":"returnData","type":"tuple[]"}],"stateMutability":"payable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call[]","name":"calls","type":"tuple[]"}],"name":"blockAndAggregate","outputs":[{"internalType":"uint256","name":"blockNumber","type":"uint256"},{"internalType":"bytes32","name":"blockHash","type":"bytes32"},{"components":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"bytes","name":"returnData","type":"bytes"}],"internalType":"struct Multicall3.Result[]","name":"returnData","type":"tuple[]"}],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"getBasefee","outputs":[{"internalType":"uint256","name":"basefee","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"getBlockHash","outputs":[{"internalType":"bytes32","name":"blockHash","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getBlockNumber","outputs":[{"internalType":"uint256","name":"blockNumber","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getChainId","outputs":[{"internalType":"uint256","name":"chainid","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCurrentBlockCoinbase","outputs":[{"internalType":"address","name":"coinbase","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCurrentBlockDifficulty","outputs":[{"internalType":"uint256","name":"difficulty","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCurrentBlockGasLimit","outputs":[{"internalType":"uint256","name":"gaslimit","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCurrentBlockTimestamp","outputs":[{"internalType":"uint256","name":"timestamp","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"getEthBalance","outputs":[{"internalType":"uint256","name":"balance","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getLastBlockHash","outputs":[{"internalType":"bytes32","name":"blockHash","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bool","name":"requireSuccess","type":"bool"},{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call[]","name":"calls","type":"tuple[]"}],"name":"tryAggregate","outputs":[{"components":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"bytes","name":"returnData","type":"bytes"}],"internalType":"struct Multicall3.Result[]","name":"returnData","type":"tuple[]"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bool","name":"requireSuccess","type":"bool"},{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call[]","name":"calls","type":"tuple[]"}],"name":"tryBlockAndAggregate","outputs":[{"internalType":"uint256","name":"blockNumber","type":"uint256"},{"internalType":"bytes32","name":"blockHash","type":"bytes32"},{"components":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"bytes","name":"returnData","type":"bytes"}],"internalType":"struct Multicall3.Result[]","name":"returnData","type":"tuple[]"}],"stateMutability":"payable","type":"function"}]`

var MulticallAddress = common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")

func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [l1-config-file]",
		Short: "Get info about the validator manager",
		Long: `
		Config file is a JSON file with the following structure:
		{
			"name": "coqnet",
			"network": "fuji",
			"subnet_id": "4YurNFwLzhGUrYyihDnUUc2L199YBnFeWP3fhJKmDDjkbvy8G",
			"blockchain_id": "EyWDF1cGmMKXRi4d5Mb1kVNnB1zHWvMPGQM5uHgizNGFTvsTn",
			"vm_subnet_id": "11111111111111111111111111111111LpoYY",
			"vm_blockchain_id": "yH8D7ThNJkxmtkuv2jgBa4P1Rn3Qpr4pPr7QYNfcdoS6k6HWp",
			"vm_address": "0x0ec8f51391b3976b406ec182c8c22e537ff14eca",
			"validators_url": "https://testnet-coqnetfuji-wd58c.avax-test.network/ext/bc/EyWDF1cGmMKXRi4d5Mb1kVNnB1zHWvMPGQM5uHgizNGFTvsTn/validators",
			"ava_url": "https://api.avax-test.network",
			"evm_url": "https://api.avax-test.network/ext/bc/C/rpc"
		}
		`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			configFile := args[0]
			fileContent, err := os.ReadFile(configFile)
			cobra.CheckErr(err)

			l1Config := L1Config{}
			err = json.Unmarshal(fileContent, &l1Config)
			cobra.CheckErr(err)

			_, err = constants.NetworkID(l1Config.Network)
			cobra.CheckErr(err)

			info := Info{Config: l1Config, ContractValidators: []Validator{}}

			client := glacier.NewClientWithConfig(glacier.Config{
				Network: l1Config.Network,
			})

			info.PchainValidators, err = client.FetchValidators(l1Config.SubnetID)
			cobra.CheckErr(err)

			uptimeClient, err := uptime.NewClient(l1Config.ValidatorsURL)
			cobra.CheckErr(err)

			info.Uptime, _ = uptimeClient.GetCurrentValidators(context.Background())
			// cobra.CheckErr(err)

			ec, err := ethclient.Dial(l1Config.EvmURL)
			cobra.CheckErr(err)

			// Use Multicall3 and batches to efficiently get validators from the contract
			validationIDs := []ids.ID{}
			for _, validator := range info.PchainValidators {
				validationIDs = append(validationIDs, validator.ValidationID)
			}

			// Process validators in batches of 50
			batchSize := 50
			var allValidators []Validator
			for i := 0; i < len(validationIDs); i += batchSize {
				end := i + batchSize
				if end > len(validationIDs) {
					end = len(validationIDs)
				}
				batch, err := getValidatorsBatch(ec, l1Config.VMAddress, validationIDs[i:end])
				if err == nil {
					allValidators = append(allValidators, batch...)
				}
			}
			info.ContractValidators = allValidators

			info.ContractSettings, err = getValidatorManagerSettings(l1Config.EvmURL, l1Config.VMAddress)
			cobra.CheckErr(err)

			infoJSON, err := json.MarshalIndent(info, "", "  ")
			cobra.CheckErr(err)
			fmt.Println(string(infoJSON))
		},
	}
	return cmd
}

func getValidatorsBatch(ec ethclient.Client, contractAddress common.Address, validationIDs []ids.ID) ([]Validator, error) {
	// Parse the validator method ABI
	methodSpec := "getValidator(bytes32 validationID)->((uint8 status,bytes nodeID,uint64 startingWeight,uint64 messageNonce,uint64 weight,uint64 startedAt,uint64 endedAt))"
	methodName, methodABI, err := ParseSpec(methodSpec, nil, false, false, false, true, ids.ID{})
	if err != nil {
		return nil, err
	}

	metadata := &bind.MetaData{
		ABI: methodABI,
	}

	contractAbi, err := metadata.GetAbi()
	if err != nil {
		return nil, err
	}

	// Parse Multicall ABI
	multicallAbi, err := abi.JSON(strings.NewReader(MulticallABI))
	if err != nil {
		return nil, err
	}

	// Prepare multicall inputs
	calls := make([]struct {
		Target   common.Address `json:"target"`
		CallData []byte         `json:"callData"`
	}, len(validationIDs))

	// Pack each validator call
	for i, validationID := range validationIDs {
		callData, err := contractAbi.Pack(methodName, validationID)
		if err != nil {
			return nil, fmt.Errorf("failed to pack call data for validator %s: %w", validationID, err)
		}

		calls[i] = struct {
			Target   common.Address `json:"target"`
			CallData []byte         `json:"callData"`
		}{
			Target:   contractAddress,
			CallData: callData,
		}
	}

	// Pack multicall input
	multicallInput, err := multicallAbi.Pack("aggregate", calls)
	if err != nil {
		return nil, fmt.Errorf("failed to pack multicall input: %w", err)
	}

	// Make the call
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msg := interfaces.CallMsg{
		To:   &MulticallAddress,
		Data: multicallInput,
	}

	result, err := ec.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("multicall failed: %w", err)
	}

	// Unpack multicall result
	multicallOutput, err := multicallAbi.Unpack("aggregate", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack multicall result: %w", err)
	}

	// Extract return data
	returnData, ok := multicallOutput[1].([][]byte)
	if !ok {
		return nil, fmt.Errorf("invalid multicall return data type")
	}

	// Process results
	validators := make([]Validator, len(returnData))
	for i, data := range returnData {
		var out interface{}
		err = contractAbi.UnpackIntoInterface(&out, methodName, data)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack validator %d result: %w", i, err)
		}
		validators[i] = *abi.ConvertType(out, new(Validator)).(*Validator)
		validators[i].ValidationID = validationIDs[i]
	}

	return validators, nil
}

// func getValidator(ec ethclient.Client, contractAddress common.Address, validationID ids.ID) (Validator, error) {
// 	spec := "getValidator(bytes32 validationID)->((uint8 status,bytes nodeID,uint64 startingWeight,uint64 messageNonce,uint64 weight,uint64 startedAt,uint64 endedAt))"
// 	validatorSlice, err := callToMethod(ec, contractAddress, spec, validationID)
// 	if err != nil {
// 		return Validator{}, fmt.Errorf("failed to get validator: %w", err)
// 	}
// 	out0 := *abi.ConvertType(validatorSlice[0], new(Validator)).(*Validator)

// 	return out0, nil
// }

// func callToMethod(ec ethclient.Client, contractAddress common.Address, methodSpec string, params ...interface{}) ([]interface{}, error) {
// 	methodName, methodABI, err := ParseSpec(methodSpec, nil, false, false, false, true, params...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	metadata := &bind.MetaData{
// 		ABI: methodABI,
// 	}

// 	abi, err := metadata.GetAbi()
// 	if err != nil {
// 		return nil, err
// 	}

// 	contract := bind.NewBoundContract(contractAddress, *abi, ec, ec, ec)
// 	var out []interface{}
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	err = contract.Call(&bind.CallOpts{Context: ctx}, &out, methodName, params...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return out, nil
// }
