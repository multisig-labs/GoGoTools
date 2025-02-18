package main

import (
	"fmt"

	"github.com/jxskiss/mcli"
	"github.com/multisig-labs/gogotools/pkg/utils"
)

func verifyBLSCmd() {
	args := struct {
		PublicKey string `cli:"#R, --bls-pubkey, Public key to verify"`
		Signature string `cli:"#R, --bls-pop, Proof of Posession signature to verify"`
	}{}
	mcli.MustParse(&args)

	err := utils.ValidateBLSKeys(args.PublicKey, args.Signature)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Signature is valid")
	}
}
