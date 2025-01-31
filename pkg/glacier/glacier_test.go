package glacier

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/stretchr/testify/require"
)

func TestFetchValidators(t *testing.T) {
	subnetID, err := ids.FromString("4YurNFwLzhGUrYyihDnUUc2L199YBnFeWP3fhJKmDDjkbvy8G")
	require.NoError(t, err)
	client := NewClientWithConfig(Config{
		Network: "fuji",
	})
	validators, err := client.FetchValidators(subnetID)
	require.NoError(t, err)
	validatorsJSON, err := json.MarshalIndent(validators, "", "  ")
	require.NoError(t, err)
	fmt.Println(string(validatorsJSON))
	t.Fail()
}
