package validatormanagercmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"os/exec"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Map to the same struct in viewer.sol
type ContractSettings struct {
	SubnetID                 ids.ID         `abi:"_subnetID"`
	ChurnPeriodSeconds       uint64         `abi:"_churnPeriodSeconds"`
	MaximumChurnPercentage   uint8          `abi:"_maximumChurnPercentage"`
	MinimumStakeAmount       *big.Int       `abi:"_minimumStakeAmount"`
	MaximumStakeAmount       *big.Int       `abi:"_maximumStakeAmount"`
	MinimumStakeDuration     uint64         `abi:"_minimumStakeDuration"`
	MinimumDelegationFeeBips uint16         `abi:"_minimumDelegationFeeBips"`
	MaximumStakeMultiplier   uint64         `abi:"_maximumStakeMultiplier"`
	WeightToValueFactor      *big.Int       `abi:"_weightToValueFactor"`
	RewardCalculator         common.Address `abi:"_rewardCalculator"`
	UptimeBlockchainID       ids.ID         `abi:"_uptimeBlockchainID"`
}

// ContractResponse represents the solc compiler output
type ContractResponse struct {
	Contracts map[string]struct {
		BinRuntime string          `json:"bin-runtime"`
		Abi        json.RawMessage `json:"abi"`
	} `json:"contracts"`
}

// JsonRPCRequest represents an Ethereum JSON-RPC request
type JsonRPCRequest struct {
	ID      int           `json:"id"`
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

//go:embed viewer.sol
var viewerSol string

func getValidatorManagerSettings(evmRpc string, contractAddress common.Address) (ContractSettings, error) {
	var err error

	// Create a temporary file for viewer.sol
	tmpFile, err := os.CreateTemp("", "viewer-*.sol")
	if err != nil {
		return ContractSettings{}, fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temp file when done

	// Write the embedded content to the temp file
	if _, err := tmpFile.WriteString(viewerSol); err != nil {
		return ContractSettings{}, fmt.Errorf("failed to write to temporary file: %v", err)
	}
	tmpFile.Close()

	// Use the temporary file with solc
	cmd := exec.Command("solc", "--no-cbor-metadata", "--combined-json", "abi,bin-runtime", tmpFile.Name())
	output, err := cmd.Output()
	if err != nil {
		return ContractSettings{}, fmt.Errorf("failed to compile viewer.sol contract (solc must be in your path): %v", err)
	}

	// Parse the solc output
	var contractResp ContractResponse
	if err := json.Unmarshal(output, &contractResp); err != nil {
		return ContractSettings{}, fmt.Errorf("failed to parse solc output: %v", err)
	}

	contract := contractResp.Contracts[tmpFile.Name()+":Viewer"]
	bytecode := contract.BinRuntime

	// Prepare the JSON-RPC request
	callData := "0x85b4bb53" // function selector "getSettings()"
	params := []interface{}{
		map[string]string{
			"to":   contractAddress.Hex(),
			"data": callData,
		},
		"latest",
		map[string]interface{}{
			contractAddress.Hex(): map[string]string{
				"code": "0x" + bytecode,
			},
		},
	}

	rpcReq := JsonRPCRequest{
		ID:      1,
		JsonRPC: "2.0",
		Method:  "eth_call",
		Params:  params,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(rpcReq)
	if err != nil {
		return ContractSettings{}, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Make the HTTP request
	req, err := http.NewRequest("POST", evmRpc, bytes.NewBuffer(jsonData))
	if err != nil {
		return ContractSettings{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ContractSettings{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ContractSettings{}, fmt.Errorf("failed to read response: %v", err)
	}

	var result struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return ContractSettings{}, fmt.Errorf("failed to parse response: %v", err)
	}

	// Decode the struct using the ABI
	parsedAbi, err := abi.JSON(bytes.NewReader(contract.Abi))
	if err != nil {
		return ContractSettings{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	data := hexutil.MustDecode(result.Result)

	// Decode the struct
	decodedStruct, err := parsedAbi.Unpack("getSettings", data)
	if err != nil {
		return ContractSettings{}, fmt.Errorf("failed to decode struct: %v", err)
	}

	// Get the anonymous struct as a map
	rawStruct := decodedStruct[0].(struct {
		SubnetID                 [32]byte       `json:"_subnetID"`
		ChurnPeriodSeconds       uint64         `json:"_churnPeriodSeconds"`
		MaximumChurnPercentage   uint8          `json:"_maximumChurnPercentage"`
		MinimumStakeAmount       *big.Int       `json:"_minimumStakeAmount"`
		MaximumStakeAmount       *big.Int       `json:"_maximumStakeAmount"`
		MinimumStakeDuration     uint64         `json:"_minimumStakeDuration"`
		MinimumDelegationFeeBips uint16         `json:"_minimumDelegationFeeBips"`
		MaximumStakeMultiplier   uint64         `json:"_maximumStakeMultiplier"`
		WeightToValueFactor      *big.Int       `json:"_weightToValueFactor"`
		RewardCalculator         common.Address `json:"_rewardCalculator"`
		UptimeBlockchainID       [32]byte       `json:"_uptimeBlockchainID"`
	})

	// Convert to ViewerSettings
	viewerSettings := ContractSettings{
		SubnetID:                 ids.ID(rawStruct.SubnetID),
		ChurnPeriodSeconds:       rawStruct.ChurnPeriodSeconds,
		MaximumChurnPercentage:   rawStruct.MaximumChurnPercentage,
		MinimumStakeAmount:       rawStruct.MinimumStakeAmount,
		MaximumStakeAmount:       rawStruct.MaximumStakeAmount,
		MinimumStakeDuration:     rawStruct.MinimumStakeDuration,
		MinimumDelegationFeeBips: rawStruct.MinimumDelegationFeeBips,
		MaximumStakeMultiplier:   rawStruct.MaximumStakeMultiplier,
		WeightToValueFactor:      rawStruct.WeightToValueFactor,
		RewardCalculator:         rawStruct.RewardCalculator,
		UptimeBlockchainID:       ids.ID(rawStruct.UptimeBlockchainID),
	}

	return viewerSettings, nil
}
