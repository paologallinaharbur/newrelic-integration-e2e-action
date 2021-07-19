DOCKER_IMG=newrelic-integration-e2e-action:spec-validator
SPEC_PATH=?

all: compile run

compile:
	docker build -t $(DOCKER_IMG) .
run:
	docker run $(DOCKER_IMG) --spec_path=$(SPEC_PATH)