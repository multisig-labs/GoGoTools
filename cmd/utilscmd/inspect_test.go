package utilscmd

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ava-labs/avalanchego/vms/platformvm/txs"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func Test_InspectCreateChainTx(t *testing.T) {
	// txb64 := "AAAAAAAPAAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAT2b2sDtHXYTMM9oDv3rGkIVnrOH1tKVDJb30o9hu+KqAAAABwAAAAAFHGJAAAAAAAAAAAAAAAABAAAAAXV7ygFb3qW0cYdejdINXVUUOYvwAAAAAYcXWiAVvQ2pjNCfZi63DN/Xj4FFqVaYbhEyVtMGGgkEAAAAAD2b2sDtHXYTMM9oDv3rGkIVnrOH1tKVDJb30o9hu+KqAAAABQAAAAALEkNAAAAAAQAAAAAAAAAAhxdaIBW9DamM0J9mLrcM39ePgUWpVphuETJW0wYaCQQABkNPUU5ldHN1Ym5ldGV2bQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKAAAAAQAAAAAAAAACAAAACQAAAAFi1fuzBfjXzB50CDMC4Uyv2DbtEFEYAk/YdBQ6eYpT5Ql7JDEdtNdwgSNOVjo0vc8/NzroDF/RQdusshcghgSTAQAAAAkAAAABFBTz7/Z2hoMIbQGItW/dm5GYl2UIbnWgAwkuh2VoKBQL3dHvjOdbl7tbhIu1Hd+JF9MmKxlco/fj5yn6kMMHUQA="
	txb64 := "AAAAAAAPAAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAT2b2sDtHXYTMM9oDv3rGkIVnrOH1tKVDJb30o9hu+KqAAAABwAAAACxyRiAAAAAAAAAAAAAAAABAAAAAXV7ygFb3qW0cYdejdINXVUUOYvwAAAAAiK5NS7GJ/aPuE29LKFzdblMG9EmnpszQyd8Wpzk6OqKAAAAAD2b2sDtHXYTMM9oDv3rGkIVnrOH1tKVDJb30o9hu+KqAAAABQAAAACywRvAAAAAAQAAAACLei+F68dpdk1o8JpEp5JynlYWIEPZnSl0+nPk1hEswAAAAAA9m9rA7R12EzDPaA796xpCFZ6zh9bSlQyW99KPYbviqgAAAAUAAAAABP3dwAAAAAEAAAAAAAAAAIcXWiAVvQ2pjNCfZi63DN/Xj4FFqVaYbhEyVtMGGgkEAAZDT1FOZXRzdWJuZXRldm0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAY0ewogICJjb25maWciOiB7CiAgICAiYnl6YW50aXVtQmxvY2siOiAwLAogICAgImNoYWluSWQiOiA2OTQyMCwKICAgICJjb25zdGFudGlub3BsZUJsb2NrIjogMCwKICAgICJjb250cmFjdE5hdGl2ZU1pbnRlckNvbmZpZyI6IHsKICAgICAgImFkbWluQWRkcmVzc2VzIjogWyIweEUzRTcyZUQwQWFCNzg4MUVENzk1OTMyMkY5NTE1MmJjMjJjNDdmODEiXSwKICAgICAgImJsb2NrVGltZXN0YW1wIjogMAogICAgfSwKICAgICJlaXAxNTBCbG9jayI6IDAsCiAgICAiZWlwMTU1QmxvY2siOiAwLAogICAgImVpcDE1OEJsb2NrIjogMCwKICAgICJmZWVDb25maWciOiB7CiAgICAgICJnYXNMaW1pdCI6IDgwMDAwMDAsCiAgICAgICJ0YXJnZXRCbG9ja1JhdGUiOiAyLAogICAgICAibWluQmFzZUZlZSI6IDI1MDAwMDAwMDAwLAogICAgICAidGFyZ2V0R2FzIjogMTUwMDAwMDAsCiAgICAgICJiYXNlRmVlQ2hhbmdlRGVub21pbmF0b3IiOiAzNiwKICAgICAgIm1pbkJsb2NrR2FzQ29zdCI6IDAsCiAgICAgICJtYXhCbG9ja0dhc0Nvc3QiOiAxMDAwMDAwLAogICAgICAiYmxvY2tHYXNDb3N0U3RlcCI6IDIwMDAwMAogICAgfSwKICAgICJmZWVNYW5hZ2VyQ29uZmlnIjogewogICAgICAiYWRtaW5BZGRyZXNzZXMiOiBbIjB4RTNFNzJlRDBBYUI3ODgxRUQ3OTU5MzIyRjk1MTUyYmMyMmM0N2Y4MSJdLAogICAgICAiYmxvY2tUaW1lc3RhbXAiOiAwCiAgICB9LAogICAgImhvbWVzdGVhZEJsb2NrIjogMCwKICAgICJpc3RhbmJ1bEJsb2NrIjogMCwKICAgICJtdWlyR2xhY2llckJsb2NrIjogMCwKICAgICJwZXRlcnNidXJnQmxvY2siOiAwLAogICAgIndhcnBDb25maWciOiB7CiAgICAgICJibG9ja1RpbWVzdGFtcCI6IDE3MTYzMTU1MjUsCiAgICAgICJxdW9ydW1OdW1lcmF0b3IiOiA2NwogICAgfQogIH0sCiAgIm5vbmNlIjogIjB4MCIsCiAgInRpbWVzdGFtcCI6ICIweDY2NGNlNTg1IiwKICAiZXh0cmFEYXRhIjogIjB4IiwKICAiZ2FzTGltaXQiOiAiMHg3YTEyMDAiLAogICJkaWZmaWN1bHR5IjogIjB4MCIsCiAgIm1peEhhc2giOiAiMHgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwIiwKICAiY29pbmJhc2UiOiAiMHgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwIiwKICAiYWxsb2MiOiB7CiAgICAiRTNFNzJlRDBBYUI3ODgxRUQ3OTU5MzIyRjk1MTUyYmMyMmM0N2Y4MSI6IHsKICAgICAgImJhbGFuY2UiOiAiMHgxMTU4ZTQ2MDkxM2QwMDAwMCIKICAgIH0KICB9LAogICJhaXJkcm9wSGFzaCI6ICIweDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLAogICJhaXJkcm9wQW1vdW50IjogbnVsbCwKICAibnVtYmVyIjogIjB4MCIsCiAgImdhc1VzZWQiOiAiMHgwIiwKICAicGFyZW50SGFzaCI6ICIweDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLAogICJiYXNlRmVlUGVyR2FzIjogbnVsbCwKICAiZXhjZXNzQmxvYkdhcyI6IG51bGwsCiAgImJsb2JHYXNVc2VkIjogbnVsbAp9CgAAAAoAAAABAAAAAAAAAAMAAAAJAAAAAe8vwN6P/HKxb0NTwb9CeIsG8NkvNtrB1mn/sd52U/kNfQWpp0aQ6D9fTwLYxmY6Q+z5Lvih7WRxgZl488zkNo4BAAAACQAAAAHvL8Dej/xysW9DU8G/QniLBvDZLzbawdZp/7HedlP5DX0FqadGkOg/X08C2MZmOkPs+S74oe1kcYGZePPM5DaOAQAAAAkAAAABFvajFN75KuTjwkoG1hfss9XF9bqQP0unKUoJKJ9SIMBy+/z+MncwmEHC43f69/KThWlsTu3w08VE7XBNr7pk+AA="
	txb, _ := b64.StdEncoding.DecodeString(txb64)
	t.Logf("%x", txb)
	txcc := &txs.Tx{}
	_, err := txs.Codec.Unmarshal(txb, txcc)
	require.NoError(t, err)
	t.Logf("%+v", txcc)

	js, err := json.Marshal(txcc)
	require.NoError(t, err)
	t.Logf("%s", js)

	gdata := "ewogICJjb25maWciOiB7CiAgICAiYnl6YW50aXVtQmxvY2siOiAwLAogICAgImNoYWluSWQiOiA2OTQyMCwKICAgICJjb25zdGFudGlub3BsZUJsb2NrIjogMCwKICAgICJjb250cmFjdE5hdGl2ZU1pbnRlckNvbmZpZyI6IHsKICAgICAgImFkbWluQWRkcmVzc2VzIjogWyIweEUzRTcyZUQwQWFCNzg4MUVENzk1OTMyMkY5NTE1MmJjMjJjNDdmODEiXSwKICAgICAgImJsb2NrVGltZXN0YW1wIjogMAogICAgfSwKICAgICJlaXAxNTBCbG9jayI6IDAsCiAgICAiZWlwMTU1QmxvY2siOiAwLAogICAgImVpcDE1OEJsb2NrIjogMCwKICAgICJmZWVDb25maWciOiB7CiAgICAgICJnYXNMaW1pdCI6IDgwMDAwMDAsCiAgICAgICJ0YXJnZXRCbG9ja1JhdGUiOiAyLAogICAgICAibWluQmFzZUZlZSI6IDI1MDAwMDAwMDAwLAogICAgICAidGFyZ2V0R2FzIjogMTUwMDAwMDAsCiAgICAgICJiYXNlRmVlQ2hhbmdlRGVub21pbmF0b3IiOiAzNiwKICAgICAgIm1pbkJsb2NrR2FzQ29zdCI6IDAsCiAgICAgICJtYXhCbG9ja0dhc0Nvc3QiOiAxMDAwMDAwLAogICAgICAiYmxvY2tHYXNDb3N0U3RlcCI6IDIwMDAwMAogICAgfSwKICAgICJmZWVNYW5hZ2VyQ29uZmlnIjogewogICAgICAiYWRtaW5BZGRyZXNzZXMiOiBbIjB4RTNFNzJlRDBBYUI3ODgxRUQ3OTU5MzIyRjk1MTUyYmMyMmM0N2Y4MSJdLAogICAgICAiYmxvY2tUaW1lc3RhbXAiOiAwCiAgICB9LAogICAgImhvbWVzdGVhZEJsb2NrIjogMCwKICAgICJpc3RhbmJ1bEJsb2NrIjogMCwKICAgICJtdWlyR2xhY2llckJsb2NrIjogMCwKICAgICJwZXRlcnNidXJnQmxvY2siOiAwLAogICAgIndhcnBDb25maWciOiB7CiAgICAgICJibG9ja1RpbWVzdGFtcCI6IDE3MTYzMTU1MjUsCiAgICAgICJxdW9ydW1OdW1lcmF0b3IiOiA2NwogICAgfQogIH0sCiAgIm5vbmNlIjogIjB4MCIsCiAgInRpbWVzdGFtcCI6ICIweDY2NGNlNTg1IiwKICAiZXh0cmFEYXRhIjogIjB4IiwKICAiZ2FzTGltaXQiOiAiMHg3YTEyMDAiLAogICJkaWZmaWN1bHR5IjogIjB4MCIsCiAgIm1peEhhc2giOiAiMHgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwIiwKICAiY29pbmJhc2UiOiAiMHgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwIiwKICAiYWxsb2MiOiB7CiAgICAiRTNFNzJlRDBBYUI3ODgxRUQ3OTU5MzIyRjk1MTUyYmMyMmM0N2Y4MSI6IHsKICAgICAgImJhbGFuY2UiOiAiMHgxMTU4ZTQ2MDkxM2QwMDAwMCIKICAgIH0KICB9LAogICJhaXJkcm9wSGFzaCI6ICIweDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLAogICJhaXJkcm9wQW1vdW50IjogbnVsbCwKICAibnVtYmVyIjogIjB4MCIsCiAgImdhc1VzZWQiOiAiMHgwIiwKICAicGFyZW50SGFzaCI6ICIweDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLAogICJiYXNlRmVlUGVyR2FzIjogbnVsbCwKICAiZXhjZXNzQmxvYkdhcyI6IG51bGwsCiAgImJsb2JHYXNVc2VkIjogbnVsbAp9Cg=="
	gdatab, err := b64.StdEncoding.DecodeString(gdata)
	require.NoError(t, err)
	t.Logf("%s", gdatab)
	g := new(core.Genesis)
	err = json.Unmarshal(gdatab, g)
	require.NoError(t, err)
	t.Logf("%+v", g)

	t.Fatal()
}

func Test_Foo(t *testing.T) {
	expectedUnsignedSimpleBaseTxBytes := []byte{
		// Codec version
		0x00, 0x00,
		// BaseTx Type ID
		0x00, 0x00, 0x00, 0x22,
		// Mainnet network ID
		0x00, 0x00, 0x00, 0x01,
		// P-chain blockchain ID
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// Number of outputs
		0x00, 0x00, 0x00, 0x00,
		// Number of inputs
		0x00, 0x00, 0x00, 0x01,
		// Inputs[0]
		// TxID
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		// Tx output index
		0x00, 0x00, 0x00, 0x01,
		// Mainnet AVAX assetID
		0x21, 0xe6, 0x73, 0x17, 0xcb, 0xc4, 0xbe, 0x2a,
		0xeb, 0x00, 0x67, 0x7a, 0xd6, 0x46, 0x27, 0x78,
		0xa8, 0xf5, 0x22, 0x74, 0xb9, 0xd6, 0x05, 0xdf,
		0x25, 0x91, 0xb2, 0x30, 0x27, 0xa8, 0x7d, 0xff,
		// secp256k1fx transfer input type ID
		0x00, 0x00, 0x00, 0x05,
		// input amount = 1 MilliAvax
		0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, 0x42, 0x40,
		// number of signatures needed in input
		0x00, 0x00, 0x00, 0x01,
		// index of signer
		0x00, 0x00, 0x00, 0x05,
		// length of memo
		0x00, 0x00, 0x00, 0x00,
	}
	fmt.Printf("%x", expectedUnsignedSimpleBaseTxBytes)
	t.Fatal()

}

func formatBech32(hrp string, payload []byte) (string, error) {
	fiveBits, err := bech32.ConvertBits(payload, 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.Encode(hrp, fiveBits)
}

func TestFrak(t *testing.T) {
	_, z, err := bech32.Decode("abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw")
	require.NoError(t, err)
	t.Logf("%x", z)
	a := "0xE7CDF6BD6AB2720A818C1CBB09B67C55D062DFB3"
	b, err := hexutil.Decode(a)
	require.NoError(t, err)
	fiveBits, err := bech32.ConvertBits(b, 8, 5, true)
	require.NoError(t, err)
	t.Logf("%x", fiveBits)
	t.Fatal()
}
