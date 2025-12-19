.PHONY: build clean install test help

# Build variables
BINARY_NAME=taskflow
BUILD_DIR=build
MAIN_FILE=main.go

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: ./$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@echo "Multi-platform build complete in $(BUILD_DIR)/"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Install to system
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "Install complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run example
run-example:
	@echo "Running example workflow..."
	@./$(BINARY_NAME) run examples/simple.yaml

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run || echo "golangci-lint not installed, skipping..."

# Show help
help:
	@echo "TaskFlow Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build         Build the binary"
	@echo "  make build-all     Build for multiple platforms"
	@echo "  make deps          Install dependencies"
	@echo "  make clean         Remove build artifacts"
	@echo "  make install       Install to /usr/local/bin"
	@echo "  make test          Run tests"
	@echo "  make run-example   Run an example workflow"
	@echo "  make fmt           Format code"
	@echo "  make lint          Lint code"
	@echo "  make help          Show this help message"
