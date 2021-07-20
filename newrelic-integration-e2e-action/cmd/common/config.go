package common

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	flagSpecPath    = "spec_path"
	flagVerboseMode = "verbose_mode"
)

type Config interface {
	SpecPath() string
	LogLevel() logrus.Level
	Validate() error
}

type config struct {
	specPath    string
	verboseMode bool
}

func (c *config) SpecPath() string {
	return c.specPath
}

func (c *config) LogLevel() logrus.Level {
	if c.verboseMode {
		return logrus.DebugLevel
	}
	return logrus.InfoLevel
}

func (c *config) Validate() error {
	if c.specPath == "" {
		return fmt.Errorf("missing required flag %s", flagSpecPath)
	}
	return nil
}

func LoadConfig() Config {
	specsPathPtr := flag.String(flagSpecPath, "", "Relative path to the spec file")
	verboseModePtr := flag.Bool(flagVerboseMode, false, "If true the debug level is enabled")
	flag.Parse()
	return &config{
		specPath:    *specsPathPtr,
		verboseMode: *verboseModePtr,
	}

}
