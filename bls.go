package main

import (
	"encoding/hex"
	"fmt"

	blst "github.com/supranational/blst/bindings/go"
)

// Domain separator for hash-to-curve, specific to BLS signature scheme
const domain = "BLS_SIG_BLS12381G1_XMD:SHA-256_SSWU_RO_POP_"

func main() {
	// Example hex inputs (replace with actual hex strings)
	pkHex := "your_pk_hex_here"   // G2 point, 192 bytes (compressed)
	popHex := "your_pop_hex_here" // G1 point, 96 bytes (compressed)

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
	pk := new(blst.P2Affine).Uncompress(pkBytes)
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

// serializeG1 converts a G1 point to [x, y] as big.Int strings
func serializeG1(point *blst.P1Affine) (string, string) {
	x := point.X.ToBig()
	y := point.Y.ToBig()
	return x.String(), y.String()
}

// serializeG2 converts a G2 point to [x_real, x_imag, y_real, y_imag] as big.Int strings
func serializeG2(point *blst.P2Affine) (string, string, string, string) {
	x := point.X.ToBigArray() // [real, imag]
	y := point.Y.ToBigArray() // [real, imag]
	return x[0].String(), x[1].String(), y[0].String(), y[1].String()
}
