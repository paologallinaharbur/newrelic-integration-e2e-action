NAME=newrelic-integration-e2e

GOLANGCI_LINT = github.com/golangci/golangci-lint/cmd/golangci-lint
GOFLAGS          = -mod=readonly

# Variables
SPEC_PATH=?
ROOT_DIR=?
VERBOSE=?
LICENSE_KEY=?
GH_ACTION_DIR=$(CURDIR)/../

all: clean validate test docker-build

clean:
	@echo "=== $(NAME) === [ clean ]: removing binaries and coverage file..."
	@rm -rfv bin coverage.xml $(TARGET)

validate:
	@printf "=== $(NAME) === [ validate ]: running golangci-lint & semgrep... "
	@go run  $(GOFLAGS) $(GOLANGCI_LINT) run --verbose
	@[ -f .semgrep.yml ] && semgrep_config=".semgrep.yml" || semgrep_config="p/golang" ; \
	docker run --rm -v "${PWD}:/src:ro" --workdir / returntocorp/semgrep -c "$$semgrep_config"


run:
	@printf "=== $(NAME) === [ run / $* ]: running the binary "
	go run $(GH_ACTION_DIR)/$(NAME)/main.go --root_dir=${GH_ACTION_DIR} --agent_dir=$(GH_ACTION_DIR)/agent_dir --license_key=$(LICENSE_KEY) --spec_path=$(SPEC_PATH) --verbose_mode=$(VERBOSE)

test:
	@echo "=== $(NAME) === [ test ]: running unit tests..."
	@go test -race ./... -count=1
