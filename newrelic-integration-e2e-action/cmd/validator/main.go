package main

import (
	"log"

	"github.com/newrelic/newrelic-integration-e2e-action/cmd/common"
	"github.com/newrelic/newrelic-integration-e2e-action/pkg/settings"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("running validator")
	cfg := common.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
	settings, err := settings.New(
		settings.WithSpecPath(cfg.SpecPath()),
		settings.WithLogLevel(cfg.LogLevel()),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	logger := settings.Logger()

	if err := settings.SpecDefinition().Validate(); err != nil {
		logger.Fatal(err)
	}
}
