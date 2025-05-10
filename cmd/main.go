package main

import (
	"fmt"
	"os"

	"github.com/jxskiss/mcli"
	"github.com/multisig-labs/gogotools/pkg/version"
)

type URLFlags struct {
	AvaUrl string `cli:"--ava-url, URL of the Avalanche node (do not include a path)." env:"AVA_RPC_URL" default:"https://api.avax.network"`
	EthUrl string `cli:"--eth-url, URL of the Ethereum endpoint (full path to evm rpc)" env:"ETH_RPC_URL" default:"https://api.avax.network/ext/bc/C/rpc"`
}

func main() {
	defer version.PanicHandler()
	// mcli.SetOptions(mcli.Options{KeepCommandOrder: true})
	mcli.Add("balance", balanceAddressCmd, "Get the balance of an address")
	mcli.Add("balance-pk", balancePKCmd, "Get the balance of an address for a private key")
	mcli.Add("bech32-decode", bech32DecodeCmd, "Decode a bech32 address")
	mcli.Add("cb58-encode", cb58EncodeCmd, "Encode a value to CB58")
	mcli.Add("cb58-decode", cb58DecodeCmd, "Decode a value from CB58 (ID or NodeID)")
	mcli.Add("cb58-decodesig", cb58DecodeSigCmd, "Decode a signature (r,s,v) from CB58")
	mcli.Add("cross-chain-tx", crossChainTransferCmd, "Transfer assets from C-Chain to P-Chain")
	mcli.Add("msgdigest", digestAvaMsgCmd, "Generate a hash digest for an Avalanche Signed Message (ERC-191)")
	mcli.Add("pk", privkeyCmd, "Show various address encodings of a private key")
	mcli.Add("rpc", rpcCmd, "Ergonomic access to avalanche node RPC APIs")
	mcli.Add("l1-validators", l1ValidatorsCmd, "Get current validators from a L1 validator RPC endpoint")
	mcli.Add("mnemonic-addrs", mnemonicAddrsCmd, "Show addresses for a BIP39 mnemonic")
	mcli.Add("mnemonic-keys", mnemonicKeysCmd, "Show keys and addresses for a BIP39 mnemonic")
	mcli.Add("mnemonic-generate", randomMnemonicCmd, "Generate a random BIP39 mnemonic")
	mcli.Add("mnemonic-insecure", mnemonicInsecureCmd, "Generate an INSECURE test BIP39 mnemonic")
	mcli.Add("random-nodeid", randomNodeIDCmd, "Generate a random node ID")
	mcli.Add("random-id", randomIDCmd, "Generate a random ID")
	mcli.Add("random-bls", randomBLSCmd, "Generate a random BLS key")
	mcli.Add("revert-reason", revertReasonCmd, "Find revert reason for a failed tx hash")
	mcli.Add("inspect-tx-p", inspectPTxCmd, "Inspect a P-Chain transaction")
	mcli.Add("vmname", vmNameCmd, "Given a vmID, try to decode the ASCII name")
	mcli.Add("vmid", vmIDCmd, "Given a vmName, try to encode the ASCII name as a vmID")
	mcli.Add("verify-bls", verifyBLSCmd, "Verify a BLS Proof of Possession")
	mcli.Add("warp-get", getWarpMsgCmd, "Get a warp message from a transaction ID")
	mcli.Add("warp-parse", parseWarpMsgCmd, "Parse a warp message")
	mcli.Add("warp-construct-uptime", constructUptimeMsgCmd, "Construct an unsgined uptime message")
	mcli.Add("warp-construct-l1-validator-registration", constructL1ValidatorRegistrationMsgCmd, "Construct an unsigned L1ValidatorRegistration msg")
	mcli.Add("warp-construct-l1-weight", constructL1WeightMsgCmd, "Construct an unsigned msg to change weight on P-Chain")
	mcli.Add("warp-aggregate-signatures", aggregateSignaturesCmd, "Aggregate signatures for a warp message")
	mcli.Add("version", versionCmd, "Version")
	mcli.AddHelp()
	mcli.AddCompletion()
	mcli.Run()
}

func versionCmd() {
	fmt.Printf("Version: %s  BuildDate: %s  Commit: %s\n", version.Version, version.BuildDate, version.GitCommit)
}

func checkErr(err interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
