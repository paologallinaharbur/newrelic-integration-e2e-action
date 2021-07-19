package pkg

import "gopkg.in/yaml.v3"

type SpecDefinition struct {
	Description string `yaml:"description"`
}

func ParseSpecFile(content []byte) (*SpecDefinition, error) {
	specDefinition := &SpecDefinition{}
	if err := yaml.Unmarshal(content, specDefinition); err != nil {
		return nil, err
	}
	return specDefinition, nil
}
