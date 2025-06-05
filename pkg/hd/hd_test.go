package hd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHD(t *testing.T) {
	mnemonic := "test test test test test test test test test test test junk"

	// Test public key derivation
	pubkey, err := XPubKey(mnemonic, EthDerivationPath)
	require.NoError(t, err)
	t.Logf("pubkey: %s", pubkey.String())

	hdkeys, err := DerivePubKeys(pubkey, EthDerivationPath, 10)
	require.NoError(t, err)

	// Test private key derivation for comparison
	privateKeys, err := DeriveHDKeys(mnemonic, EthDerivationPath, 10)
	require.NoError(t, err)

	// Verify that public key derivation matches private key derivation
	require.Equal(t, len(privateKeys), len(hdkeys), "Number of keys should match")

	for i, pubKeyResult := range hdkeys {
		privateKeyResult := privateKeys[i]
		t.Logf("Index %d - Public key addr: %s, Private key addr: %s",
			i, pubKeyResult.EthAddr(), privateKeyResult.EthAddr())

		// The addresses should match
		require.Equal(t, privateKeyResult.EthAddr(), pubKeyResult.EthAddr(),
			"Address at index %d should match between public and private key derivation", i)
	}
}
