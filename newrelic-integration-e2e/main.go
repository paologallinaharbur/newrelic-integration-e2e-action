package main

import (
	_ "embed"
	"flag"

	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/newrelic"

	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal/executor"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/agent"
	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/settings"
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
	specsPath := *flag.String(flagSpecPath, "", "Relative path to the spec file")
	licenseKey := *flag.String(flagLicenseKey, "", "New Relic License Key")
	agentDir := *flag.String(flagAgentDir, "", "Directory used to deploy the agent")
	rootDir := *flag.String(flagRootDir, "", "workspace directory")
	verboseMode := *flag.Bool(flagVerboseMode, false, "If true the debug level is enabled")
	apiKey := *flag.String(flagApiKey, "", "New Relic Api Key")
	accountID := *flag.Int(flagAccountID, 0, "New Relic accountID to be used")

	flag.Parse()
	if licenseKey == "" {
		logrus.Fatalf("missing required license_key")
	}
	if specsPath == "" {
		logrus.Fatalf("missing required spec_path")
	}
	if rootDir == "" {
		logrus.Fatalf("missing required root_dir")
	}
	if accountID == 0 {
		logrus.Fatalf("missing required accountID")
	}
	if apiKey == "" {
		logrus.Fatalf("missing required apiKey")
	}

	logLevel := logrus.InfoLevel
	if verboseMode {
		logLevel = logrus.DebugLevel
	}
	return licenseKey, specsPath, rootDir, agentDir, apiKey, accountID, logLevel

}

func main() {
	logrus.Info("running executor")

	licenseKey, specsPath, rootDir, agentDir, apiKey, accountID, logLevel := processCliArgs()
	s, err := settings.New(
		settings.WithSpecPath(specsPath),
		settings.WithLogLevel(logLevel),
		settings.WithLicenseKey(licenseKey),
		settings.WithAgentDir(agentDir),
		settings.WithRootDir(rootDir),
		settings.WithApiKey(apiKey),
		settings.WithAccountID(accountID),
	)
	if err != nil {
		logrus.Fatalf("error loading s: %s", err)
	}
	logger := s.Logger()

	logger.Debug("validating the spec definition")
	if err := s.Spec().Validate(); err != nil {
		logger.Fatalf("error validating the spec definition: %s", err)
	}

	nrc := newrelic.NewNrClient(s.ApiKey(), s.AccountID())
	ag := agent.NewAgent(s)

	if err := executor.Exec(ag, nrc, s); err != nil {
		logger.Fatal(err)
	}

	logger.Info("execution completed successfully!")
}
