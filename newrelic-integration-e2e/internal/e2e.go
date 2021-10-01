package e2e

import (
	"fmt"

	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/agent"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/newrelic"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/runtime"
)

func Exec(settings Settings) error {
	settings.Logger().Debug("validating the spec definition")

	if err := settings.SpecDefinition().Validate(); err != nil {
		return fmt.Errorf("error validating the spec definition: %s", err)
	}

	nrClient := newrelic.NewNrClient(settings.ApiKey(), settings.AccountID())
	entitiesTester := runtime.NewEntitiesTester(nrClient, settings.Logger())
	metricsTester := runtime.NewMetricsTester(nrClient, settings.Logger(), settings.SpecParentDir())
	nrqlTester := runtime.NewNRQLTester(nrClient, settings.Logger())

	runner := runtime.NewRunner(
		agent.NewAgent(settings),
		[]runtime.Tester{
			entitiesTester,
			metricsTester,
			nrqlTester,
		},
		settings,
	)

	return runner.Run()
}
