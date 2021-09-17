package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseSpecFile(t *testing.T) {
	var sample = `
description: |
  End-to-end tests for PowerDNS integration

agent:
  integrations:
    nri-prometheus:  bin/nri-prometheus

scenarios:
  - description: |
      Scenario Description.
    before:
      - docker-compose -f deps/docker-compose.yml up -d
    after:
      - docker-compose -f deps/docker-compose.yml down -v
    integrations:
      - name: nri-powerdns
        binary_path: bin/nri-powerdns
        exporter_binary_path: bin/nri-powerdns-exporter
        config:
          powerdns_url: http://localhost:8081/api/v1/
    tests:
      nrqls:
        - query: "a-query"
      entities:
        - type: "POWERDNS_AUTHORITATIVE"
          data_type: "Metric"
          metric_name: "powerdns_authoritative_up"
      metrics:
        - source: "powerdns.yml"
          except:
            - powerdns_authoritative_answers_bytes_total
          additionals: ""`
	spec, err := ParseSpecFile([]byte(sample))
	assert.Nil(t, err)
	assert.Equal(t, "End-to-end tests for PowerDNS integration\n", spec.Description)

	expecedAgentOverrides := Agent{
		Integrations: map[string]string{
			"nri-prometheus": "bin/nri-prometheus",
		},
	}
	assert.Equal(t, &expecedAgentOverrides, spec.AgentOverrides)

	expectedScenarios := []Scenario{
		{
			Description: "Scenario Description.\n",
			Integrations: []Integration{
				{
					Name:               "nri-powerdns",
					BinaryPath:         "bin/nri-powerdns",
					ExporterBinaryPath: "bin/nri-powerdns-exporter",
					Config: map[string]interface{}{
						"powerdns_url": "http://localhost:8081/api/v1/",
					},
				},
			},
			Before: []string{"docker-compose -f deps/docker-compose.yml up -d"},
			After:  []string{"docker-compose -f deps/docker-compose.yml down -v"},
			Tests: Tests{
				NRQLs: []NRQL{{Query: "a-query"}},
				Entities: []Entity{
					{
						Type:       "POWERDNS_AUTHORITATIVE",
						DataType:   "Metric",
						MetricName: "powerdns_authoritative_up",
					},
				},
				Metrics: []Metrics{
					{
						Source: "powerdns.yml",
						Except: []string{"powerdns_authoritative_answers_bytes_total"},
					},
				},
			},
		},
	}
	assert.Equal(t, expectedScenarios, spec.Scenarios)
}
