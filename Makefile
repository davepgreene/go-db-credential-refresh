PHONY: all

GOLANG_CI_LINT_VERSION=1.45.2
GO_BIN=$(shell pwd)/.build
all: test lint

build:
	go build ./...

test:
	go test ./... -count=1 -coverprofile=cover.out

bench:
	go test ./... -bench -count=1 -coverprofile=cover.out

cover: test
	go tool cover -html=cover.out

lint-setup: _bin
	@# Make sure linter is up to date
	$(eval CURRENT_VERSION := $(strip $(shell $(GO_BIN)/bin/golangci-lint version 2>&1 | sed 's/[^0-9.]*\([0-9.]*\).*/\1/')))
	@if [ "$(CURRENT_VERSION)" != "$(GOLANG_CI_LINT_VERSION)" ]; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(GO_BIN)" "v$(GOLANG_CI_LINT_VERSION)" ; \
	fi

lint: lint-setup
	"$(GO_BIN)/golangci-lint" run ./...

_bin:
	mkdir -p "$(GO_BIN)"
