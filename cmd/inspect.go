package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jxskiss/mcli"
)

func inspectPTxCmd() {
	args := struct {
		Tx string `cli:"#R, tx, Transaction to inspect (hex or base64)"`
	}{}
	mcli.MustParse(&args)

	txStr := args.Tx
	var txb []byte
	var err error
	if strings.HasPrefix(txStr, "0x") {
		txb, err = hexutil.Decode(txStr)
		checkErr(err)
	} else {
		txb, err = b64.StdEncoding.DecodeString(txStr)
		checkErr(err)
	}

	tx := &txs.Tx{}
	_, err = txs.Codec.Unmarshal(txb, tx)
	checkErr(err)

	js, err := json.Marshal(tx)
	checkErr(err)

	fmt.Println(string(js))
}
