package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ResolveAmounts(t *testing.T) {
	in := []string{"wootether", "10000", "1ether", "0.1ether", "23948ether"}
	out := ResolveAmounts(in)
	require.ElementsMatch(t, out, []string{"wootether", "10000", "1000000000000000000", "100000000000000000", "23948000000000000000000"})
}

func Test_DownloadAva(t *testing.T) {
	url, destFile, err := DownloadAvalanchego("/tmp", "v1.9.7")
	t.Log(url)
	t.Log(destFile)
	require.NoError(t, err)
}

func Test_DownloadSubnetevm(t *testing.T) {
	url, destFile, err := DownloadSubnetevm("/tmp", "v0.4.7")
	t.Log(url)
	t.Log(destFile)
	require.NoError(t, err)
}
