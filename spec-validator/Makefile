DOCKER_IMG=newrelic-integration-e2e-action:spec-validator
SPEC_PATH=?
ROOT_DIR=?
VERBOSE=?

all: compile run

compile:
	docker build -t $(DOCKER_IMG) .
run:
	docker run -v $(ROOT_DIR):/opt/code $(DOCKER_IMG) --spec_path=/opt/code/$(SPEC_PATH) --verbose_mode=$(VERBOSE)