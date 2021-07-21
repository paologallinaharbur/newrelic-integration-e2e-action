package settings

import (
	"io/ioutil"
	"path/filepath"

	"github.com/newrelic/newrelic-integration-e2e/pkg/spec"
	"github.com/sirupsen/logrus"
)

var defaultSettingsOptions = settingOptions{
	logLevel: logrus.InfoLevel,
}

type settingOptions struct {
	logLevel      logrus.Level
	specPath      string
	specParentDir string
	licenseKey    string
	agentDir      string
}

type Option func(*settingOptions)

func WithSpecPath(specPath string) Option {
	return func(o *settingOptions) {
		o.specPath = specPath
		o.specParentDir = filepath.Dir(specPath)
	}
}

func WithLogLevel(logLevel logrus.Level) Option {
	return func(o *settingOptions) {
		o.logLevel = logLevel
	}
}

func WithLicenseKey(licenseKey string) Option {
	return func(o *settingOptions) {
		o.licenseKey = licenseKey
	}
}

func WithAgentDir(agentDir string) Option {
	return func(o *settingOptions) {
		o.agentDir = agentDir
	}
}

type Settings interface {
	Logger() *logrus.Logger
	Spec() *spec.Definition
	AgentDir() string
	RootDir() string
	LicenseKey() string
}

type settings struct {
	logger        *logrus.Logger
	spec          *spec.Definition
	specParentDir string
	agentDir      string
	licenseKey    string
}

func (s *settings) Logger() *logrus.Logger {
	return s.logger
}

func (s *settings) LicenseKey() string {
	return s.licenseKey
}

func (s *settings) Spec() *spec.Definition {
	return s.spec
}

func (s *settings) AgentDir() string {
	return s.agentDir
}

func (s *settings) RootDir() string {
	return s.specParentDir
}

// New returns a Scheduler
func New(
	opts ...Option) (Settings, error) {
	options := defaultSettingsOptions
	for _, opt := range opts {
		opt(&options)
	}
	logger := logrus.New()
	logger.SetLevel(options.logLevel)
	content, err := ioutil.ReadFile(options.specPath)
	if err != nil {
		return nil, err
	}
	logger.Debug("parsing the content of the spec file")
	spec, err := spec.ParseSpecFile(content)
	if err != nil {
		return nil, err
	}
	logger.Debug("return with settings")
	return &settings{
		logger:        logger,
		spec:          spec,
		agentDir:      options.agentDir,
		specParentDir: options.specParentDir,
		licenseKey:    options.licenseKey,
	}, nil
}
