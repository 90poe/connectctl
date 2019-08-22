built_at := $(shell date +%s)
git_commit := $(shell git describe --dirty --always)

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "-X github.com/90poe/connectctl/pkg/version.GitHash=$(git_commit) -X github.com/90poe/connectctl/pkg/version.BuildDate=$(built_at)" ./cmd/connectctl

.PHONY: install-deps
install-deps:
	go mod download
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.17.1

.PHONY: local-release
local-release:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: test
test:
	@go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: lint
lint:
	./bin/golangci-lint run

.PHONY: ci
ci: build test lint
