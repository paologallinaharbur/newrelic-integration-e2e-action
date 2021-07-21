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

type Integrations []Integration

func (i Integrations) MarshalYAML() (interface{}, error) {
	type integration struct {
		Name   string                 `yaml:"name"`
		Config map[string]interface{} `yaml:"config"`
	}
	type outputIntegrations struct {
		Integrations []integration `yaml:"integrations"`
	}
	out := outputIntegrations{
		Integrations:make([]integration, len(i)),
	}
	for index:=range i{
		out.Integrations[index] = integration{
			Name:   i[index].Name,
			Config: i[index].Config,
		}
	}
	content, err := yaml.Marshal(out)
	if err != nil {
		return nil, err
	}

	return string(content), nil
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
