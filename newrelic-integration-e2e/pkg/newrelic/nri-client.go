package newrelic

import (
	"fmt"
	"log"

	newrelicgo "github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

type DataClient interface {
	FindEntityGUID(sample string, metricName string, entityTag string) (*entities.EntityGUID, error)
	FindEntityByGUID(guid *entities.EntityGUID) (entities.EntityInterface, error)
}

type nrClient struct {
	accountID int
	apiKey    string
	client    *newrelicgo.NewRelic
}

func (nrc *nrClient) FindEntityGUID(sample string, metricName string, entityTag string) (*entities.EntityGUID, error) {

	query := fmt.Sprintf("SELECT * from %s where metricName = '%s' where tags.testKey = '%s' limit 1", sample, metricName, entityTag)

	a, err := nrc.client.Nrdb.Query(nrc.accountID, nrdb.NRQL(query))
	if err != nil {
		return nil, fmt.Errorf("executing query to fetch entity GUID %s, %w", query, err.Error())
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

	entity, err := nrc.client.Entities.GetEntity(*guid)
	if err != nil {
		return nil, fmt.Errorf("get entity: %w", err)
	}

	if entity == nil {
		return nil, fmt.Errorf("impossible ot deferentiate entity: it is nil")
	}

	return *entity, nil
}

func NewNrClient(apiKey string, accountID int) *nrClient {

	client, err := newrelicgo.New(newrelicgo.ConfigPersonalAPIKey(apiKey))
	if err != nil {
		log.Fatal("error initializing client:", err)
	}
	return &nrClient{
		client:    client,
		apiKey:    apiKey,
		accountID: accountID,
	}
}
