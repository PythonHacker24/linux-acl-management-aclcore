PROJECT_NAME := aclcore
GO_PACKAGES   := ./cmd/... ./internal/...
BUILD_DIR     := bin

CORE_BIN       := $(BUILD_DIR)/aclcore

.PHONY: all
all: fmt vet build

.PHONY: build
build: $(CORE_BIN)

$(CORE_BIN):
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/aclcore

.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	go fmt $(GO_PACKAGES)

.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet $(GO_PACKAGES)

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: install
install: build
	@echo "Installing binaries to /usr/local/bin..."
	install -m 755 $(API_BIN)  /usr/local/bin/aclcore

.PHONY: docker-api

docker-core:
	docker build -f Dockerfile.core  -t aclcore:latest .

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

.PHONY: run-core

run-core: $(API_CORE)
	@echo "Running aclcore (as aclcore user)..."
	su -s /bin/bash aclcore -c "$(API_BIN)"

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  all         : fmt, vet, build"
	@echo "  build       : build both binaries"
	@echo "  fmt         : format Go code"
	@echo "  vet         : run go vet"
	@echo "  lint        : run golangci-lint"
	@echo "  test        : run tests"
	@echo "  install     : install binaries to /usr/local/bin"
	@echo "  docker-core : build Docker image for core daemon"
	@echo "  clean       : remove build artifacts"
	@echo "  run-core    : run core daemon as root"
	@echo "  help        : this help message"
