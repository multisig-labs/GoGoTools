package uptime

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetCurrentValidators(t *testing.T) {
	client, err := NewClient("https://blah")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	validators, err := client.GetCurrentValidators(ctx)
	require.NoError(t, err)

	validatorsJSON, err := json.MarshalIndent(validators, "", "  ")
	require.NoError(t, err)
	fmt.Println(string(validatorsJSON))
	t.Fail()
}
