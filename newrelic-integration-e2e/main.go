package main

import (
	_ "embed"
	"flag"

	"github.com/newrelic/newrelic-integration-e2e/internal/executor"
	"github.com/newrelic/newrelic-integration-e2e/pkg/agent"
	"github.com/newrelic/newrelic-integration-e2e/pkg/settings"
	"github.com/sirupsen/logrus"
)

const (
	flagSpecPath    = "spec_path"
	flagVerboseMode = "verbose_mode"
	flagLicenseKey  = "license_key"
	flagAgentDir    = "agent_dir"
)

//go:embed resources/docker-compose.yml.tmpl
var dockerComposeTemplate string

func processCliArgs() (string, string, string, logrus.Level) {
	specsPathPtr := flag.String(flagSpecPath, "", "Relative path to the spec file")
	licenseKeyPtr := flag.String(flagLicenseKey, "", "New Relic License Key")
	agentDir := flag.String(flagAgentDir, "", "Directory used to deploy the agent")
	verboseModePtr := flag.Bool(flagVerboseMode, false, "If true the debug level is enabled")
	flag.Parse()
	if *licenseKeyPtr == "" {
		logrus.Fatalf("missing required license_key")
	}
	if *specsPathPtr == "" {
		logrus.Fatalf("missing required spec_path")
	}
	logLevel := logrus.InfoLevel
	if *verboseModePtr {
		logLevel = logrus.DebugLevel
	}
	return *licenseKeyPtr, *specsPathPtr, *agentDir, logLevel

}

func main() {
	logrus.Info("running executor")

	licenseKey, specPath, agentDir, logLevel := processCliArgs()
	settings, err := settings.New(
		settings.WithSpecPath(specPath),
		settings.WithLogLevel(logLevel),
		settings.WithLicenseKey(licenseKey),
		settings.WithAgentDir(agentDir),
	)
	if err != nil {
		logrus.Fatalf("error loading settings: %s", err)
	}
	logger := settings.Logger()
	logger.Debug("validating the spec definition")
	if err := settings.Spec().Validate(); err != nil {
		logger.Fatalf("error validating the spec definition: %s", err)
	}
	ag := agent.NewAgent(settings, dockerComposeTemplate)
	if err := executor.Exec(ag, settings); err != nil {
		logger.Fatal(err)
	}
	logger.Info("execution completed successfully!")
}
