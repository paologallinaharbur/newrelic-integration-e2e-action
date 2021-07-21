#!/bin/bash

LICENSE_KEY=$1
SPEC_PATH=/Users/icorrales/Repositories/github.com/newrelic/newrelic-integration-e2e-action/samples/powerdns_e2e.yml
make -C newrelic-integration-e2e LICENSE_KEY="${LICENSE_KEY}" ROOT_DIR=$(PWD) SPEC_PATH=${SPEC_PATH} VERBOSE=true run