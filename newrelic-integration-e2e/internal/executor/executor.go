package executor

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/newrelic"

	"github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/agent"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/settings"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/spec"
)

func Exec(ag agent.Agent, nrc newrelic.DataClient, settings settings.Settings) error {
	testSpec := settings.Spec()
	logger := settings.Logger()
	for _, scenario := range testSpec.Scenarios {
		// TODO Improve tag with more info from each scenario
		scenarioTag := RandStringRunes(10)
		logger.Debugf("[scenario]: %s, [Tag]: %s", scenario.Description, scenarioTag)

		if err := ag.SetUp(settings.Logger(), scenario); err != nil {
			return err
		}

		if err := executeOSCommands(settings, scenario.Before); err != nil {
			return err
		}

		if err := ag.Launch(scenarioTag); err != nil {
			return err
		}

		if err := executeTests(settings, nrc, scenario.Tests, scenarioTag); err != nil {
			return err
		}

		if err := executeOSCommands(settings, scenario.After); err != nil {
			logger.Error(err)
		}

		if err := ag.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func executeOSCommands(settings settings.Settings, statements []string) error {
	logger := settings.Logger()
	rootDir := settings.RootDir()
	for _, stmt := range statements {
		logger.Debugf("execute command '%s' from path '%s'", stmt, rootDir)
		cmd := exec.Command("bash", "-c", stmt)
		cmd.Dir = rootDir
		stdout, err := cmd.Output()
		logrus.Debug(stdout)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO Interface to specify it? needed?

func executeTests(settings settings.Settings, nrc newrelic.DataClient, tests spec.Tests, scenarioTag string) error {

	err := retry(settings.Logger(), 10, 1*time.Minute, func() []error {

		testEntities(tests.Entities, nrc, scenarioTag)
		//testNRQLs(tests.NRQLs)
		//testMetrics(tests.Metrics)
		return nil
	})

	return err
}

func testEntities(entities []spec.Entities, nrc newrelic.DataClient, tag string) []error {
	var errors []error
	for _, e := range entities {
		guid, err := nrc.FindEntityGUID(e.Type, e.MetricName, tag)
		if err != nil {
			errors = append(errors, fmt.Errorf("finding entity guid: %w", err))
			continue
		}
		entity, err := nrc.FindEntityByGUID(guid)
		if err != nil {
			errors = append(errors, fmt.Errorf("finding entity guid: %w", err))
			continue
		}

		if entity.GetType() != e.Type {
			errors = append(errors, fmt.Errorf("enttiy type is not matching: %s!=%s", entity.GetType(), e.Type))
			continue
		}
	}
	return errors
}
