PHONY: all

GOCMD=go
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GOCOVER=$(GOCMD) tool cover

# Common flags
GOTEST_NOCACHEFLAG=-count=1
GOTEST_COVERFILE=cover.out
GOTEST_COVERFLAG=-coverprofile=$(GOTEST_COVERFILE)

GOLANG_CI_LINT_VERSION=1.27.0

all: test lint

build:
	$(GOBUILD) ./...

test:
	$(GOTEST) ./... $(GOTEST_NOCACHEFLAG) $(GOTEST_COVERFLAG)

bench:
	$(GOTEST) ./... -bench $(GOTEST_NOCACHEFLAG) $(GOTEST_COVERFLAG)

cover: test
	$(GOCOVER) -html=$(GOTEST_COVERFILE)

lint-setup:
	@# Make sure linter is up to date
	$(eval CURRENT_VERSION := $(strip $(shell ${GOPATH}/bin/golangci-lint version 2>&1 | sed 's/[^0-9.]*\([0-9.]*\).*/\1/')))
	@if [ "$(CURRENT_VERSION)" != "$(GOLANG_CI_LINT_VERSION)" ]; then \
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin v$(GOLANG_CI_LINT_VERSION) ; \
	fi

lint: lint-setup
	${GOPATH}/bin/golangci-lint run ./...
