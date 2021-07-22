package executor

import (
	"fmt"
	"os/exec"
	"time"

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
		if err := executeOSCommands(scenario.Before); err != nil {
			return err
		}
		if err:=ag.Launch();err!=nil{
			return err
		}

		/**
		This block will be used to run the tests
		 */




		time.Sleep(1 * time.Minute)



		if err := executeOSCommands(scenario.After); err != nil {
			println(err.Error())
		}
		if err:=ag.Stop();err!=nil{
			return err
		}
	}

	return nil
}

func executeOSCommands(statements []string) error {
	for i := range statements {
		stmt := statements[i]
		fmt.Println(stmt)
		cmd := exec.Command("bash", "-c", stmt)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}
	return nil
}

