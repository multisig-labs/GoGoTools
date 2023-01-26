package utilscmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
)

// Same algo as Avax wallet
// msg is the message, returns the hash of the full msg with prefix
func digestAvaMsg(msg string) []byte {
	msgb := []byte(msg)
	l := uint32(len(msgb))
	lb := make([]byte, 4)
	binary.BigEndian.PutUint32(lb, l)
	prefix := []byte("\x1AAvalanche Signed Message:\n")

	buf := new(bytes.Buffer)
	buf.Write(prefix)
	buf.Write(lb)
	buf.Write(msgb)
	fullmsg := buf.Bytes()
	h := sha256.Sum256(fullmsg)
	return h[:]
}

func newMsgDigestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "msgdigest [msg]",
		Short: "Generate a hash digest of a message",
		Long: `Construct an Avalanche Signed Message and return the hash, that can then
be signed and verified in the Avalanche web wallet.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := digestAvaMsg(args[0])
			h := hex.EncodeToString(b)
			fmt.Printf("0x%s\n", h)
			return nil
		},
	}
	return cmd
}
