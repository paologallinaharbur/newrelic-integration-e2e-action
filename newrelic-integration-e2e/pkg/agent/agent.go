package agent

import (
	_ "embed"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/newrelic/newrelic-integration-e2e/internal/docker"
	"github.com/newrelic/newrelic-integration-e2e/pkg/settings"
	"github.com/newrelic/newrelic-integration-e2e/pkg/spec"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	integrationsCfgDir  = "integrations.d"
	integrationsBinDir  = "bin"
	dockerCompose       = "docker-compose.yml"
	defConfigFile       = "nri-config.yml"
	varIntegrationNames = "integrations"
	varLicenseKey       = "licenseKey"
)

type Agent interface {
	SetUp(scenario spec.Scenario) error
	Launch() error
}

type agent struct {
	agentDir              string
	pathDockerCompose     string
	configsDir            string
	binsDir               string
	dockerComposeTemplate string
	licenseKey            string
	defConfigFile         string
	rootDir               string
	logger                *logrus.Logger
}

func NewAgent(settings settings.Settings, dockerComposeTemplate string) *agent {
	agentDir := settings.AgentDir()
	return &agent{
		rootDir:               settings.RootDir(),
		agentDir:              agentDir,
		pathDockerCompose:     filepath.Join(agentDir, dockerCompose),
		configsDir:            filepath.Join(agentDir, integrationsCfgDir),
		binsDir:               filepath.Join(agentDir, integrationsBinDir),
		defConfigFile:         filepath.Join(agentDir, integrationsCfgDir, defConfigFile),
		dockerComposeTemplate: dockerComposeTemplate,
		licenseKey:            settings.LicenseKey(),
		logger:                settings.Logger(),
	}
}

func (a *agent) initialize() error {
	a.logger.Debug("removing the content of the root folder")
	if err := removeDirectoryContent(a.agentDir); err != nil {
		return err
	}
	a.logger.Debug("creating the required folders by the agent")
	return makeDirs(0777, a.configsDir, a.binsDir)
}

func (a *agent) addIntegration(integration spec.Integration) error {
	if integration.Path == "" {
		return nil
	}
	if _, err := copyFile(filepath.Join(a.rootDir, integration.Path), filepath.Join(a.agentDir, integrationsBinDir, filepath.Base(integration.Path))); err != nil {
		return err
	}
	return nil
}

func (a *agent) addIntegrationsConfigFile(integrations spec.Integrations) error {
	content, err := yaml.Marshal(integrations)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(a.agentDir, integrationsCfgDir, "nri-config.yml"), content, 0777)
}

func (a *agent) addDockerCompose(licenseKey string, integrations []string) error {
	t, err := template.New("").Parse(a.dockerComposeTemplate)
	if err != nil {
		return err
	}
	vars := map[string]interface{}{
		varIntegrationNames: integrations,
		varLicenseKey:       licenseKey,
	}
	return processTemplate(t, vars, a.pathDockerCompose)
}

func (a *agent) SetUp(scenario spec.Scenario) error {
	if err := a.initialize(); err != nil {
		return err
	}
	integrations := scenario.Integrations
	a.logger.Debugf("there are %d integrations", len(integrations))
	integrationsNames := make([]string, len(integrations))
	for i := range integrations {
		if err := a.addIntegration(integrations[i]); err != nil {
			return err
		}
		integrationsNames[i] = integrations[i].Name
	}
	if err := a.addIntegrationsConfigFile(integrations); err != nil {
		return err
	}
	return a.addDockerCompose(a.licenseKey, integrationsNames)
}

func (a *agent) Launch() error {
	return docker.DockerComposeUp(filepath.Join(a.agentDir, dockerCompose))
}
