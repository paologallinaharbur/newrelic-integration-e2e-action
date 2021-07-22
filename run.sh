#!/bin/bash

LICENSE_KEY=$1
SPEC_PATH=$(PWD)/samples/nri-powerdns/powerdns-e2e.yml
make -C newrelic-integration-e2e LICENSE_KEY="${LICENSE_KEY}" ROOT_DIR=$(PWD) SPEC_PATH=${SPEC_PATH} VERBOSE=true run