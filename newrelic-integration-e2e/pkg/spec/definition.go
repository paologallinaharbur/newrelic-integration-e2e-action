package spec

import (
	"gopkg.in/yaml.v3"
)

type Definition struct {
	Description string     `yaml:"description"`
	BeforeAll   string     `yaml:"before_all"`
	AfterAll    string     `yaml:"after_all"`
	Scenarios   []Scenario `yaml:"scenarios"`
}

func (def *Definition) Validate() error {
	for i := range def.Scenarios {
		if err := def.Scenarios[i].validate(); err != nil {
			return err
		}
	}
	return nil
}

type Scenario struct {
	Description  string        `yaml:"description"`
	Integrations []Integration `yaml:"integrations"`
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
	Name   string                 `yaml:"name"`
	Path   string                 `yaml:"path"`
	Config map[string]interface{} `yaml:"config"`
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
