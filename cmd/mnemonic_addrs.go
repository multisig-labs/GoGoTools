package main

import (
	"fmt"

	"github.com/jxskiss/mcli"
	"github.com/multisig-labs/gogotools/pkg/hd"

	"github.com/tyler-smith/go-bip39"
)

func mnemonicAddrsCmd() {
	args := struct {
		Mnemonic string `cli:"#R, mnemonic, BIP39 mnemonic"`
		Hrp      string `cli:"--hrp, hrp, Human-readable part (avax, fuji, local, etc)" default:"avax"`
		NumKeys  int    `cli:"--num-keys, num-keys, Number of keys to generate" default:"10"`
	}{}
	mcli.MustParse(&args)

	if ok := bip39.IsMnemonicValid(args.Mnemonic); !ok {
		checkErr("invaid mnemonic")
	}

	fmt.Printf("=== BIP39 Mnemonic ===\n%s\n\n", args.Mnemonic)

	fmt.Println("=== Ethereum Derivation Path ===")
	hdkeys, err := hd.DeriveHDKeys(args.Mnemonic, hd.EthDerivationPath, args.NumKeys)
	checkErr(err)

	fmtStr := "%-18s %-42s %-45s\n"

	fmt.Printf(fmtStr, "Path", "EVM Addr", "Ava Addr")
	for _, k := range hdkeys {
		fmt.Printf(fmtStr,
			k.Path,
			k.EthAddr(),
			k.AvaAddr("P", args.Hrp),
		)
	}

	fmt.Println("\n=== Avalanche Derivation Path ===")
	hdkeys, err = hd.DeriveHDKeys(args.Mnemonic, hd.AvaDerivationPath, args.NumKeys)
	checkErr(err)

	fmt.Printf(fmtStr, "Path", "EVM Addr", "Ava Addr")
	for _, k := range hdkeys {
		fmt.Printf(fmtStr,
			k.Path,
			k.EthAddr(),
			k.AvaAddr("P", args.Hrp),
		)
	}
}
