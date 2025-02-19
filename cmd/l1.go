package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jxskiss/mcli"
	"github.com/multisig-labs/gogotools/pkg/uptime"
)

func l1ValidatorsCmd() {
	args := struct {
		URL string `cli:"#R, --url, URL Of the L1 Validator RPC Endpoint (i.e. https://node.myl1.network/ext/bc/<blockchainID>/validators)"`
	}{}
	mcli.MustParse(&args)

	client, err := uptime.NewClient(args.URL)
	checkErr(err)

	validators, err := client.GetCurrentValidators(context.Background())
	checkErr(err)

	js, err := json.Marshal(validators)
	checkErr(err)

	fmt.Println(string(js))
}
