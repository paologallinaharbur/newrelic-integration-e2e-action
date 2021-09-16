package main

import (
	_ "embed"
	"flag"

	e2e "github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/agent"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/executor"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/newrelic"
	"github.com/sirupsen/logrus"
)

const (
	flagSpecPath    = "spec_path"
	flagVerboseMode = "verbose_mode"
	flagApiKey      = "api_key"
	flagAccountID   = "account_id"
	flagLicenseKey  = "license_key"
	flagAgentDir    = "agent_dir"
	flagRootDir     = "root_dir"
)

func processCliArgs() (string, string, string, string, string, int, logrus.Level) {
	specsPath := flag.String(flagSpecPath, "", "Relative path to the spec file")
	licenseKey := flag.String(flagLicenseKey, "", "New Relic License Key")
	agentDir := flag.String(flagAgentDir, "", "Directory used to deploy the agent")
	rootDir := flag.String(flagRootDir, "", "workspace directory")
	verboseMode := flag.Bool(flagVerboseMode, false, "If true the debug level is enabled")
	apiKey := flag.String(flagApiKey, "", "New Relic Api Key")
	accountID := flag.Int(flagAccountID, 0, "New Relic accountID to be used")
	flag.Parse()

	if *licenseKey == "" {
		logrus.Fatalf("missing required license_key")
	}
	if *specsPath == "" {
		logrus.Fatalf("missing required spec_path")
	}
	if *rootDir == "" {
		logrus.Fatalf("missing required root_dir")
	}
	if *accountID == 0 {
		logrus.Fatalf("missing required accountID")
	}
	if *apiKey == "" {
		logrus.Fatalf("missing required apiKey")
	}

	logLevel := logrus.InfoLevel
	if *verboseMode {
		logLevel = logrus.DebugLevel
	}
	return *licenseKey, *specsPath, *rootDir, *agentDir, *apiKey, *accountID, logLevel

}

func main() {
	logrus.Info("running executor")

	licenseKey, specsPath, rootDir, agentDir, apiKey, accountID, logLevel := processCliArgs()
	s, err := e2e.NewSettings(
		e2e.SettingsWithSpecPath(specsPath),
		e2e.SettingsWithLogLevel(logLevel),
		e2e.SettingsWithLicenseKey(licenseKey),
		e2e.SettingsWithAgentDir(agentDir),
		e2e.SettingsWithRootDir(rootDir),
		e2e.SettingsWithApiKey(apiKey),
		e2e.SettingsWithAccountID(accountID),
	)
	if err != nil {
		logrus.Fatalf("error loading s: %s", err)
	}
	logger := s.Logger()

	logger.Debug("validating the spec definition")
	if err := s.SpecDefinition().Validate(); err != nil {
		logger.Fatalf("error validating the spec definition: %s", err)
	}

	e2eExecutor := executor.NewExecutor(
		agent.NewAgent(s),
		newrelic.NewNrClient(s.ApiKey(), s.AccountID()),
		s,
	)

	if err := e2eExecutor.Exec(); err != nil {
		logger.Fatal(err)
	}

	logger.Info("execution completed successfully!")
}
