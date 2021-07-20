package settings

import (
	"io/ioutil"

	"github.com/newrelic/newrelic-integration-e2e-action/pkg/spec"
	"github.com/sirupsen/logrus"
)

var defaultSettingsOptions = settingOptions{
	logLevel: logrus.InfoLevel,
}

type settingOptions struct {
	logLevel logrus.Level
	specPath string
}

type Option func(*settingOptions)

func WithSpecPath(specPath string) Option {
	return func(o *settingOptions) {
		o.specPath = specPath
	}
}

func WithLogLevel(logLevel logrus.Level) Option {
	return func(o *settingOptions) {
		o.logLevel = logLevel
	}
}

type Settings interface {
	Logger() *logrus.Logger
	SpecDefinition() *spec.SpecDefinition
}

type settings struct {
	logger *logrus.Logger
	spec   *spec.SpecDefinition
}

func (s *settings) Logger() *logrus.Logger {
	return s.logger
}

func (s *settings) SpecDefinition() *spec.SpecDefinition {
	return s.spec
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
	return &settings{
		logger: logger,
		spec:   spec,
	}, nil
}