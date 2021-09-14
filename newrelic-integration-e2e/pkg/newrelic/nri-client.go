package newrelic

import (
	"fmt"
	"log"

	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

type DataClient interface {
	FindEntityGUID(sample string, metricName string, entityTag string) (*entities.EntityGUID, error)
	FindEntityByGUID(guid *entities.EntityGUID) (entities.EntityInterface, error)
}

type nrClient struct {
	accountID int
	apiKey    string
	client    ApiClient
}

func NewNrClient(apiKey string, accountID int) *nrClient {

	client, err := NewApiClientWrapper(apiKey)
	if err != nil {
		log.Fatal("error initializing client:", err)
	}
	return &nrClient{
		client:    client,
		apiKey:    apiKey,
		accountID: accountID,
	}
}

func (nrc *nrClient) FindEntityGUID(sample string, metricName string, entityTag string) (*entities.EntityGUID, error) {

	query := fmt.Sprintf("SELECT * from %s where metricName = '%s' where tags.testKey = '%s' limit 1", sample, metricName, entityTag)

	a, err := nrc.client.Query(nrc.accountID, query)
	if err != nil {
		return nil, fmt.Errorf("executing query to fetch entity GUID %s, %w", query, err)
	}
	if len(a.Results) == 0 {
		return nil, fmt.Errorf("query to fetch entity GUID did not return any result %s", query)
	}
	firstResult := a.Results[0]
	guid := entities.EntityGUID(fmt.Sprintf("%+v", firstResult["entity.guid"]))
	return &guid, nil
}

func (nrc *nrClient) FindEntityByGUID(guid *entities.EntityGUID) (entities.EntityInterface, error) {

	if guid == nil {
		return nil, fmt.Errorf("impossible ot find entity: guid is nil")
	}

	entity, err := nrc.client.GetEntity(guid)
	if err != nil {
		return nil, fmt.Errorf("get entity: %w", err)
	}

	if entity == nil {
		return nil, fmt.Errorf("impossible ot deferentiate entity: it is nil")
	}

	return *entity, nil
}
