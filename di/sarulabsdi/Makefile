GO_DIR ?= $(shell pwd)
GO_PKG ?= $(shell go list -e -f "{{ .ImportPath }}")

GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0

override GO111MODULE = on

.PHONY: build
build: ## Build all binaries
	@echo "Build binaries"
	@mkdir -p $(GO_DIR)/bin
	@errors=$$(GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} \
        go build -o "$(GO_DIR)/bin/app"); if [ "$${errors}" != "" ]; then echo "$${errors}"; fi

.PHONY: generate
generate: generate-code

.PHONY: generate-code
generate-code: ## Run generate code
	@echo "Generate code"
	@go generate ./...

.PHONY: test
test: ## Run unit tests
	@echo "Run unit tests"
	@go test -v ./...

.PHONY: fix
fix: fix-format fix-import

.PHONY: fix-import
fix-import: ## Fix imports in code
	@echo "Fix imports"
	@errors=$$(goimports -l -w -local $(GO_PKG) $$(go list -f "{{ .Dir }}" ./...)); if [ "$${errors}" != "" ]; then echo "$${errors}"; fi

.PHONY: fix-format
fix-format: ## Fix format in code
	@echo "Fix formatting"
	@gofmt -w ${GO_FMT_FLAGS} $$(go list -f "{{ .Dir }}" ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; fi

.PHONY: lint
lint: lint-format lint-import lint-style

.PHONY: lint-format
lint-format:
	@echo "Check formatting"
	@errors=$$(gofmt -l ${GO_FMT_FLAGS} $$(go list -f "{{ .Dir }}" ./...)); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

.PHONY: lint-import
lint-import:
	@echo "Check imports"
	@errors=$$(goimports -l -local $(GO_PKG) $$(go list -f "{{ .Dir }}" ./...)); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

.PHONY: lint-style
lint-style: ## execute linter
	@echo "Check code style"
	@errors=$$(golangci-lint run --no-config --issues-exit-code=0 --deadline=30m \
                             --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
                             --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
                             --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck \
                              $$(go list -f "{{ .Dir }}" ./...)); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

.PHONY: clean
clean:
	@echo "Cleanup"
	@rm -rf ${GO_DIR}/bin/*

.PHONY: tidy
tidy: ## Add missing and remove unused modules
	@echo 'run go mod tidy'
	@go mod tidy

.PHONY: vendor
vendor: ## Download modules in vendor folder
	@echo 'run go mod vendor'
	@go mod vendor

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
