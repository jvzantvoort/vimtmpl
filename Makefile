NAME := git-ci
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'
GOIMPORTS ?= goimports
GOCILINT ?= golangci-lint
GO ?= GO111MODULE=on go
.DEFAULT_GOAL := help

.PHONY: fmt
fmt: ## Formatting source codes.
	find . -type f -name '*.go' -not -path '*/vendor/*' -exec $(GOIMPORTS) -w "{}" \;

.PHONY: clean
clean:
	@rm -f $(COMMANDS) || true

.PHONY: refresh
refresh: tags
	@go-bindata -pkg vimtmpl messages; \

.PHONY: tags
tags:
	@find "$${PWD}" -type f -name '*.go' -not -path '*/vendor/*'| sed "s,$${PWD}/,," | xargs gotags >tags

.PHONY: pretty
pretty:
	@find "$${PWD}" -type f -name '*.go' -not -path '*/vendor/*' -exec goimports -w "{}" \;; \
	find "$${PWD}" -type f -name '*.go' -not -path '*/vendor/*' -exec gofmt -w "{}" \;

.PHONY: lint
lint: ## Run golint and go vet.
	@$(GOCILINT) run --no-config --disable-all --enable=goimports --enable=misspell ./...

.PHONY: test
test:  ## Run the tests.
	@$(GO) test ./...

.PHONY: build
build: main.go  ## Build a binary.
	$(foreach cmd,$(COMMANDS), $(GO) build -ldflags "$(LDFLAGS)" ./cmd/$(cmd);)

.PHONY: cross
cross: main.go  ## Build binaries for cross platform.
	mkdir -p pkg
	@# darwin
	@for arch in "amd64"; do \
		GOOS=darwin GOARCH=$${arch} make build; \
		zip pkg/git-ci_$(VERSION)_darwin_$${arch}.zip git-ci; \
	done;
	@# linux
	@for arch in "amd64" "386" "arm64" "arm"; do \
		GOOS=linux GOARCH=$${arch} make build; \
		zip pkg/git-ci_$(VERSION)_linux_$${arch}.zip git-ci; \
	done;
