package agent

import "github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/pkg/spec"

type agentIntegration struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}
type agentIntegrationsList struct {
	Integrations []agentIntegration `yaml:"integrations"`
}

func createAgentIntegrationModel(integrations []spec.Integration) *agentIntegrationsList {
	out := &agentIntegrationsList{
		Integrations: make([]agentIntegration, len(integrations)),
	}
	for index := range integrations {
		integration := integrations[index]
		out.Integrations[index] = agentIntegration{
			Name:   integration.Name,
			Config: integration.Config,
		}
	}
	return out

}
