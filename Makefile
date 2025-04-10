# Variables
APP_NAME = execution-service
CMD_DIR = ./cmd/main.go
BUILD_DIR = ./bin

# Default target
.PHONY: all
all: build

# Build the project
.PHONY: build
build:
    @echo "Building the project..."
    @mkdir -p $(BUILD_DIR)
    @go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)
    @echo "Build complete. Binary is located at $(BUILD_DIR)/$(APP_NAME)"

# Run the project
.PHONY: run
run:
    @echo "Running the project..."
    @go run $(CMD_DIR)

# Clean up build artifacts
.PHONY: clean
clean:
    @echo "Cleaning up..."
    @rm -rf $(BUILD_DIR)
    @echo "Clean complete."

# Format the code
.PHONY: fmt
fmt:
    @echo "Formatting code..."
    @go fmt ./...

# Install dependencies
.PHONY: deps
deps:
    @echo "Installing dependencies..."
    @go mod tidy

# Test the project
.PHONY: test
test:
    @echo "Running tests..."
    @go test ./... -v