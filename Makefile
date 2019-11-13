built_at := $(shell date +%s)
git_commit := $(shell git describe --dirty --always)

BIN:=./bin

OS := $(shell uname)
GOLANGCI_LINT_VERSION?=1.19.1
ifeq ($(OS),Darwin)
	GOLANGCI_LINT_ARCHIVE=golangci-lint-$(GOLANGCI_LINT_VERSION)-darwin-amd64.tar.gz
else
	GOLANGCI_LINT_ARCHIVE=golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64.tar.gz
endif

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "-X github.com/90poe/connectctl/pkg/version.GitHash=$(git_commit) -X github.com/90poe/connectctl/pkg/version.BuildDate=$(built_at)" ./cmd/connectctl

.PHONY: local-release
local-release:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: release
release: test lint
	curl -sL https://git.io/goreleaser | bash

.PHONY: test
test:
	@go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: ci
ci: build test lint

.PHONY: lint
lint: $(BIN)/golangci-lint/golangci-lint ## lint
	$(BIN)/golangci-lint/golangci-lint run

$(BIN)/golangci-lint/golangci-lint:
	curl -OL https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VERSION)/$(GOLANGCI_LINT_ARCHIVE)
	mkdir -p $(BIN)/golangci-lint/
	tar -xf $(GOLANGCI_LINT_ARCHIVE) --strip-components=1 -C $(BIN)/golangci-lint/
	chmod +x $(BIN)/golangci-lint
	rm -f $(GOLANGCI_LINT_ARCHIVE)
