package main

import (
	"encoding/json"
	"fmt"

	"github.com/ava-labs/avalanchego/utils/formatting/address"
	ethcrypto "github.com/ava-labs/libevm/crypto"
	"github.com/jxskiss/mcli"
	"github.com/multisig-labs/gogotools/pkg/utils"
)

func addrsCmd() {
	args := struct {
		Key  string `cli:"#R, key, Public key"`
		Json bool   `cli:"--json, Output in JSON format" default:"false"`
	}{}
	mcli.MustParse(&args)

	pubkey, err := utils.ParsePublicKey(args.Key)
	checkErr(err)

	ethAddr := ethcrypto.PubkeyToAddress(*pubkey.ToECDSA()).String()

	avaAddrMainnet, err := address.Format("P", "avax", pubkey.Address().Bytes())
	checkErr(err)

	avaAddrFuji, err := address.Format("P", "fuji", pubkey.Address().Bytes())
	checkErr(err)

	avaAddrLocal, err := address.Format("P", "local", pubkey.Address().Bytes())
	checkErr(err)

	if args.Json {
		jsonBytes, err := json.MarshalIndent(struct {
			EthAddr        string `json:"eth_addr"`
			AvaAddrMainnet string `json:"ava_addr_mainnet"`
			AvaAddrFuji    string `json:"ava_addr_fuji"`
			AvaAddrLocal   string `json:"ava_addr_local"`
		}{
			EthAddr:        ethAddr,
			AvaAddrMainnet: avaAddrMainnet,
			AvaAddrFuji:    avaAddrFuji,
			AvaAddrLocal:   avaAddrLocal,
		}, "", "  ")
		checkErr(err)
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("%-20s %x\n", "PubKey (Compressed):", pubkey.Bytes())
		fmt.Printf("%-20s %s\n", "Eth addr:", ethAddr)
		fmt.Printf("%-20s %s\n", "Ava addr:", avaAddrMainnet)
		fmt.Printf("%-20s %s\n", "Ava addr:", avaAddrFuji)
		fmt.Printf("%-20s %s\n", "Ava addr:", avaAddrLocal)
	}
}
