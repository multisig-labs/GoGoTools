package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/jxskiss/mcli"
)

// Same algo as Avax wallet
// msg is the message, returns the hash of the full msg with prefix
func digestAvaMsgCmd() {
	args := struct {
		Msg string `cli:"#R, msg, Message to digest"`
	}{}
	mcli.MustParse(&args)

	msgb := []byte(args.Msg)
	l := uint32(len(msgb))
	lb := make([]byte, 4)
	binary.BigEndian.PutUint32(lb, l)
	prefix := []byte("\x1AAvalanche Signed Message:\n")

	buf := new(bytes.Buffer)
	buf.Write(prefix)
	buf.Write(lb)
	buf.Write(msgb)
	fullmsg := buf.Bytes()
	hash := sha256.Sum256(fullmsg)
	h := hex.EncodeToString(hash[:])
	fmt.Printf("0x%s\n", h)
}
