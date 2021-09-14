package agent

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/dockercompose"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/settings"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/spec"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	integrationsCfgDir = "integrations.d"
	exportersDir       = "exporters"
	integrationsBinDir = "bin"
	dockerCompose      = "docker-compose.yml"
	defConfigFile      = "nri-config.yml"
	containerName      = "agent"
	infraAgentDir      = "newrelic-infra-agent"
	customTagKey       = "testKey"
)

type Agent interface {
	SetUp(logger *logrus.Logger, scenario spec.Scenario) error
	Launch(scenarioTag string) error
	Stop() error
}

type agent struct {
	scenario          spec.Scenario
	agentDir          string
	configsDir        string
	exportersDir      string
	binsDir           string
	licenseKey        string
	defConfigFile     string
	rootDir           string
	dockerComposePath string
	logger            *logrus.Logger
	overrides         *spec.Agent
}

func NewAgent(settings settings.Settings) *agent {
	agentDir := settings.AgentDir()
	return &agent{
		rootDir:           settings.RootDir(),
		agentDir:          agentDir,
		configsDir:        filepath.Join(agentDir, infraAgentDir, integrationsCfgDir),
		exportersDir:      filepath.Join(agentDir, infraAgentDir, exportersDir),
		binsDir:           filepath.Join(agentDir, infraAgentDir, integrationsBinDir),
		defConfigFile:     filepath.Join(agentDir, infraAgentDir, integrationsCfgDir, defConfigFile),
		dockerComposePath: filepath.Join(agentDir, dockerCompose),
		licenseKey:        settings.LicenseKey(),
		logger:            settings.Logger(),
		overrides:         settings.Spec().AgentOverrides,
	}
}

func (a *agent) initialize() error {
	a.logger.Debug("removing temporary folders")
	if err := removeDirectories(a.exportersDir, a.configsDir, a.binsDir); err != nil {
		return err
	}
	a.logger.Debug("creating folders required by the agent")
	return makeDirs(0777, a.exportersDir, a.configsDir, a.binsDir)
}

func (a *agent) addIntegration(logger *logrus.Logger, integration spec.Integration) error {
	if integration.BinaryPath == "" {
		return nil
	}
	source := filepath.Join(a.rootDir, integration.BinaryPath)
	destination := filepath.Join(a.binsDir, integration.Name)
	logger.Debugf("copy file from '%s' to '%s'", source, destination)
	return copyFile(source, destination)
}

func (a *agent) addPrometheusExporter(logger *logrus.Logger, integration spec.Integration) error {
	if integration.ExporterBinaryPath == "" {
		return nil
	}
	exporterName := filepath.Base(integration.ExporterBinaryPath)
	source := filepath.Join(a.rootDir, integration.ExporterBinaryPath)
	destination := filepath.Join(a.exportersDir, exporterName)
	logger.Debugf("copy file from '%s' to '%s'", source, destination)
	return copyFile(source, destination)
}

func (a *agent) addIntegrationsConfigFile(logger *logrus.Logger, integrations []spec.Integration) error {
	content, err := yaml.Marshal(createAgentIntegrationModel(integrations))
	if err != nil {
		return err
	}
	cfgPath := filepath.Join(a.configsDir, defConfigFile)
	logger.Debugf("create config file '%s' in  '%s'", defConfigFile, cfgPath)
	return ioutil.WriteFile(cfgPath, content, 0777)
}

func (a *agent) SetUp(logger *logrus.Logger, scenario spec.Scenario) error {
	a.scenario = scenario
	if err := a.initialize(); err != nil {
		return err
	}
	integrations := scenario.Integrations
	a.logger.Debugf("there are %d integrations", len(integrations))
	integrationsNames := make([]string, len(integrations))
	for i := range integrations {
		integration := integrations[i]
		if err := a.addIntegration(logger, integration); err != nil {
			return err
		}
		if err := a.addPrometheusExporter(logger, integration); err != nil {
			return err
		}
		integrationsNames[i] = integration.Name
	}
	if err := a.addIntegrationsConfigFile(logger, integrations); err != nil {
		return err
	}
	for k, v := range a.overrides.Integrations {
		source := filepath.Join(a.rootDir, v)
		destination := filepath.Join(a.binsDir, k)
		return copyFile(source, destination)
	}
	return nil
}

func (a *agent) Launch(scenarioTag string) error {

	return dockercompose.Run(a.dockerComposePath, containerName, map[string]string{
		"NRIA_VERBOSE":           "1",
		"NRIA_LICENSE_KEY":       a.licenseKey,
		"NRIA_CUSTOM_ATTRIBUTES": fmt.Sprintf(`'{"%s":"%s"}'`, customTagKey, scenarioTag),
	})
}

func (a *agent) Stop() error {

	return dockercompose.Down(a.dockerComposePath)
}
