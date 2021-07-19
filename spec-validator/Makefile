DOCKER_IMG=newrelic-integration-e2e-action:spec-validator
SPECS=?

all: compile run

compile:
	docker build -t $(DOCKER_IMG) .
run:
	docker run $(DOCKER_IMG) -e SPECS=$(SPECS)