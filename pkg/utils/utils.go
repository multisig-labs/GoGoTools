package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// Given nodeID in either 0x123 or NodeID-123 format, return the [20]byte ids.NodeID format
func ParseNodeID(nodeID string) (ids.NodeID, error) {
	var nodeShortID ids.NodeID
	var err error

	if strings.HasPrefix(nodeID, "NodeID-") {
		nodeShortID, err = ids.NodeIDFromString(nodeID)
		if err != nil {
			return ids.NodeID{}, fmt.Errorf("error decoding nodeID %s: %w", nodeID, err)
		}
		return nodeShortID, nil
	}

	if strings.HasPrefix(nodeID, "0x") {
		b := HexToBytes(nodeID)
		b20 := common.BytesToAddress(b)
		return ids.NodeID(b20), nil
	}

	return ids.NodeID{}, fmt.Errorf("invalid nodeID format %s: %w", nodeID, err)
}

// Given nodeID in [20]bytes address format, return the [20]byte ids.NodeID format
func AddressToNodeID(nodeID common.Address) ids.NodeID {
	return ids.NodeID(nodeID)
}

// returns the bytes represented by the hexadecimal string s, may be prefixed with "0x".
func HexToBytes(s string) []byte {
	b, _ := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	return b
}

// returns the hexadecimal string representation of the bytes b, prefixed with "0x".
func BytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

// Converts a '0x'-prefixed hex string or cb58-encoded string to an ID.
func HexOrCB58ToID(s string) (ids.ID, error) {
	if strings.HasPrefix(s, "0x") {
		bytes, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
		if err != nil {
			return ids.ID{}, err
		}
		return ids.ToID(bytes)
	}
	return ids.FromString(s)
}

// Support PrivateKey-(cb58) or hex string with optional 0x prefix
func ParsePrivateKey(pkStr string) (avaKey *secp256k1.PrivateKey, ethKey *ecdsa.PrivateKey, err error) {
	var pkBytes []byte
	if strings.HasPrefix(pkStr, "PrivateKey-") {
		pkBytes, err = cb58.Decode(strings.TrimPrefix(pkStr, "PrivateKey-"))
		if err != nil {
			return nil, nil, err
		}
	} else {
		pkBytes, err = hex.DecodeString(strings.TrimPrefix(pkStr, "0x"))
		if err != nil {
			return nil, nil, err
		}
	}
	avaKey, err = secp256k1.ToPrivateKey(pkBytes)
	if err != nil {
		return nil, nil, err
	}
	ethKey, err = crypto.HexToECDSA(fmt.Sprintf("%x", pkBytes))
	if err != nil {
		return nil, nil, err
	}
	return avaKey, ethKey, nil
}

// Parse a private key and return the ava and eth addresses
// network is the network name (mainnet, fuji, etc)
func ParsePrivateKeyToAddresses(privateKeyStr string, network string) (string, string, error) {
	pchainKey, ethKey, err := ParsePrivateKey(privateKeyStr)
	if err != nil {
		return "", "", err
	}
	avaAddr, err := address.Format("P", network, pchainKey.PublicKey().Address().Bytes())
	if err != nil {
		return "", "", err
	}
	ethAddr := crypto.PubkeyToAddress(ethKey.PublicKey)

	return avaAddr, ethAddr.String(), nil
}

func ValidateBLSKeys(blsPubKey string, blsPop string) error {
	publicKey, err := bls.PublicKeyFromCompressedBytes(HexToBytes(blsPubKey))
	if err != nil {
		return fmt.Errorf("error decoding blsPubKey %s: %w", blsPubKey, err)
	}
	signature, err := bls.SignatureFromBytes(HexToBytes(blsPop))
	if err != nil {
		return fmt.Errorf("error decoding blsPop %s: %w", blsPop, err)
	}
	if !bls.VerifyProofOfPossession(publicKey, signature, HexToBytes(blsPubKey)) {
		return fmt.Errorf("invalid proof of possession")
	}

	return nil
}

func ConvertNanoAvaxToWei(nanoAvax int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(nanoAvax), big.NewInt(1e9))
}

func ConvertAvaxToWei(avax float64) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(avax*1e9)), big.NewInt(1e9))
}

func FetchRPCGJSON(url string, method string, params string) (*gjson.Result, error) {
	s, err := FetchRPC(url, method, params)
	if err != nil {
		return nil, err
	}
	out := gjson.Parse(s)
	return &out, nil
}

func FetchRPC(url string, method string, params string) (string, error) {
	client := resty.New()
	// client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(30 * time.Second)

	var resp *resty.Response
	var err error

	if params == "" {
		params = "{}"
	}

	body := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id"     : %d,
		"method" : "%s",
		"params" : %s
	}`, time.Now().Unix(), method, params)

	resp, err = client.R().
		SetDebug(IsDebug()).
		EnableTrace().
		EnableGenerateCurlOnDebug().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(body).
		Post(url)

	if resp.IsError() {
		return "", fmt.Errorf("fetch error %d: %s %s", resp.StatusCode(), url, body)
	}
	return resp.String(), err
}

func IsDebug() bool {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		return false
	}
	return debug
}

// DecodeError decodes an error from an ABI string and an error.
// usage: err = DecodeError(abiStr, err)
func DecodeError(abiStr string, err error) error {
	if err == nil {
		return nil
	}
	parsedABI, _ := abi.JSON(strings.NewReader(abiStr))
	// Try to decode the revert reason using the ABI
	if revertErr, ok := err.(interface{ ErrorData() interface{} }); ok {
		if data := revertErr.ErrorData(); data != nil {
			// Get the raw error data
			errData := data.(string)
			// Convert hex string to bytes
			if errBytes, hexErr := hex.DecodeString(strings.TrimPrefix(errData, "0x")); hexErr == nil {
				var errBytes4 [4]byte
				copy(errBytes4[:], errBytes[:4])
				if abiError, findErr := parsedABI.ErrorByID(errBytes4); findErr == nil {
					// If there's no data to unpack (len == 4 for just the selector)
					if len(errBytes) == 4 {
						return fmt.Errorf("transaction reverted: %v (decoded error: %s)",
							err, abiError.Name)
					}
					// Try to unpack data if available
					if errorData, unpackErr := abiError.Unpack(errBytes[4:]); unpackErr == nil {
						return fmt.Errorf("transaction reverted: %v (decoded error: %s%v)",
							err, abiError.Name, errorData)
					}
				}
			}
		}
	}
	return err
}

// Better JSON output for RegisterL1Validator
func RegisterL1ValidatorToJSON(m *message.RegisterL1Validator) ([]byte, error) {

	type pChainOwner struct {
		Threshold uint32   `serialize:"true" json:"threshold"`
		Addresses []string `serialize:"true" json:"addresses"`
	}

	rbOwners := []string{}
	dbOwners := []string{}
	for _, owner := range m.RemainingBalanceOwner.Addresses {
		addr, _ := address.Format("P", "avax", owner.Bytes())
		rbOwners = append(rbOwners, addr)
	}
	for _, owner := range m.DisableOwner.Addresses {
		addr, _ := address.Format("P", "avax", owner.Bytes())
		dbOwners = append(dbOwners, addr)
	}

	rv := struct {
		SubnetID              ids.ID      `serialize:"true" json:"subnetID"`
		NodeID                ids.NodeID  `serialize:"true" json:"nodeID"`
		BLSPublicKey          string      `serialize:"true" json:"blsPublicKey"`
		Expiry                uint64      `serialize:"true" json:"expiry"`
		RemainingBalanceOwner pChainOwner `serialize:"true" json:"remainingBalanceOwner"`
		DisableOwner          pChainOwner `serialize:"true" json:"disableOwner"`
		Weight                uint64      `serialize:"true" json:"weight"`
	}{
		SubnetID:     m.SubnetID,
		NodeID:       ids.NodeID(m.NodeID),
		BLSPublicKey: hex.EncodeToString(m.BLSPublicKey[:]),
		Expiry:       m.Expiry,
		RemainingBalanceOwner: pChainOwner{
			Threshold: m.RemainingBalanceOwner.Threshold,
			Addresses: rbOwners,
		},
		DisableOwner: pChainOwner{
			Threshold: m.DisableOwner.Threshold,
			Addresses: dbOwners,
		},
		Weight: m.Weight,
	}

	jsonBytes, err := json.MarshalIndent(rv, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
