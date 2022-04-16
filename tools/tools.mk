.PHONY: all

all: test lint

ifndef MODULE
	$(error You must define the MODULE variable)
endif

ifdef GO_BIN
	ifeq ($(strip $(GO_BIN)),)
		$(error GO_BIN must not be empty)
	endif
else
	$(error You must define the GO_BIN variable)
endif

build:
	go build ./...

test:
	go test ./... -count=1 -coverprofile=cover.out

bench:
	go test ./... -bench -count=1 -coverprofile=cover.out

cover: test
	go tool cover -html=cover.out -o "$(GO_BIN)/coverage/$(MODULE).html"

lint:
	"$(GO_BIN)/golangci-lint" run ./...
