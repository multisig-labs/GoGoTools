package main

import (
	"encoding/json"
	"fmt"

	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/jxskiss/mcli"
	"github.com/multisig-labs/gogotools/pkg/utils"
)

func privkeyCmd() {
	args := struct {
		Key  string `cli:"#R, key, Private key"`
		Json bool   `cli:"--json, Output in JSON format" default:"false"`
	}{}
	mcli.MustParse(&args)

	avaKey, ethKey, err := utils.ParsePrivateKey(args.Key)
	checkErr(err)

	avaKeyCB58, err := cb58.Encode(avaKey.Bytes())
	checkErr(err)
	avaKeyCB58 = "PrivateKey-" + avaKeyCB58

	avaKeyBytes := avaKey.Bytes()

	ethAddr := ethcrypto.PubkeyToAddress(ethKey.PublicKey).String()

	avaAddrMainnet, err := address.Format("P", "avax", avaKey.PublicKey().Address().Bytes())
	checkErr(err)

	avaAddrFuji, err := address.Format("P", "fuji", avaKey.PublicKey().Address().Bytes())
	checkErr(err)

	avaAddrLocal, err := address.Format("P", "local", avaKey.PublicKey().Address().Bytes())
	checkErr(err)

	if args.Json {
		jsonBytes, err := json.MarshalIndent(struct {
			PrivKeyHex     string `json:"pk_hex"`
			PrivKeyCB58    string `json:"pk_cb58"`
			EthAddr        string `json:"eth_addr"`
			AvaAddrMainnet string `json:"ava_addr_mainnet"`
			AvaAddrFuji    string `json:"ava_addr_fuji"`
			AvaAddrLocal   string `json:"ava_addr_local"`
		}{
			PrivKeyHex:     utils.BytesToHex(avaKeyBytes),
			PrivKeyCB58:    avaKeyCB58,
			EthAddr:        ethAddr,
			AvaAddrMainnet: avaAddrMainnet,
			AvaAddrFuji:    avaAddrFuji,
			AvaAddrLocal:   avaAddrLocal,
		}, "", "  ")
		checkErr(err)
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("%-14s %#x\n", "PrivKey Hex:", avaKeyBytes)
		fmt.Printf("%-14s %s\n", "PrivKey CB58:", avaKeyCB58)
		fmt.Printf("%-14s %s\n", "Eth addr:", ethAddr)
		fmt.Printf("%-14s %s\n", "Ava addr:", avaAddrMainnet)
		fmt.Printf("%-14s %s\n", "Ava addr:", avaAddrFuji)
		fmt.Printf("%-14s %s\n", "Ava addr:", avaAddrLocal)
	}
}
