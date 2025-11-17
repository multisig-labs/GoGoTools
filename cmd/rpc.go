package main

import (
	"fmt"
	"strings"

	"github.com/jxskiss/mcli"
	"github.com/multisig-labs/gogotools/pkg/utils"
)

func rpcCmd() {
	args := struct {
		Endpoint string `cli:"#R, endpoint, RPC endpoint"`
		Method   string `cli:"#R, method, RPC method"`
		Params   string `cli:"params, RPC params"`
		AvaUrl   string `cli:"--ava-url, URL of the Avalanche node (do not include a path)." env:"AVA_RPC_URL,PCHAIN_RPC_URL,NODE_URL" default:"https://api.avax.network"`
	}{}
	mcli.MustParse(&args)

	var url string
	if strings.HasPrefix(args.Endpoint, "http") {
		url = args.Endpoint
	} else {
		url = fmt.Sprintf("%s%s", args.AvaUrl, args.Endpoint)
	}
	result, err := utils.FetchRPCGJSON(url, args.Method, args.Params)
	checkErr(err)
	fmt.Println(result.String())
}
