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

build-default:
	go build ./...

test-default:
	go test ./... -count=1 -coverprofile=cover.out

bench-default:
	go test ./... -bench -count=1 -coverprofile=cover.out

cover-default: test-default
	go tool cover -html=cover.out -o "$(GO_BIN)/coverage/$(MODULE).html"

lint-default:
	"$(GO_BIN)/golangci-lint" run ./...

%:  %-default
	@  true
