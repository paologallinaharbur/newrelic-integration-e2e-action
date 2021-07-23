package executor

import (
	"os/exec"
	"time"

	"github.com/newrelic/newrelic-integration-e2e/pkg/agent"
	"github.com/newrelic/newrelic-integration-e2e/pkg/settings"
	"github.com/sirupsen/logrus"
)

func Exec(ag agent.Agent, settings settings.Settings) error {
	spec := settings.Spec()
	for i := range spec.Scenarios {
		scenario := spec.Scenarios[i]
		settings.Logger().Debugf("[scenario]: %s", scenario.Description)
		if err := ag.SetUp(settings.Logger(), scenario); err != nil {
			return err
		}
		if err := executeOSCommands(settings, scenario.Before); err != nil {
			return err
		}
		if err := ag.Launch(); err != nil {
			return err
		}

		/**
		This block will be used to run the tests
		*/

		time.Sleep(1 * time.Minute)

		if err := executeOSCommands(settings, scenario.After); err != nil {
			println(err.Error())
		}
		if err := ag.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func executeOSCommands(settings settings.Settings, statements []string) error {
	logger := settings.Logger()
	rootDir := settings.RootDir()
	for i := range statements {
		stmt := statements[i]
		logger.Debugf("execute command '%s' from path '%s'", stmt, rootDir)
		cmd := exec.Command("bash", "-c", stmt)
		cmd.Dir = rootDir
		stdout, err := cmd.Output()
		if err != nil {
			logrus.Error(stdout)
			return err
		}
	}
	return nil
}
