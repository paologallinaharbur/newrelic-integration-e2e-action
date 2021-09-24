package executor

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"

	e2e "github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/agent"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/newrelic"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/spec"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/retrier"
	"github.com/sirupsen/logrus"
)

const dmTableName = "Metric"

type Executor struct {
	agent         agent.Agent
	nrClient      newrelic.Client
	logger        *logrus.Logger
	spec          *spec.Definition
	specParentDir string
}

func NewExecutor(agent agent.Agent, nrClient newrelic.Client, settings e2e.Settings) *Executor {
	return &Executor{
		agent:         agent,
		nrClient:      nrClient,
		logger:        settings.Logger(),
		spec:          settings.SpecDefinition(),
		specParentDir: settings.SpecParentDir(),
	}
}

func (ex *Executor) Exec() error {
	for _, scenario := range ex.spec.Scenarios {
		// TODO Improve tag with more info from each scenario, like GH commit
		if err := ex.agent.SetUp(scenario); err != nil {
			return err
		}
		ex.logger.Debugf("[scenario]: %s, [Tag]: %s", scenario.Description, ex.agent.GetCustomTagValue())

		if err := ex.executeOSCommands(scenario.Before); err != nil {
			return err
		}

		if err := ex.agent.Run(); err != nil {
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
		ex.logger.Debugf("execute command '%s' from path '%s'", stmt, ex.specParentDir)
		cmd := exec.Command("bash", "-c", stmt)
		cmd.Dir = ex.specParentDir
		stdout, err := cmd.Output()
		logrus.Debug(stdout)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ex *Executor) executeTests(tests spec.Tests) error {
	return retrier.Retry(ex.logger, 10, 60*time.Second, func() []error {
		errors := ex.testEntities(tests.Entities)
		if len(errors) == 0 {
			errors = append(
				errors,
				ex.testNRQLs(tests.NRQLs)...,
			)
			errors = append(
				errors,
				ex.testMetrics(tests.Metrics)...,
			)
		}
		return errors
	})
}

func (ex *Executor) testEntities(entities []spec.TestEntity) []error {
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
			errors = append(errors, fmt.Errorf("entity type is not matching: %s!=%s", entity.GetType(), en.Type))
			continue
		}
	}
	return errors
}

func (ex *Executor) testNRQLs(nrqls []spec.TestNRQL) []error {
	var errors []error
	for _, nrql := range nrqls {
		err := ex.nrClient.NRQLQuery(nrql.Query, ex.agent.GetCustomTagKey(), ex.agent.GetCustomTagValue())
		if err != nil {
			errors = append(errors, fmt.Errorf("querying: %w", err))
			continue
		}
	}
	return errors
}

func (ex *Executor) testMetrics(testMetrics []spec.TestMetrics) []error {
	var errors []error
	for _, tm := range testMetrics {
		content, err := ioutil.ReadFile(filepath.Join(ex.specParentDir, tm.Source))
		if err != nil {
			errors = append(errors, fmt.Errorf("reading metrics source file: %w", err))
			continue
		}
		ex.logger.Debug("parsing the content of the metrics source file")
		metrics, err := spec.ParseMetricsFile(content)
		if err != nil {
			errors = append(errors, fmt.Errorf("unmarshaling metrics source file: %w", err))
			continue
		}

		queriedMetrics, err := ex.nrClient.FindEntityMetrics(dmTableName, ex.agent.GetCustomTagKey(), ex.agent.GetCustomTagValue())
		if err != nil {
			errors = append(errors, fmt.Errorf("finding keyset: %w", err))
			continue
		}

		for _, entity := range metrics.Entities {
			if isEntityException(entity.EntityType, tm.ExceptEntities) {
				continue
			}

			for _, metric := range entity.Metrics {
				if isMetricException(metric.Name, tm.ExceptMetrics) {
					continue
				}

				if containsMetric(metric.Name, queriedMetrics) {
					continue
				} else {
					errors = append(errors, fmt.Errorf("finding Metric: %v", metric.Name))
					continue
				}
			}
		}
	}
	return errors
}

func isEntityException(entity string, entitiesList []string) bool {
	for _, entityType := range entitiesList {
		if entityType == entity {
			return true
		}
	}
	return false
}

func isMetricException(metric string, exceptionMetricsList []string) bool {
	for _, exceptMetric := range exceptionMetricsList {
		if exceptMetric == metric {
			return true
		}
	}
	return false
}

func containsMetric(metric string, queriedMetricsList []string) bool {
	for _, queriedMetric := range queriedMetricsList {
		if queriedMetric == metric {
			return true
		}
	}
	return false
}
