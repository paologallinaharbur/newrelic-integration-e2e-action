package main

import (
	"log"

	"github.com/newrelic/newrelic-integration-e2e-action/cmd/common"
	"github.com/newrelic/newrelic-integration-e2e-action/pkg/executor"
	"github.com/newrelic/newrelic-integration-e2e-action/pkg/settings"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("running executor")
	cfg := common.LoadConfig()
	if err := cfg.Validate(); err != nil {
		logrus.Fatalf("error validating the flags: %s",err)
		log.Fatal(err)
	}
	settings, err := settings.New(
		settings.WithSpecPath(cfg.SpecPath()),
		settings.WithLogLevel(cfg.LogLevel()),
	)
	if err != nil {
		logrus.Fatalf("error loading settings: %s",err)
	}
	logger := settings.Logger()
	specDefinition:=settings.SpecDefinition()
	logger.Debug("validating the spec definition")
	if err := specDefinition.Validate(); err != nil {
		logger.Fatalf("error validating the spec definition: %s",err)
	}
	logger.Debug("executing the spec")
	if err:=executor.Execute(settings);err!=nil{
		logger.Fatalf("error running the spec: %s",err)
	}

}
