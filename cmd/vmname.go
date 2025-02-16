package main

import (
	"fmt"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/jxskiss/mcli"
)

func vmNameCmd() {
	args := struct {
		ID string `cli:"#R, id, VM ID"`
	}{}
	mcli.MustParse(&args, mcli.DisableGlobalFlags())

	id, err := ids.FromString(args.ID)
	checkErr(err)
	out := strings.Builder{}
	for _, v := range id {
		out.Write([]byte{v})
	}
	fmt.Println(out.String())
}

func vmIDCmd() {
	args := struct {
		Name string `cli:"#R, name, VM Name"`
	}{}
	mcli.MustParse(&args, mcli.DisableGlobalFlags())

	nameBytes := []byte(args.Name)
	if len(nameBytes) > 32 {
		fmt.Printf("Warning: VM name exceeds 32 bytes, will be truncated\n")
	}

	paddedBytes := [32]byte{}
	copy(paddedBytes[:], nameBytes)
	id, err := ids.ToID(paddedBytes[:])
	checkErr(err)
	fmt.Println(id.String())
}
