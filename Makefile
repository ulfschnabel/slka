BINARY_READ = slka-read
BINARY_WRITE = slka-write
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.Version=$(VERSION)"

.PHONY: all build test test-unit test-integration test-all clean install lint skill

all: build

build: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

build-linux-amd64:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_READ)-linux-amd64 ./cmd/slka-read
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_WRITE)-linux-amd64 ./cmd/slka-write

build-linux-arm64:
	mkdir -p dist
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_READ)-linux-arm64 ./cmd/slka-read
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_WRITE)-linux-arm64 ./cmd/slka-write

build-darwin-amd64:
	mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_READ)-darwin-amd64 ./cmd/slka-read
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_WRITE)-darwin-amd64 ./cmd/slka-write

build-darwin-arm64:
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_READ)-darwin-arm64 ./cmd/slka-read
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_WRITE)-darwin-arm64 ./cmd/slka-write

build-local:
	mkdir -p dist
	go build $(LDFLAGS) -o dist/$(BINARY_READ) ./cmd/slka-read
	go build $(LDFLAGS) -o dist/$(BINARY_WRITE) ./cmd/slka-write

test: test-unit

test-unit:
	@echo "Running unit tests..."
	go test -v -race -cover ./internal/... ./pkg/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./test/integration/...

test-all: test-unit test-integration
	@echo "All tests passed!"

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	golangci-lint run --timeout=5m

clean:
	rm -rf dist/
	rm -f coverage.out coverage.html

install: build-local
	cp dist/$(BINARY_READ) $(GOPATH)/bin/
	cp dist/$(BINARY_WRITE) $(GOPATH)/bin/
	@echo "Installed to $(GOPATH)/bin"

deps:
	go mod download
	go mod tidy

skill:
	cd skill && ./package.sh

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build             - Build for all platforms"
	@echo "  make build-local       - Build for current platform"
	@echo "  make test              - Run unit tests (alias for test-unit)"
	@echo "  make test-unit         - Run unit tests only"
	@echo "  make test-integration  - Run integration tests with mock Slack server"
	@echo "  make test-all          - Run all tests (unit + integration)"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make lint              - Run linter"
	@echo "  make clean             - Remove build artifacts"
	@echo "  make install           - Install to GOPATH/bin"
	@echo "  make deps              - Download dependencies"
	@echo "  make skill             - Package the agent skill"
