package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNrClient_FindEntityGUID(t *testing.T) {
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

	entityGUID, err := client.FindEntityGUID("Metric", "windowsService.service.status", "testKey", "xxxx")
	require.NoError(t, err)
	require.NotEmpty(t, entityGUID)

	entity, err := client.FindEntityByGUID(entityGUID)
	require.NoError(t, err)
	require.NotEmpty(t, entity)
}

func TestNrClient_FindEntityByGUID(t *testing.T) {
	t.Skipped()
}
