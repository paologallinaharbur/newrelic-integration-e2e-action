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

type Executor struct {
	agent    agent.Agent
	nrClient newrelic.DataClient
	logger   *logrus.Logger
	spec     *spec.Definition
	rootDir  string
}

func NewExecutor(agent agent.Agent, nrClient newrelic.DataClient, settings settings.Settings) *Executor {
	return &Executor{
		agent:    agent,
		nrClient: nrClient,
		logger:   settings.Logger(),
		spec:     settings.Spec(),
		rootDir:  settings.RootDir(),
	}
}

func (ex *Executor) Exec() error {
	for _, scenario := range ex.spec.Scenarios {
		// TODO Improve tag with more info from each scenario, like GH commit
		if err := ex.agent.SetUp(scenario); err != nil {
			return err
		}
		ex.logger.Debugf("[scenario]: %s, [Tag]: %s", scenario.Description)

		if err := ex.executeOSCommands(scenario.Before); err != nil {
			return err
		}

		if err := ex.agent.Launch(); err != nil {
			return err
		}

		errAssertions := ex.executeTests(scenario.Tests)

		if err := ex.executeOSCommands(scenario.After); err != nil {
			ex.logger.Error(err)
		}

		if err := ex.agent.Stop(); err != nil {
			return err
		}

		if errAssertions != nil {
			return errAssertions
		}
	}

	return nil
}

func (ex *Executor) executeOSCommands(statements []string) error {
	for _, stmt := range statements {
		ex.logger.Debugf("execute command '%s' from path '%s'", stmt, ex.rootDir)
		cmd := exec.Command("bash", "-c", stmt)
		cmd.Dir = ex.rootDir
		stdout, err := cmd.Output()
		logrus.Debug(stdout)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO Interface to specify it? needed?

func (ex *Executor) executeTests(tests spec.Tests) error {
	return retry(ex.logger, 10, 60*time.Second, func() []error {
		errors := ex.testEntities(tests.Entities)
		errors = append(
			errors,
			ex.testNRQLs(tests.NRQLs)...,
		)
		errors = append(
			errors,
			ex.testMetrics(tests.Metrics)...,
		)
		return errors
	})
}

func (ex *Executor) testEntities(entities []spec.Entity) []error {
	var errors []error
	for _, en := range entities {
		guid, err := ex.nrClient.FindEntityGUID(en.DataType, en.MetricName, ex.agent.GetCustomTagKey(), ex.agent.GetCustomTagValue())
		if err != nil {
			errors = append(errors, fmt.Errorf("finding entity guid: %w", err))
			continue
		}
		entity, err := ex.nrClient.FindEntityByGUID(guid)
		if err != nil {
			errors = append(errors, fmt.Errorf("finding entity guid: %w", err))
			continue
		}

		if entity.GetType() != en.Type {
			errors = append(errors, fmt.Errorf("enttiy type is not matching: %s!=%s", entity.GetType(), en.Type))
			continue
		}
	}
	return errors
}

func (ex *Executor) testNRQLs(nrqls []spec.NRQL) []error {
	var errors []error
	return errors
}

func (ex *Executor) testMetrics(metrics []spec.Metrics) []error {
	var errors []error
	return errors
}
