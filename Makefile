.PHONY: all

GOLANG_CI_LINT_VERSION=1.59.1
GO_BIN=$(shell pwd)/.build
export

all: test lint

build:
	@echo "Building AWS RDS Store"
	cd store/awsrds && $(MAKE) build
	@echo "Building Vault store"
	cd store/vault && $(MAKE) build
	@echo "Building main module"
	go build ./...

test:
	@echo "Testing AWS RDS store"
	cd store/awsrds && $(MAKE) test
	@echo "Testing Vault store"
	cd store/vault && $(MAKE) test
	@echo "Testing main module"
	go test ./... -count=1 -coverprofile=cover.out

lint: lint-setup
	@echo "Linting AWS RDS store"
	cd store/awsrds && $(MAKE) lint
	@echo "Linting Vault store"
	cd store/vault && $(MAKE) lint
	@echo "Linting main module"
	"$(GO_BIN)/golangci-lint" run ./...

bench:
	@echo "Benching AWS RDS store"
	cd store/awsrds && $(MAKE) bench
	@echo "Benching Vault store"
	cd store/vault && $(MAKE) bench
	@echo "Benching main module"
	go test ./... -bench -count=1 -coverprofile=cover.out

cover: test
	mkdir -p "$(GO_BIN)/coverage"
	@echo "Generating coverage for AWS RDS store"
	cd store/awsrds && $(MAKE) cover
	@echo "Generating coverage for Vault store"
	cd store/vault && $(MAKE) cover
	@echo "Generating coverage for main module"
	go tool cover -html=cover.out -o "$(GO_BIN)/coverage/main.html"

lint-setup: _bin
	@# Make sure linter is up to date
	$(eval CURRENT_VERSION := $(strip $(shell $(GO_BIN)/bin/golangci-lint version 2>&1 | sed 's/[^0-9.]*\([0-9.]*\).*/\1/')))
	@if [ "$(CURRENT_VERSION)" != "$(GOLANG_CI_LINT_VERSION)" ]; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(GO_BIN)" "v$(GOLANG_CI_LINT_VERSION)" ; \
	fi

_bin:
	mkdir -p "$(GO_BIN)"

format:
	gofmt -w -e .

tidy:
	go mod tidy
	cd store/awsrds && go mod tidy
	cd store/vault && go mod tidy
