package sigagg

import (
	"testing"

	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/multisig-labs/gogotools/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	c, err := NewClient("https://glacier-api.avax.network/v1/signatureAggregator/mainnet/aggregateSignatures", map[string]string{"x-glacier-api-key": ""})
	require.NoError(t, err)

	data := utils.HexToBytes("0x00000000000100000000000000000000000000000000000000000000000000000000000000000000003500000000000100000000000000270000000000023d396f65226fb525e6820be86fd18286a70b033765dc1a3f7e2c09cb36d3982001")
	msg, err := warp.ParseUnsignedMessage(data)
	require.NoError(t, err)

	msgSigned, err := c.AggregateSignatures(msg, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, msgSigned)
}
