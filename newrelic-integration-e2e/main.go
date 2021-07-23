package main

import (
	_ "embed"
	"flag"

	"github.com/newrelic/newrelic-integration-e2e/pkg/agent"
	"github.com/newrelic/newrelic-integration-e2e/pkg/settings"
	"github.com/newrelic/newrelic-integration-e2e/pkg/spec"
	"github.com/sirupsen/logrus"
)

const (
	flagSpecPath    = "spec_path"
	flagVerboseMode = "verbose_mode"
	flagLicenseKey  = "license_key"
	flagAgentDir    = "agent_dir"
	flagRootDir     = "root_dir"
)

func processCliArgs() (string, string, string, string, logrus.Level) {
	specsPathPtr := flag.String(flagSpecPath, "", "Relative path to the spec file")
	licenseKeyPtr := flag.String(flagLicenseKey, "", "New Relic License Key")
	agentDir := flag.String(flagAgentDir, "", "Directory used to deploy the agent")
	flagRootDirPtr := flag.String(flagRootDir, "", "workspace directory")
	verboseModePtr := flag.Bool(flagVerboseMode, false, "If true the debug level is enabled")
	flag.Parse()
	if *licenseKeyPtr == "" {
		logrus.Fatalf("missing required license_key")
	}
	if *specsPathPtr == "" {
		logrus.Fatalf("missing required spec_path")
	}
	if *flagRootDirPtr == "" {
		logrus.Fatalf("missing required root_dir")
	}
	logLevel := logrus.InfoLevel
	if *verboseModePtr {
		logLevel = logrus.DebugLevel
	}
	return *licenseKeyPtr, *specsPathPtr, *flagRootDirPtr, *agentDir, logLevel

}

func main() {
	logrus.Info("running executor")

	licenseKey, specPath, rootDir, agentDir, logLevel := processCliArgs()
	s, err := settings.New(
		settings.WithSpecPath(specPath),
		settings.WithLogLevel(logLevel),
		settings.WithLicenseKey(licenseKey),
		settings.WithAgentDir(agentDir),
		settings.WithRootDir(rootDir),
	)
	if err != nil {
		logrus.Fatalf("error loading s: %s", err)
	}
	logger := s.Logger()
	logger.Debug("validating the spec definition")
	if err := s.Spec().Validate(); err != nil {
		logger.Fatalf("error validating the spec definition: %s", err)
	}
	ag := agent.NewAgent(s)
	if err := spec.Exec(ag, s); err != nil {
		logger.Fatal(err)
	}
	logger.Info("execution completed successfully!")
}
