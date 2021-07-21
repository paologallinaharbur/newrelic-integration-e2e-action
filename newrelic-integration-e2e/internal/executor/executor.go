package executor

import (
	"github.com/newrelic/newrelic-integration-e2e/pkg/agent"
	"github.com/newrelic/newrelic-integration-e2e/pkg/settings"
)

func Exec(ag agent.Agent, settings settings.Settings) error{
	spec := settings.Spec()
	for i := range spec.Scenarios {
		scenario := spec.Scenarios[i]
		settings.Logger().Debugf("[scenario]: %s", scenario.Description)
		if err:=ag.SetUp(scenario);err!=nil{
			return err
		}
		if err:=ag.Launch();err!=nil{
			return err
		}
	}

	return nil
}
