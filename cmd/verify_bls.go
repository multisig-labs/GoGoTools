package main

import (
	"fmt"

	"github.com/multisig-labs/gogotools/pkg/utils"
)

func verifyBLSCmd() {
	args := struct {
		PublicKey string `cli:"--public-key, Public key to verify"`
		Signature string `cli:"--signature, Signature to verify"`
	}{}

	err := utils.ValidateBLSKeys(args.PublicKey, args.Signature)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Signature is valid")
	}
}
