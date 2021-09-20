package newrelic

import (
	"errors"
	"fmt"
	"log"

	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

type Client interface {
	FindEntityGUID(sample, metricName, customTagKey, entityTag string) (*entities.EntityGUID, error)
	FindEntityByGUID(guid *entities.EntityGUID) (entities.EntityInterface, error)
	FindEntityMetrics()
}

var (
	ErrNilEntity = errors.New("nil entity, impossible to dereference")
	ErrNilGUID   = errors.New("GUID is nil, impossible to find entity")
	ErrNoResult  = errors.New("query to fetch entity GUID did not return any result")
)

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

func (nrc *nrClient) FindEntityGUID(sample, metricName, customTagKey, entityTag string) (*entities.EntityGUID, error) {
	query := fmt.Sprintf("SELECT * from %s where metricName = '%s' where %s = '%s' limit 1", sample, metricName, customTagKey, entityTag)

	a, err := nrc.client.Query(nrc.accountID, query)
	if err != nil {
		return nil, fmt.Errorf("executing query to fetch entity GUID %s, %w", query, err)
	}
	if len(a.Results) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrNoResult, query)
	}
	firstResult := a.Results[0]
	if firstResult["entity.guid"] == nil {
		return nil, ErrNilGUID
	}

	guid := entities.EntityGUID(fmt.Sprintf("%+v", firstResult["entity.guid"]))
	return &guid, nil
}

func (nrc *nrClient) FindEntityByGUID(guid *entities.EntityGUID) (entities.EntityInterface, error) {
	if guid == nil {
		return nil, ErrNilGUID
	}

	entity, err := nrc.client.GetEntity(guid)
	if err != nil {
		return nil, fmt.Errorf("get entity: %w", err)
	}

	if entity == nil {
		return nil, ErrNilEntity
	}

	return *entity, nil
}

func (nrc *nrClient) FindEntityMetrics(sample, metricName, customTagKey, entityTag string) ([]string, error) {

	query := fmt.Sprintf("SELECT keyset() from %s where metricName = '%s' where %s = '%s' limit 1", sample, metricName, customTagKey, entityTag)

	a, err := nrc.client.Query(nrc.accountID, query)
	if err != nil {
		return nil, fmt.Errorf("executing query to fetch entity GUID %s, %w", query, err)
	}
	if len(a.Results) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrNoResult, query)
	}
	firstResult := a.Results[0]
	guid := entities.EntityGUID(fmt.Sprintf("%+v", firstResult["entity.guid"]))
	return &guid, nil
}
