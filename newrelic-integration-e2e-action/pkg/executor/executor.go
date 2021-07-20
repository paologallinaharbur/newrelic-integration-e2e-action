package executor

import (
	"github.com/newrelic/newrelic-integration-e2e-action/pkg/settings"
	"github.com/newrelic/newrelic-integration-e2e-action/pkg/spec"
)

func setUpIntegrations(settings settings.Settings, integration []spec.Integration) error{
	ag:= agent{rootDir: settings.RootDir()}
	if err:=ag.initialize(settings.Logger());err!=nil{
		return err
	}
	return nil
}

func Execute(settings settings.Settings) error {
	definition := settings.SpecDefinition()
	for i := range definition.Scenarios {
		scenario := definition.Scenarios[i]
		setUpIntegrations(settings, scenario.Integrations)
	}
	return nil
}
