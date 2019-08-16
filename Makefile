

.PHONY: install-deps
install-deps:
	go mod download
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.17.1

.PHONY: local-release
local-release:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: lint
lint:
	./bin/golangci-lint run