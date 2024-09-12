# Name of the binary
BINARY_NAME=service

# Build and binary settings
BUILD_DIR=bin
BUILD_PATH=$(BUILD_DIR)/$(BINARY_NAME)

# Default Go commands
GO=go
GOREBUILD=go build -o $(BUILD_PATH)

# Build the Go binary
build:
	@echo "Building the binary..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=arm $(GO) build -o $(BUILD_PATH) ./main.go
	@echo "Build complete: $(BUILD_PATH)"

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -rf docker/$(BUILD_DIR)
	@echo "Clean complete."

# Lint Go code (requires golint to be installed)
lint:
	@echo "Linting Go code..."

# Test Go code
test:
	@echo "Running tests..."
	@$(GO) test ./...
	@echo "Tests complete."

container: lint test build
	mkdir -p docker/bin
	cp bin/service docker/bin
	docker build -t service:0.1 docker
