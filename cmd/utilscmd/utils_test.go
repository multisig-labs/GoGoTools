package utilscmd

import (
	"testing"

	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/coreth/plugin/evm"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func Test_AddrVariants(t *testing.T) {
	j, err := addrVariants("P-avax1tnuesf6cqwnjw7fxjyk7lhch0vhf0v95wj5jvy")
	t.Log(j)
	require.NoError(t, err)
	t.Fatal()
}

var factory secp256k1.Factory

func Test_RecoverAddr(t *testing.T) {
	// c := blocks.GenesisCodec
	// txId := "2uuhQuhbFDXZCsFk1AcTr43HYnYbsWBG6x1WHUWLM5uaUYVM8u"
	unsignedTxHex := "0x0000000000110000000100000000000000000000000000000000000000000000000000000000000000000000000121e67317cbc4be2aeb00677ad6462778a8f52274b9d605df2591b23027a87dff000000070000000562ccb8ba000000000000000000000001000000015cf998275803a7277926912defdf177b2e97b0b40000000000000000ed5f38341e436e5d46e2bb00b45d62ae97d1b050c64bc634ae10626739e35c4b0000000167fdba44845e082140152be52304bbb74edd72d5cd6a22b7327258de9bb45efb0000000021e67317cbc4be2aeb00677ad6462778a8f52274b9d605df2591b23027a87dff000000050000000562dbfafa0000000100000000"
	unsignedTxBytes, err := formatting.Decode(formatting.HexNC, unsignedTxHex)
	require.NoError(t, err)

	txHash := hashing.ComputeHash256(unsignedTxBytes)
	sigBytes, err := formatting.Decode(formatting.HexNC, "0x11aeb54f353570a83cf626817eb321ead5ec0ca1be11d95f8342d2df56d50f7e4e9c15435e6dbaef11946a87cbffc9f8a3de0e6becec39bf2de1c6a6ffc41b5501")
	require.NoError(t, err)
	pk, err := factory.RecoverHashPublicKey(txHash, sigBytes[:])
	require.NoError(t, err)
	t.Logf("%+v", pk)
	spew.Dump(pk)
	addr, err := address.FormatBech32("avax", pk.Address().Bytes())
	require.NoError(t, err)
	t.Logf("%v", addr)
	caddy := evm.PublicKeyToEthAddress(pk)
	t.Logf("%s", caddy)
	// common.BytesToAddress(Keccak256(pubBytes[1:])[12:])
	t.Fatal()
}

// {
//   "jsonrpc": "2.0",
//   "result": {
//     "tx": {
//       "unsignedTx": {
//         "networkID": 1,
//         "blockchainID": "11111111111111111111111111111111LpoYY",
//         "outputs": [
//           {
//             "assetID": "FvwEAhmxKfeiG8SnEvq42hc6whRyY3EFYAvebMqDNDGCgxN5Z",
//             "fxID": "spdxUxVJQbX85MGxMHbKw1sHxMnSqJ3QBzDyDYEP3h6TLuxqQ",
//             "output": {
//               "addresses": [
//                 "P-avax1tnuesf6cqwnjw7fxjyk7lhch0vhf0v95wj5jvy"
//               ],
//               "amount": 77849476551,
//               "locktime": 0,
//               "threshold": 1
//             }
//           }
//         ],
//         "inputs": [
//           {
//             "txID": "ao2ggs6wamRcJS68LdVxt95iqWmGp22G4FxiTnmNvMwUQVBNa",
//             "outputIndex": 0,
//             "assetID": "FvwEAhmxKfeiG8SnEvq42hc6whRyY3EFYAvebMqDNDGCgxN5Z",
//             "fxID": "spdxUxVJQbX85MGxMHbKw1sHxMnSqJ3QBzDyDYEP3h6TLuxqQ",
//             "input": {
//               "amount": 106013920254,
//               "signatureIndices": [
//                 0
//               ]
//             }
//           }
//         ],
//         "memo": "0x",
//         "validator": {
//           "nodeID": "NodeID-7AnwdDA9QLTgTwzKoFL7j9wsBvourjD3S",
//           "start": 1692050894,
//           "end": 1693260494,
//           "weight": 28164443703
//         },
//         "stake": [
//           {
//             "assetID": "FvwEAhmxKfeiG8SnEvq42hc6whRyY3EFYAvebMqDNDGCgxN5Z",
//             "fxID": "spdxUxVJQbX85MGxMHbKw1sHxMnSqJ3QBzDyDYEP3h6TLuxqQ",
//             "output": {
//               "addresses": [
//                 "P-avax1tnuesf6cqwnjw7fxjyk7lhch0vhf0v95wj5jvy"
//               ],
//               "amount": 28164443703,
//               "locktime": 0,
//               "threshold": 1
//             }
//           }
//         ],
//         "rewardsOwner": {
//           "addresses": [
//             "P-avax19zfygxaf59stehzedhxjesads0p5jdvfeedal0"
//           ],
//           "locktime": 0,
//           "threshold": 1
//         }
//       },
//       "credentials": [
//         {
//           "signatures": [
//             "0xdb7347e94300275f036d7ef0c8eb7ad91857f7a2a0254780b8339abda051305c39738c85143dda4f9b0614d02a0e4bc8771c1218e8eaf44a1bd72a433d59342100"
//           ]
//         }
//       ],
//       "id": "2QKru99n5skgVjQnpjD2CoXSGMXuf9ge1W5nJjJkfRs5dti2wa"
//     },
//     "encoding": "json"
//   },
//   "id": 0
// }
