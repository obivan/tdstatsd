.PHONY: release
release:
	goreleaser --rm-dist

.PHONY: lint
lint:
	golangci-lint run
