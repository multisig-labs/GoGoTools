package main

import (
	"fmt"
	"strings"

	"github.com/tyler-smith/go-bip39"
)

func mnemonicInsecureCmd() {
	entropy, _ := bip39.NewEntropy(128)
	phrase, _ := bip39.NewMnemonic(entropy)
	words := strings.Split(phrase, " ")[0:11]
	words[0] = "test"
	words[1] = "test"
	wordList := bip39.GetWordList()
	var tryMnemonic string
	for _, word := range wordList {
		// Construct a possible mnemonic by adding each word as the 12th word
		tryMnemonic = strings.Join(words, " ") + " " + word
		// Check if the constructed mnemonic is valid
		if bip39.IsMnemonicValid(tryMnemonic) {
			break
		}
	}

	fmt.Println(tryMnemonic)
}
