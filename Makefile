.PHONY: all help
.DEFAULT_GOAL := help

# COLORS
YELLOW = \033[33m
GREEN  = \033[32m
WHITE  = \033[37m
RESET  = \033[0m

GOLANG_CI_LINT_VERSION=2.4.0
GO_BIN=$(shell pwd)/.build
export

all: test lint

build: ## Build library and sample stores
	@printf "$(GREEN)Building AWS RDS Store$(RESET)\n"
	@cd store/awsrds && $(MAKE) -s build
	@printf "$(GREEN)Building Vault store$(RESET)\n"
	@cd store/vault && $(MAKE) -s build
	@printf "$(GREEN)Building main module$(RESET)\n"
	@go build ./...

test: ## Test everything
	@printf "$(GREEN)Testing AWS RDS store$(RESET)\n"
	@cd store/awsrds && $(MAKE) -s test
	@printf "\n$(GREEN)Testing Vault store$(RESET)\n"
	@cd store/vault && $(MAKE) -s test
	@printf "\n$(GREEN)Testing main module$(RESET)\n"
	@go test github.com/davepgreene/go-db-credential-refresh/driver -count=1 -coverprofile=cover.out

lint: lint-setup ## Lint everything
	@printf "$(GREEN)Linting AWS RDS store$(RESET)\n"
	@cd store/awsrds && $(MAKE) -s lint
	@printf "\n$(GREEN)Linting Vault store$(RESET)\n"
	@cd store/vault && $(MAKE) -s lint
	@printf "\n$(GREEN)Linting main module$(RESET)\n"
	@"$(GO_BIN)/golangci-lint" run ./...

bench: ## Go test with benchmarks
	@printf "$(GREEN)Benching AWS RDS store$(RESET)\n"
	@cd store/awsrds && $(MAKE) -s bench
	@printf "\n$(GREEN)Benching Vault store$(RESET)\n"
	@cd store/vault && $(MAKE) -s bench
	@printf "\n$(GREEN)Benching main module$(RESET)\n"
	@go test ./... -bench -count=1 -coverprofile=cover.out

cover: test ## Go test with coverage
	@mkdir -p "$(GO_BIN)/coverage"
	@echo "Generating coverage for AWS RDS store"
	@cd store/awsrds && $(MAKE) -s cover
	@echo "Generating coverage for Vault store"
	@cd store/vault && $(MAKE) -s cover
	@echo "Generating coverage for main module"
	@go tool cover -html=cover.out -o "$(GO_BIN)/coverage/main.html"

lint-setup: _bin
	@# Make sure linter is up to date
	$(eval CURRENT_VERSION := $(strip $(shell $(GO_BIN)/golangci-lint version 2>&1 | sed 's/[^0-9.]*\([0-9.]*\).*/\1/')))
	@if [ "$(CURRENT_VERSION)" != "$(GOLANG_CI_LINT_VERSION)" ]; then \
		echo "Updating golangci-lint from $(CURRENT_VERSION) to version $(GOLANG_CI_LINT_VERSION)"; \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(GO_BIN)" "v$(GOLANG_CI_LINT_VERSION)" ; \
	fi

_bin:
	@mkdir -p "$(GO_BIN)"

format: ## Gofmt
	@gofmt -w -e .

tidy: ## Tidy up go modules
	@go mod tidy
	@cd store/awsrds && go mod tidy
	@cd store/vault && go mod tidy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
