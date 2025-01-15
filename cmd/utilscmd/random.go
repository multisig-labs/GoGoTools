package utilscmd

import (
	"crypto/rand"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
)

func newRandomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "random",
		Short: "Generate random ids of various types",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmd.AddCommand(newRandomNodeIDCmd())
	cmd.AddCommand(newRandomBLSCmd())

	return cmd
}

func newRandomNodeIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node_id",
		Short: "Create random NodeID with 10 leading zero bytes so we know its a fake nodeID",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			r := make([]byte, 20)
			_, err := rand.Read(r)
			cobra.CheckErr(err)
			zeroSlice := make([]byte, 10)
			copy(r, zeroSlice)
			nodeid := ids.NodeID(r)
			fmt.Println(nodeid)
		},
	}
	return cmd
}

func newRandomBLSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bls",
		Short: "Create random bls keys",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			pop := &signer.ProofOfPossession{}
			sk, err := bls.NewSigner()
			cobra.CheckErr(err)
			pop = signer.NewProofOfPossession(sk)
			err = pop.Verify()
			cobra.CheckErr(err)
			popjs, err := pop.MarshalJSON()
			cobra.CheckErr(err)

			skBytes := fmt.Sprintf("0x%x", sk.ToBytes())

			out, err := sjson.SetBytes(popjs, "privateKey", skBytes)
			cobra.CheckErr(err)

			fmt.Println(string(out))
		},
	}
	return cmd
}
