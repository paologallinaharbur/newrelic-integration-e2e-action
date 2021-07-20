package spec

import (
	"gopkg.in/yaml.v3"
)

type SpecDefinition struct {
	Description string     `yaml:"description"`
	BeforeAll   string     `yaml:"before_all"`
	AfterAll    string     `yaml:"before_all"`
	Scenarios   []Scenario `yaml:"scenarios"`
}

func (s *SpecDefinition) Validate() error {
	for i := range s.Scenarios {
		if err := s.Scenarios[i].validate(); err != nil {
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
	Config map[string]interface{} `yaml:"config"`
}

func (i *Integration) validate() error {
	return nil
}

func ParseSpecFile(content []byte) (*SpecDefinition, error) {
	specDefinition := &SpecDefinition{}
	if err := yaml.Unmarshal(content, specDefinition); err != nil {
		return nil, err
	}
	return specDefinition, nil
}
