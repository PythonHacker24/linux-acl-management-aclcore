PROJECT_NAME := aclcore
GO_PACKAGES  := ./cmd/... ./internal/...
BUILD_DIR    := bin
TAR_BUILD_DIR := build

CORE_BIN     := $(BUILD_DIR)/$(PROJECT_NAME)

TARGETS := \
	linux_amd64 \
	linux_arm64

.PHONY: all
all: fmt vet build

.PHONY: build
build: $(CORE_BIN)

$(CORE_BIN):
	@echo "Building $(PROJECT_NAME)... (online)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/$(PROJECT_NAME)

.PHONY: build-offline
build-offline: vendor
	@echo "Building $(PROJECT_NAME)... (offline using vendor)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -mod=vendor -o $(CORE_BIN) ./cmd/$(PROJECT_NAME)

.PHONY: build-cross
build-cross:
	@echo "Cross building for: $(TARGETS) (online)"
	@mkdir -p $(BUILD_DIR)
	@for target in $(TARGETS); do \
		OS=$${target%_*}; \
		ARCH=$${target#*_}; \
		OUT=$(BUILD_DIR)/$(PROJECT_NAME)-$$OS-$$ARCH; \
		echo "Building $$OUT..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -o $$OUT ./cmd/$(PROJECT_NAME); \
	done

.PHONY: vendor
vendor:
	@echo "Vendoring dependencies..."
	go mod vendor

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
	@echo "Installing binary to /usr/local/bin..."
	install -m 755 $(CORE_BIN) /usr/local/bin/$(PROJECT_NAME)

.PHONY: docker-core
docker-core:
	docker build -f Dockerfile.core -t $(PROJECT_NAME):latest .

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf vendor

.PHONY: run-core
run-core: $(CORE_BIN)
	@echo "Running $(PROJECT_NAME) as user 'aclcore'..."
	su -s /bin/bash aclcore -c "$(CORE_BIN)"

.PHONY: package
package: vendor
	@echo "Packaging project source with vendor..."
	@mkdir -p $(TAR_BUILD_DIR)
	@TMP_DIR=$$(mktemp -d); \
	NAME=$(PROJECT_NAME)-source; \
	echo "Copying files to $$TMP_DIR/$$NAME..."; \
	mkdir -p $$TMP_DIR/$$NAME; \
	rsync -a --exclude '$(TAR_BUILD_DIR)' --exclude '*.tar.gz' ./ $$TMP_DIR/$$NAME; \
	TARBALL=$(PROJECT_NAME)-source.tar.gz; \
	tar -czf $(TAR_BUILD_DIR)/$$TARBALL -C $$TMP_DIR $$NAME; \
	echo "Created $(TAR_BUILD_DIR)/$$TARBALL"; \
	rm -rf $$TMP_DIR

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  all            : fmt, vet, build"
	@echo "  build          : build for local OS/arch (online)"
	@echo "  build-offline  : build using vendor (offline)"
	@echo "  build-cross    : cross build for $(TARGETS) (online)"
	@echo "  vendor         : vendor dependencies"
	@echo "  package        : create tarball of project + vendor"
	@echo "  fmt            : format Go code"
	@echo "  vet            : run go vet"
	@echo "  lint           : run golangci-lint"
	@echo "  test           : run tests"
	@echo "  install        : install binary to /usr/local/bin"
	@echo "  docker-core    : build Docker image for core daemon"
	@echo "  clean          : clean build artifacts and vendor"
	@echo "  run-core       : run core daemon as 'aclcore' user"
	@echo "  help           : show this help"
