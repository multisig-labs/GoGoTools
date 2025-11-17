package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/jxskiss/mcli"
)

func inspectPTxCmd() {
	args := struct {
		Tx      string `cli:"#R, tx, P-Chain transaction to inspect (hex or base64)"`
		HexFlag bool   `cli:"--hex, Output unsigned tx as hex string" default:"false"`
	}{}
	mcli.MustParse(&args)

	txStr := args.Tx
	var txb []byte
	var err error
	if strings.HasPrefix(txStr, "0x") {
		txb, err = formatting.Decode(formatting.HexC, txStr)
		checkErr(err)
	} else {
		txb, err = b64.StdEncoding.DecodeString(txStr)
		checkErr(err)
	}

	tx := &txs.Tx{}
	_, err = txs.Codec.Unmarshal(txb, tx)
	checkErr(err)

	if args.HexFlag {
		unsignedBytes, err := txs.Codec.Marshal(txs.CodecVersion, &tx.Unsigned)
		checkErr(err)
		fmt.Printf("0x%x\n", unsignedBytes)
	} else {
		js, err := json.Marshal(tx)
		checkErr(err)

		fmt.Println(string(js))
	}
}
