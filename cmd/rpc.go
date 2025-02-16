package main

import (
	"fmt"

	"github.com/jxskiss/mcli"
	"github.com/multisig-labs/gogotools/pkg/utils"
)

func rpcCmd() {
	args := struct {
		Endpoint string `cli:"#R, endpoint, RPC endpoint"`
		Method   string `cli:"#R, method, RPC method"`
		Params   string `cli:"params, RPC params"`
		URLFlags
	}{}
	mcli.MustParse(&args)

	url := fmt.Sprintf("%s%s", args.AvaUrl, args.Endpoint)
	result, err := utils.FetchRPCGJSON(url, args.Method, args.Params)
	checkErr(err)
	fmt.Println(result.String())
}
