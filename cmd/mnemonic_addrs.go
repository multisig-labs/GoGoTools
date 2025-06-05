package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
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

func xpubkeyCmd() {
	args := struct {
		Mnemonic       string `cli:"#R, mnemonic, BIP39 mnemonic"`
		DerivationPath string `cli:"--path, Derivation path" default:"m/44'/60'/0'/0/0"`
	}{}
	mcli.MustParse(&args)

	if ok := bip39.IsMnemonicValid(args.Mnemonic); !ok {
		checkErr("invaid mnemonic")
	}

	path, err := accounts.ParseDerivationPath(args.DerivationPath)
	checkErr(err)

	xpub, err := hd.XPubKey(args.Mnemonic, path)
	checkErr(err)

	fmt.Println(xpub.String())
}

func xpubAddrsCmd() {
	args := struct {
		XPub           string `cli:"#R, xpub, Extended public key"`
		DerivationPath string `cli:"--path, Derivation path" default:"m/44'/60'/0'/0/0"`
		Hrp            string `cli:"--hrp, hrp, Human-readable part (avax, fuji, local, etc)" default:"avax"`
		NumKeys        int    `cli:"--num-keys, num-keys, Number of keys to generate" default:"10"`
		Json           bool   `cli:"--json, json, Output in JSON format"`
	}{}
	mcli.MustParse(&args)

	xpub, err := hdkeychain.NewKeyFromString(args.XPub)
	checkErr(err)

	path, err := accounts.ParseDerivationPath(args.DerivationPath)
	checkErr(err)

	hdkeys, err := hd.DerivePubKeys(xpub, path, args.NumKeys)
	checkErr(err)

	if args.Json {
		out := struct {
			Xpub     string   `json:"xpub"`
			Path     string   `json:"path"`
			EvmAddrs []string `json:"evm_addrs"`
			AvaAddrs []string `json:"ava_addrs"`
		}{}
		out.Path = args.DerivationPath[:len(args.DerivationPath)-1] + "*"
		out.Xpub = args.XPub
		out.EvmAddrs = make([]string, 0, len(hdkeys))
		out.AvaAddrs = make([]string, 0, len(hdkeys))
		for _, k := range hdkeys {
			out.EvmAddrs = append(out.EvmAddrs, k.EthAddr())
			out.AvaAddrs = append(out.AvaAddrs, k.AvaAddr("P", args.Hrp))
		}
		json.NewEncoder(os.Stdout).Encode(out)
		return
	}

	fmtStr := "%-18s %-42s %-45s\n"

	fmt.Printf(fmtStr, "Path", "EVM Addr", "Ava Addr")
	for _, k := range hdkeys {
		fmt.Printf(fmtStr,
			k.Path,
			k.EthAddr(),
			k.AvaAddr("P", args.Hrp),
		)
	}
}
