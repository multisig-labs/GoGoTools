package main

import (
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/libevm/common"
	"github.com/jxskiss/mcli"
)

func bech32DecodeCmd() {
	args := struct {
		Value string `cli:"#R, value, bech32 address to decode to bytes"`
	}{}
	mcli.MustParse(&args)

	_, addrBytes, err := address.ParseBech32(strings.TrimPrefix(args.Value, "P-"))
	checkErr(err)
	fmt.Printf("0x%x\n", addrBytes)
}

func bech32EncodeCmd() {
	args := struct {
		Hrp   string `cli:"--hrp, hrp, Human-readable part (avax, fuji, local, etc)" default:"avax"`
		Value string `cli:"#R, value, hex string value to encode"`
	}{}
	mcli.MustParse(&args)

	addr, err := address.FormatBech32(args.Hrp, common.FromHex(args.Value))
	checkErr(err)
	fmt.Printf("%s\n", addr)
}
