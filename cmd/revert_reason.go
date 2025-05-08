package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc" // Required for rpc.Error type assertion
	"github.com/jxskiss/mcli"
)

func revertReasonCmd() {
	args := struct {
		TxHash string `cli:"#R, tx, Tx hash"`
		URLFlags
	}{}
	mcli.MustParse(&args)

	c, err := ethclient.Dial(args.EthUrl)
	checkErr(err)

	s, err := getRevertReason(c, common.HexToHash(args.TxHash))
	checkErr(err)

	fmt.Println(s)
}

// getRevertReason attempts to retrieve the revert reason for a failed transaction.
func getRevertReason(client *ethclient.Client, txHash common.Hash) (string, error) {
	ctx := context.Background()

	// 1. Get the transaction receipt
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction receipt for %s: %w", txHash.Hex(), err)
	}

	// Check if the transaction actually failed
	if receipt.Status == types.ReceiptStatusSuccessful {
		return "", fmt.Errorf("transaction %s was successful, no revert reason", txHash.Hex())
	}

	// 2. Get the original transaction details
	tx, _, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction by hash for %s: %w", txHash.Hex(), err)
	}

	// 3. Re-simulate the transaction using eth_call
	// We need the sender address. This can be derived from the transaction's signature.
	// types.Sender uses the EIP155 signer, which is common.
	// Ensure tx.ChainId() is correct for the network.
	var signer types.Signer
	if tx.ChainId() != nil {
		signer = types.LatestSignerForChainID(tx.ChainId())
	} else {
		// Handle older, non-EIP155 transactions if necessary, or assume a default.
		// For simplicity, this example might require a chain ID.
		// On some networks/nodes, tx.From() might be populated directly by the RPC.
		// If not, deriving it is crucial.
		return "", fmt.Errorf("transaction chain ID is nil, cannot determine signer")
	}
	from, err := types.Sender(signer, tx)
	if err != nil {
		return "", fmt.Errorf("failed to get sender from transaction %s: %w", txHash.Hex(), err)
	}

	callMsg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(), // For legacy transactions
		// GasFeeCap: tx.GasFeeCap(), // For EIP-1559 transactions
		// GasTipCap: tx.GasTipCap(), // For EIP-1559 transactions
		Value:      tx.Value(),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
	}

	// Perform the call against the state of the block in which the transaction was included
	// (and failed). receipt.BlockNumber is the correct block for this.
	result, err := client.CallContract(ctx, callMsg, receipt.BlockNumber)

	if err == nil {
		// This shouldn't happen if the original tx failed and we're replaying it.
		// It might mean the state changed or the node doesn't support revert reasons well.
		return "", fmt.Errorf("eth_call succeeded for failed tx %s (result: %x), unexpected", txHash.Hex(), result)
	}

	// 4. Decode the revert reason from the error
	// The error from CallContract often contains the revert reason.
	// It might be directly in the error message or in a structured data field.

	// Try to parse it as an ABI-encoded revert reason string "Error(string)"
	revertReason, unpackErr := tryUnpackRevertError(err)
	if unpackErr == nil {
		return revertReason, nil
	}

	// If unpacking failed, return the raw error message, which might contain the reason.
	// Some nodes return it like "execution reverted: <reason>"
	errMsg := err.Error()
	if strings.Contains(errMsg, "execution reverted") {
		// Simple extraction, might need refinement
		parts := strings.SplitN(errMsg, "execution reverted", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
			reason := strings.TrimPrefix(strings.TrimSpace(parts[1]), ": ")
			return reason, nil
		}
	}

	// Fallback: return the full error and the unpack error if any
	return "", fmt.Errorf("failed to get revert reason for %s. Call error: %w. Unpack error: %v", txHash.Hex(), err, unpackErr)
}

// tryUnpackRevertError attempts to decode an ABI-encoded revert reason from the error data.
// Revert reasons are often encoded as Error(string), which has a 4-byte selector 0x08c379a0.
func tryUnpackRevertError(callErr error) (string, error) {
	// Check if the error is an rpc.Error, which might contain a Data field.
	if rpcErr, ok := callErr.(rpc.DataError); ok {
		if data, ok := rpcErr.ErrorData().(string); ok {
			// Data is often a hex string "0x..."
			decodedReason, err := abi.UnpackRevert(common.FromHex(data))
			if err == nil {
				return decodedReason, nil
			}
			return "", fmt.Errorf("failed to abi.UnpackRevert data '%s': %w", data, err)
		}
		return "", fmt.Errorf("rpc.Error.Data() is not a string: %T", rpcErr.ErrorData())
	}

	// The error might be an abi.CallError (though less common directly from ethclient.CallContract)
	// For example:
	// if callError, ok := callErr.(abi.CallError); ok {
	//    return abi.UnpackRevert(callError.Data)
	// }

	return "", fmt.Errorf("error is not an rpc.Error or known revert error type: %T, err: %v", callErr, callErr)
}
