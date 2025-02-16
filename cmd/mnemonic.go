package main

import (
	"fmt"

	"github.com/tyler-smith/go-bip39"
)

func randomMnemonicCmd() {
	entropy, _ := bip39.NewEntropy(256)
	phrase, _ := bip39.NewMnemonic(entropy)

	fmt.Println(phrase)
}
