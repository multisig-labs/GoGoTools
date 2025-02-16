package main

import (
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/jxskiss/mcli"
)

func bech32DecodeCmd() {
	args := struct {
		Value string `cli:"#R,value, bech32 address to decode to bytes"`
	}{}
	mcli.MustParse(&args)

	_, addrBytes, err := address.ParseBech32(strings.TrimPrefix(args.Value, "P-"))
	checkErr(err)
	fmt.Printf("0x%x\n", addrBytes)
}
