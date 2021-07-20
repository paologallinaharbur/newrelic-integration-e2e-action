package executor

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	integrationsCfgDir = "integrations.d"
	integrationsBinDir = "bin"
)

type agent struct {
	rootDir string
}

func (a *agent) initialize(logger *logrus.Logger) error {
	logger.Debugf("setup the workspace on %s", a.rootDir)
	if err := removeDirectoryContent(a.rootDir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(a.rootDir, integrationsCfgDir), fs.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(a.rootDir, integrationsBinDir), fs.ModePerm); err != nil {
		return err
	}
	return nil
}

func (a *agent) copyIntegrationBinaries(integrationPath string) error {
	if _, err := copyFile(integrationPath, filepath.Join(a.rootDir, integrationsBinDir, filepath.Base(integrationPath))); err != nil {
		return err
	}
	return nil
}
