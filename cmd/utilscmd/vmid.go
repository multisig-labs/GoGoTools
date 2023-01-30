package utilscmd

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/spf13/cobra"
)

// var errIllegalNameCharacter = errors.New(
// 	"illegal name character: only letters/numbers/spaces, no special characters allowed")

// func checkInvalidNames(name string) error {
// 	if len(name) > 32 {
// 		return fmt.Errorf("name must be <= 32 bytes, found %d", len(name))
// 	}
// 	// this is currently exactly the same code as in avalanchego/vms/platformvm/create_chain_tx.go
// 	for _, r := range name {
// 		if r > unicode.MaxASCII || !(unicode.IsLetter(r) || unicode.IsNumber(r) || r == ' ') {
// 			return errIllegalNameCharacter
// 		}
// 	}

// 	return nil
// }

func newVMIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vmid [name]",
		Short: "Generate a valid ID from a name string (max 32 chars)",
		Long:  `An ID is a Base58 + Checksum encode of a 32 byte, zero-padded ASCII string`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// err := checkInvalidNames(args[0])
			// if err != nil {
			// 	return err
			// }
			paddedBytes := [32]byte{}
			copy(paddedBytes[:], []byte(args[0]))
			id, err := ids.ToID(paddedBytes[:])
			if err != nil {
				return err
			}
			fmt.Println(id.String())
			return nil
		},
	}
	return cmd
}
