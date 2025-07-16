# Default target
all: build

WORK_DIR ?= ./_work
BIN_DIR ?= $(WORK_DIR)/bin
KEYS_DIR ?= $(WORK_DIR)/keys
BIN_NAME ?= data-api-server

# Build the Go application from example/main.go
build:
	@mkdir -p $(BIN_DIR)
	@rm -f $(BIN_DIR)/$(BIN_NAME)
	@echo "Building Go application..."
	@echo "Using binary name: $(BIN_NAME)"
	@echo "Output directory: $(BIN_DIR)"
	@echo "Compiling..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BIN_NAME) .
	@ls -lh $(BIN_DIR)/$(BIN_NAME)

# Run the Go application with optional ARGS
run: build
	@$(BIN_DIR)/$(BIN_NAME) $(ARGS)

 
# Remove built binary
clean:
	@rm -f $(BIN_DIR)/$(BIN_NAME)
