package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNrClient_FindEntityByGUID(t *testing.T) {
	t.Skipped()
	const apiKey = "XXXXXXXX"
	const accountID = 2762945

	apiClient, err := NewApiClientWrapper(apiKey)
	require.NoError(t, err)

	client := nrClient{
		accountID,
		apiKey,
		apiClient,
	}

	a, err := client.FindEntityGUID("Metric", "haproxy.frontend.connectionsPerSecond", "")
	require.NoError(t, err)
	require.NotEmpty(t, a)
}

func TestNrClient_FindEntityGUID(t *testing.T) {
	t.Skipped()
}
