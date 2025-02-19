package main

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/tidwall/sjson"
	"github.com/tyler-smith/go-bip39"
)

func randomIDCmd() {
	id, err := utils.RandomID()
	checkErr(err)
	fmt.Println(id)
}

func randomNodeIDCmd() {
	nodeid, err := utils.RandomNodeID()
	checkErr(err)
	fmt.Println(nodeid)
}

func randomBLSCmd() {
	sk, pop, err := utils.RandomBLS()
	checkErr(err)
	popjs, err := pop.MarshalJSON()
	checkErr(err)

	skBytes := fmt.Sprintf("0x%x", sk.ToBytes())

	out, err := sjson.SetBytes(popjs, "privateKey", skBytes)
	checkErr(err)

	fmt.Println(string(out))
}

func randomMnemonicCmd() {
	entropy, _ := bip39.NewEntropy(256)
	phrase, _ := bip39.NewMnemonic(entropy)

	fmt.Println(phrase)
}
