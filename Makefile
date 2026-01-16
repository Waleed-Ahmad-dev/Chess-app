# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=chess
BINARY_UNIX=$(BINARY_NAME)_unix

# WebAssembly parameters
WASM_DIR=web
WASM_MAIN=cmd/wasm/main.go
WASM_OUT=$(WASM_DIR)/main.wasm

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/chess

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(WASM_OUT)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/chess
	./$(BINARY_NAME)

# WASM Build Commands
init-wasm:
	cp $$(go env GOROOT)/lib/wasm/wasm_exec.js $(WASM_DIR)/

build-wasm: init-wasm
	GOOS=js GOARCH=wasm $(GOBUILD) -o $(WASM_OUT) $(WASM_MAIN)

serve-wasm: build-wasm
	@echo "Serving WASM at http://localhost:8080"
	# Using a simple python server for serving static files if standard go tool isn't preferred for serving
	# Or we can use the main app's web server mode if it supports serving the static directory
	# For independent WASM testing, a simple static server is often easiest:
	goexec 'http.ListenAndServe(":8080", http.FileServer(http.Dir("./web")))' || \
	python3 -m http.server 8080 --directory ./web

help:
	@echo "Makefile commands:"
	@echo "  make build       - Build the CLI application"
	@echo "  make run         - Run the CLI application"
	@echo "  make build-wasm  - Build the WebAssembly application"
	@echo "  make serve-wasm  - Build and serve the WASM application"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
