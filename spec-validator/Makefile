DOCKER_IMG=newrelic-integration-e2e-action:spec-validator
ROOT_DIR=?
SPECS_PATH=?

all: compile run

compile:
	docker build -t $(DOCKER_IMG) .
run:
	docker run $(DOCKER_IMG) --root_dir=$(ROOT_DIR) --specs_path=$(SPECS_PATH)