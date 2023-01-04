BIN?=lmi

# Utility settings
TOOLS_DIR := .tools
GOLANGCI_LINT_VERSION = v1.50.1

# Container build settings
CONTAINER_BUILD_CMD?=docker build

# Container settings
CONTAINER_REPO?=ghcr.io/infratographer/lmi
LMI_CONTAINER_IMAGE_NAME = $(CONTAINER_REPO)/lmi
CONTAINER_TAG?=latest

# OpenAPI
OAPI_CODEGEN_CMD?=go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen

## Targets

.PHONY: build
build:
	go build -o $(BIN) ./main.go

.PHONY: test
test:
	@echo Running unit tests...
	@go test -timeout 30s -cover -short  -tags testtools ./...

.PHONY: coverage
coverage:
	@echo Generating coverage report...
	@go test -timeout 30s -tags testtools ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

lint: golint

golint: | vendor $(TOOLS_DIR)/golangci-lint
	@echo Linting Go files...
	@$(TOOLS_DIR)/golangci-lint run

clean:
	@echo Cleaning...
	@rm -rf coverage.out
	@go clean -testcache
	@rm -r $(TOOLS_DIR)

vendor:
	@go mod download
	@go mod tidy

image: lmi-image

lmi-image:
	$(CONTAINER_BUILD_CMD) -f images/lmi/Dockerfile . -t $(LMI_CONTAINER_IMAGE_NAME):$(CONTAINER_TAG)

.PHONY: generate
generate: openapi

.PHONY: openapi
openapi: openapi-types openapi-client openapi-spec

.PHONY: openapi-types
openapi-types:
	@echo Generating OpenAPI types...
	@$(OAPI_CODEGEN_CMD) -package v1 \
		-generate types \
		-o api/v1/types.gen.go openapi-v1.yaml

# Note that due to a limitation in oapi-codegen, we need to generate the client
# in the same package as the types.
.PHONY: openapi-client
openapi-client:
	@echo Generating OpenAPI client...
	@$(OAPI_CODEGEN_CMD) -package v1 \
		-generate client \
		-o api/v1/client.gen.go openapi-v1.yaml

.PHONY: openapi-spec
openapi-spec:
	@echo Generating OpenAPI spec...
	@$(OAPI_CODEGEN_CMD) -package v1 \
		-generate spec \
		-o api/v1/openapi.gen.go openapi-v1.yaml

# Tools setup
$(TOOLS_DIR):
	mkdir -p $(TOOLS_DIR)

$(TOOLS_DIR)/golangci-lint: $(TOOLS_DIR)
	export \
		VERSION=$(GOLANGCI_LINT_VERSION) \
		URL=https://raw.githubusercontent.com/golangci/golangci-lint \
		BINDIR=$(TOOLS_DIR) && \
	curl -sfL $$URL/$$VERSION/install.sh | sh -s $$VERSION
	$(TOOLS_DIR)/golangci-lint version
	$(TOOLS_DIR)/golangci-lint linters
