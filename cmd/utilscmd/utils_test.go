package utilscmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AddrVariants(t *testing.T) {
	j, err := addrVariants("X-fuji1wycv8n7d2fg9aq6unp23pnj4q0arv03ysya8jw")
	t.Log(j)
	require.NoError(t, err)
	// t.Fatal()
}
