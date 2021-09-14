package spec

import yaml "gopkg.in/yaml.v3"

type Definition struct {
	Description    string     `yaml:"description"`
	Scenarios      []Scenario `yaml:"scenarios"`
	AgentOverrides *Agent     `yaml:"agent"`
}

func (def *Definition) Validate() error {
	for i := range def.Scenarios {
		if err := def.Scenarios[i].validate(); err != nil {
			return err
		}
	}
	return nil
}

type Agent struct {
	Integrations map[string]string `yaml:"integrations"`
}

type Scenario struct {
	Description  string        `yaml:"description"`
	Integrations []Integration `yaml:"integrations"`
	Before       []string      `yaml:"before"`
	After        []string      `yaml:"after"`
	Tests        Tests         `yaml:"tests"`
}

func (s *Scenario) validate() error {
	for i := range s.Integrations {
		if err := s.Integrations[i].validate(); err != nil {
			return err
		}
	}
	return nil
}

type Integration struct {
	Name               string                 `yaml:"name"`
	BinaryPath         string                 `yaml:"binary_path"`
	ExporterBinaryPath string                 `yaml:"exporter_binary_path"`
	Config             map[string]interface{} `yaml:"config"`
}

type Metrics struct {
	Source    string   `yaml:"source"`
	Including []string `yaml:"including"`
	Except    []string `yaml:"except"`
}
type Tests struct {
	NRQLs    []NRQL    `yaml:"nrqls"`
	Entities []Entity  `yaml:"entities"`
	Metrics  []Metrics `yaml:"metrics"`
}

type NRQL struct {
	Query string `yaml:"query"`
}

type Entity struct {
	Type       string `yaml:"type"`
	DataType   string `yaml:"data_type"`
	MetricName string `yaml:"metric_name"`
}

func (i *Integration) validate() error {
	return nil
}

func ParseSpecFile(content []byte) (*Definition, error) {
	specDefinition := &Definition{}
	if err := yaml.Unmarshal(content, specDefinition); err != nil {
		return nil, err
	}
	return specDefinition, nil
}
