package settings

import (
	"io/ioutil"
	"path/filepath"

	"github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/spec"
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
	rootDir       string
	accountID     int
	apiKey        string
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

func WithRootDir(rootDir string) Option {
	return func(o *settingOptions) {
		o.rootDir = rootDir
	}
}

func WithAccountID(accountID int) Option {
	return func(o *settingOptions) {
		o.accountID = accountID
	}
}

func WithApiKey(apiKey string) Option {
	return func(o *settingOptions) {
		o.apiKey = apiKey
	}
}

type Settings interface {
	Logger() *logrus.Logger
	Spec() *spec.Definition
	AgentDir() string
	RootDir() string
	LicenseKey() string
	ApiKey() string
	AccountID() int
}

type settings struct {
	logger        *logrus.Logger
	spec          *spec.Definition
	specParentDir string
	agentDir      string
	licenseKey    string
	accountID     int
	apiKey        string
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

func (s *settings) ApiKey() string {
	return s.apiKey
}

func (s *settings) AccountID() int {
	return s.accountID
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
	s, err := spec.ParseSpecFile(content)
	if err != nil {
		return nil, err
	}
	logger.Debug("return with settings")
	return &settings{
		logger:        logger,
		spec:          s,
		agentDir:      options.agentDir,
		specParentDir: options.specParentDir,
		licenseKey:    options.licenseKey,
		apiKey:        options.apiKey,
		accountID:     options.accountID,
	}, nil
}
