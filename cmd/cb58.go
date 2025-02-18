package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jxskiss/mcli"
)

func cb58EncodeCmd() {
	args := struct {
		Value string `cli:"#R, value, Value to encode"`
	}{}
	mcli.MustParse(&args)

	b, err := hexutil.Decode(args.Value)
	checkErr(err)
	cb, err := cb58.Encode(b)
	checkErr(err)
	fmt.Printf("%s\n", cb)
}

func cb58DecodeCmd() {
	args := struct {
		Value string `cli:"#R, value, Value to decode (address or NodeID)"`
	}{}
	mcli.MustParse(&args)

	strippedValue := strings.TrimPrefix(args.Value, "NodeID-")
	b, err := cb58.Decode(strippedValue)
	checkErr(err)
	fmt.Printf("0x%x\n", b)
}

func cb58DecodeSigCmd() {
	args := struct {
		Value string `cli:"#R, value, Value to decode"`
	}{}
	mcli.MustParse(&args)

	b, err := cb58.Decode(args.Value)
	checkErr(err)
	sig := struct {
		R string `json:"r"`
		S string `json:"s"`
		V string `json:"v"`
	}{}
	sig.R = fmt.Sprintf("0x%x", b[0:32])
	sig.S = fmt.Sprintf("0x%x", b[32:64])
	sig.V = fmt.Sprintf("0x%x", b[64:])
	j, _ := json.Marshal(sig)
	fmt.Printf("%s\n", j)
}
