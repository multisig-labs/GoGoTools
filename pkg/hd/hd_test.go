package hd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHD(t *testing.T) {
	mnemonic := "test test test test test test test test test test test junk"
	pubkey, err := XPubKey(mnemonic, EthDerivationPath)
	require.NoError(t, err)
	t.Logf("pubkey: %s", pubkey.String())
	hdkeys, err := DerivePubKeys(pubkey, EthDerivationPath, 10)
	require.NoError(t, err)
	for _, k := range hdkeys {
		t.Logf("addr: %s", k.EthAddr())
	}
	t.Fail()
}
