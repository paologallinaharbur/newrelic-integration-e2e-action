NAME=newrelic-integration-e2e-action/spec-validator
DOCKER_IMG=newrelic-integration-e2e-action:spec-validator

GOLANGCI_LINT = github.com/golangci/golangci-lint/cmd/golangci-lint
GOFLAGS          = -mod=readonly

# Variables
SPEC_PATH=?
ROOT_DIR=?
VERBOSE=?

all: clean validate test docker-build

clean:
	@echo "=== $(NAME) === [ clean ]: removing binaries and coverage file..."
	@rm -rfv bin coverage.xml $(TARGET)

validate:
	@printf "=== $(NAME) === [ validate ]: running golangci-lint & semgrep... "
	@go run  $(GOFLAGS) $(GOLANGCI_LINT) run --verbose
	@[ -f .semgrep.yml ] && semgrep_config=".semgrep.yml" || semgrep_config="p/golang" ; \
	docker run --rm -v "${PWD}:/src:ro" --workdir / returntocorp/semgrep -c "$$semgrep_config"

docker-build:
	@printf "=== $(NAME) === [ build ]: running docker build... "
	docker build -t $(DOCKER_IMG) .

docker-run: docker-build
	@printf "=== $(NAME) === [ run ]: running the application "
	docker run -v $(ROOT_DIR):/opt/code $(DOCKER_IMG) --spec_path=/opt/code/$(SPEC_PATH) --verbose_mode=$(VERBOSE)


test:
	@echo "=== $(NAME) === [ test ]: running unit tests..."
	@go test -race ./... -count=1
