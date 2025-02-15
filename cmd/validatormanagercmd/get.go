package validatormanagercmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ava-labs/coreth/accounts/abi/bind"
	"github.com/ava-labs/coreth/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [contractAddress] [validationID]",
		Short: "Get info about the validator manager",
		Long: `
		`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			contractAddr := common.HexToAddress(args[0])
			validationID, err := ids.FromString(args[1])
			cobra.CheckErr(err)
			rpcURL := viper.GetString("node-url")
			ec, err := ethclient.Dial(rpcURL)
			cobra.CheckErr(err)

			vdr, err := getValidator(ec, contractAddr, validationID)
			cobra.CheckErr(err)

			j, err := json.MarshalIndent(vdr, "", "  ")
			cobra.CheckErr(err)
			fmt.Println(string(j))
		},
	}
	return cmd
}

func getValidator(ec ethclient.Client, contractAddress common.Address, validationID ids.ID) (*Validator, error) {
	spec := "getValidator(bytes32 validationID)->((uint8 status,bytes nodeID,uint64 startingWeight,uint64 messageNonce,uint64 weight,uint64 startedAt,uint64 endedAt))"
	validatorSlice, err := callToMethod(ec, contractAddress, spec, validationID)
	if err != nil {
		return &Validator{}, fmt.Errorf("failed to get validator: %w", err)
	}
	out0 := *abi.ConvertType(validatorSlice[0], new(Validator)).(*Validator)
	out0.ValidationID = validationID
	return &out0, nil
}

func callToMethod(ec ethclient.Client, contractAddress common.Address, methodSpec string, params ...interface{}) ([]interface{}, error) {
	methodName, methodABI, err := ParseSpec(methodSpec, nil, false, false, false, true, params...)
	if err != nil {
		return nil, err
	}

	metadata := &bind.MetaData{
		ABI: methodABI,
	}

	abi, err := metadata.GetAbi()
	if err != nil {
		return nil, err
	}

	contract := bind.NewBoundContract(contractAddress, *abi, ec, ec, ec)
	var out []interface{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = contract.Call(&bind.CallOpts{Context: ctx}, &out, methodName, params...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
