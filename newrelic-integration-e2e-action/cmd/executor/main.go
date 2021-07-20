package main

import (
	"log"

	"github.com/newrelic/newrelic-integration-e2e-action/cmd/common"
	"github.com/newrelic/newrelic-integration-e2e-action/pkg/settings"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("running executor")
	cfg := common.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
	settings, err := settings.New(
		settings.WithSpecPath(cfg.SpecPath()),
		settings.WithLogLevel(cfg.LogLevel()),
	)
	logger := settings.Logger()
	if err != nil {
		logger.Fatal(err)
	}
	if err := settings.SpecDefinition().Validate(); err != nil {
		logger.Fatal(err)
	}
}
