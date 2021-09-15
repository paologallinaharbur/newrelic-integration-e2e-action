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
		// TODO Improve tag with more info from each scenario, like GH commit
		if err := ag.SetUp(settings.Logger(), scenario); err != nil {
			return err
		}
		logger.Debugf("[scenario]: %s, [Tag]: %s", scenario.Description)

		if err := executeOSCommands(settings, scenario.Before); err != nil {
			return err
		}

		if err := ag.Launch(); err != nil {
			return err
		}

		errAssertions := executeTests(ag, settings, nrc, scenario.Tests)

		if err := executeOSCommands(settings, scenario.After); err != nil {
			logger.Error(err)
		}

		if err := ag.Stop(); err != nil {
			return err
		}

		if errAssertions != nil {
			return errAssertions
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

func executeTests(ag agent.Agent, settings settings.Settings, nrc newrelic.DataClient, tests spec.Tests) error {
	return retry(settings.Logger(), 10, 60*time.Second, func() []error {
		errors := testEntities(tests.Entities, nrc, ag)
		errors = append(
			errors,
			testNRQLs(tests.NRQLs, nrc, ag)...,
		)
		errors = append(
			errors,
			testMetrics(tests.Metrics, nrc, ag)...,
		)
		return errors
	})
}

func testEntities(entities []spec.Entity, nrc newrelic.DataClient, ag agent.Agent) []error {
	var errors []error
	for _, e := range entities {
		guid, err := nrc.FindEntityGUID(e.DataType, e.MetricName, ag.GetCustomTagKey(), ag.GetCustomTagValue())
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

func testNRQLs(nrqls []spec.NRQL, nrc newrelic.DataClient, ag agent.Agent) []error {
	var errors []error
	return errors
}

func testMetrics(metrics []spec.Metrics, nrc newrelic.DataClient, ag agent.Agent) []error {
	var errors []error
	return errors
}
