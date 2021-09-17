package agent

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/oshelper"

	e2e "github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/internal"
	"github.com/stretchr/testify/require"
)

func TestAgent_SetUp(t *testing.T) {
	agentDir := t.TempDir()
	err := oshelper.MakeDirs(0777, filepath.Join(agentDir, infraAgentDir))
	require.NoError(t, err)

	rootDir := t.TempDir()

	_, err = os.Create(filepath.Join(rootDir, "/nri-powerdns"))
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(rootDir, "/nri-powerdns-exporter"))
	require.NoError(t, err)
	_, err = os.Create(filepath.Join(rootDir, "/nri-prometheus"))
	require.NoError(t, err)

	settings, err := e2e.NewSettings(
		e2e.SettingsWithSpecPath("testdata/spec_file.yml"),
		e2e.SettingsWithAgentDir(agentDir),
		e2e.SettingsWithRootDir(rootDir),
	)
	require.NoError(t, err)

	t.Run("Given a scenario with 1 integration, the correct files should be in the AgentDir", func(t *testing.T) {
		sut := NewAgent(settings)
		require.NotEmpty(t, sut)

		err := sut.SetUp(settings.Spec().Scenarios[0])
		require.NoError(t, err)

		// nri-integration and exporter
		binaryFiles, err := ioutil.ReadDir(filepath.Join(agentDir, infraAgentDir, integrationsBinDir))
		require.NoError(t, err)
		require.Equal(t, 2, len(binaryFiles))

		// nri-prometheus
		exporterFiles, err := ioutil.ReadDir(filepath.Join(agentDir, infraAgentDir, exportersDir))
		require.NoError(t, err)
		require.Equal(t, 1, len(exporterFiles))

		// config file
		configFiles, err := ioutil.ReadDir(filepath.Join(agentDir, infraAgentDir, integrationsCfgDir))
		require.NoError(t, err)
		require.Equal(t, 1, len(configFiles))
	})
}
