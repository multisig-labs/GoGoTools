package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/defiweb/go-eth/crypto"
	"github.com/stretchr/testify/require"
	blst "github.com/supranational/blst/bindings/go"
)

func TestVerifyBLSKeys(t *testing.T) {
	pubkey := "0x8d543b279b9bd69c5b6754a09bce1eab2de2d9135eff7e391e42583fca4c19c6007e864971c2baba777dfa312ca7994e"
	pop := "0xa99743e050b543f2482c1010e2908b848c5894080a5e5ac9db96111b721d753efbe106eb12496b56be14f11feaeb4d9605507ca0cb726b3832c45448603025de7ba425d7c054f5e664922d6a5dfab8b9f7df681941d25be420ea4f973b79d041"

	err := ValidateBLSKeys(pubkey, pop)
	require.NoError(t, err)
}

func RecoverPublicKeyP256(hash []byte, r, s *big.Int, recoveryID byte) (*ecdsa.PublicKey, error) {
	// Ensure recoveryID is valid (0 or 1)
	if recoveryID > 1 {
		return nil, errors.New("invalid recovery ID: must be 0 or 1")
	}

	// Get the P-256 curve
	curve := elliptic.P256()

	// Check that r and s are within valid range
	if r.Cmp(curve.Params().N) >= 0 || s.Cmp(curve.Params().N) >= 0 {
		return nil, errors.New("invalid signature: r or s out of range")
	}

	// Step 1: Recover X coordinate from r and recoveryID
	x := new(big.Int).Set(r)
	if recoveryID&1 == 1 {
		// If recoveryID is odd, add curve order to x
		x.Add(x, curve.Params().N)
	}

	// Step 2: Compute Y coordinate
	// y² = x³ - 3x + b (mod p)
	x3 := new(big.Int).Mul(x, x)                 // x²
	x3.Mul(x3, x)                                // x³
	threeX := new(big.Int).Mul(x, big.NewInt(3)) // 3x
	y2 := new(big.Int).Sub(x3, threeX)           // x³ - 3x
	y2.Add(y2, curve.Params().B)                 // x³ - 3x + b
	y2.Mod(y2, curve.Params().P)                 // mod p

	// Compute square root to get y
	y := new(big.Int).ModSqrt(y2, curve.Params().P)
	if y == nil {
		return nil, errors.New("invalid signature: failed to compute y coordinate")
	}

	// If recoveryID bit 0 doesn't match y's parity, use the other root
	if (y.Bit(0) == 0) != (recoveryID&1 == 0) {
		y.Sub(curve.Params().P, y)
	}

	// Verify the point is on the curve
	if !curve.IsOnCurve(x, y) {
		return nil, errors.New("recovered point is not on curve")
	}

	// Step 3: Verify the signature
	pub := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	// Convert hash to correct length for ecdsa.Verify
	if len(hash) > 32 {
		hash = hash[:32]
	}

	if !ecdsa.Verify(pub, hash, r, s) {
		return nil, errors.New("signature verification failed")
	}

	return pub, nil
}

func RecoverAddressP256(hash []byte, r, s *big.Int, recoveryID byte) ([20]byte, error) {
	pub, err := RecoverPublicKeyP256(hash, r, s, recoveryID)
	if err != nil {
		return [20]byte{}, err
	}

	// Convert to Ethereum-style address
	pubBytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	hashPub := crypto.Keccak256(pubBytes[1:]) // Skip the 0x04 prefix
	var addr [20]byte
	copy(addr[:], hashPub[12:]) // Last 20 bytes
	return addr, nil
}

func TestFoo(t *testing.T) {
	// Example usage (you'd need actual values here)
	hash := []byte("example hash to sign")
	r := new(big.Int) // Set your r value
	s := new(big.Int) // Set your s value
	recoveryID := byte(0)

	pubKey, err := RecoverPublicKeyP256(hash, r, s, recoveryID)
	require.NoError(t, err)
	fmt.Println(pubKey.X, pubKey.Y)
}

// Domain separator for hash-to-curve, specific to BLS signature scheme

const (
	BLST_FP_BYTES  = 48 // Size of Fp in bytes
	BLST_FP2_BYTES = 96 // Size of Fp2 in bytes (2 * 48)
)

func TestBLS(t *testing.T) {
	domain := "BLS_SIG_BLS12381G1_XMD:SHA-256_SSWU_RO_POP_"
	// Example hex inputs (replace with actual hex strings)
	pkHex := "0x8d543b279b9bd69c5b6754a09bce1eab2de2d9135eff7e391e42583fca4c19c6007e864971c2baba777dfa312ca7994e"                                                                                                  // G2 point, 192 bytes (compressed)
	popHex := "0xa99743e050b543f2482c1010e2908b848c5894080a5e5ac9db96111b721d753efbe106eb12496b56be14f11feaeb4d9605507ca0cb726b3832c45448603025de7ba425d7c054f5e664922d6a5dfab8b9f7df681941d25be420ea4f973b79d041" // G1 point, 96 bytes (compressed)

	// Decode hex inputs into byte slices
	pkBytes, err := hex.DecodeString(pkHex)
	if err != nil {
		panic("Invalid pk hex: " + err.Error())
	}
	popBytes, err := hex.DecodeString(popHex)
	if err != nil {
		panic("Invalid pop hex: " + err.Error())
	}

	// Deserialize pk into a G2 point
	pk := new(blst.P2Affine)
	pk.Uncompress(pkBytes)
	if pk == nil {
		panic("Invalid pk: failed to deserialize")
	}

	// Deserialize pop into a G1 point
	pop := new(blst.P1Affine).Uncompress(popBytes)
	if pop == nil {
		panic("Invalid pop: failed to deserialize")
	}

	// Compute H(pk): hash the public key to a G1 point
	pkSerialized := pk.Compress() // Serialize pk to bytes for hashing
	H_pk := new(blst.P1Affine).HashTo(pkSerialized, domain)

	// Serialize points for Solidity
	H_pk_x, H_pk_y := serializeG1(H_pk)                           // H_m for Solidity
	pop_x, pop_y := serializeG1(pop)                              // sigma for Solidity
	pk_x_real, pk_x_imag, pk_y_real, pk_y_imag := serializeG2(pk) // pk for Solidity

	// Output the values in Solidity-compatible format
	fmt.Println("H_m (H(pk)):")
	fmt.Printf("  x: %s\n  y: %s\n", H_pk_x, H_pk_y)
	fmt.Println("sigma (pop):")
	fmt.Printf("  x: %s\n  y: %s\n", pop_x, pop_y)
	fmt.Println("pk:")
	fmt.Printf("  x_real: %s\n  x_imag: %s\n  y_real: %s\n  y_imag: %s\n",
		pk_x_real, pk_x_imag, pk_y_real, pk_y_imag)
}

func bytesToBigInt(bytes []byte) *big.Int {
	return new(big.Int).SetBytes(bytes)
}

// serializeG1 extracts x and y coordinates from a G1 point as *big.Int
func serializeG1(point *blst.P1Affine) (x, y *big.Int) {
	// Serialize returns 96 bytes: 48 bytes for x, 48 bytes for y
	serialized := point.Serialize()
	if len(serialized) != 96 {
		panic("Invalid serialized G1 point length")
	}
	xBytes := serialized[0:48]  // First 48 bytes are x
	yBytes := serialized[48:96] // Next 48 bytes are y
	return bytesToBigInt(xBytes), bytesToBigInt(yBytes)
}

// serializeG2 extracts x and y coordinates (real and imag parts) from a G2 point as *big.Int
func serializeG2(point *blst.P2Affine) (xReal, xImag, yReal, yImag *big.Int) {
	// Serialize returns 192 bytes: 96 bytes for x (real + imag), 96 bytes for y (real + imag)
	serialized := point.Serialize()
	if len(serialized) != 192 {
		panic("Invalid serialized G2 point length")
	}
	xRealBytes := serialized[0:48]    // x real part
	xImagBytes := serialized[48:96]   // x imaginary part
	yRealBytes := serialized[96:144]  // y real part
	yImagBytes := serialized[144:192] // y imaginary part
	return bytesToBigInt(xRealBytes), bytesToBigInt(xImagBytes),
		bytesToBigInt(yRealBytes), bytesToBigInt(yImagBytes)
}
