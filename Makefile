.PHONY: build install clean test run

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Linker flags to set version info
LDFLAGS := -ldflags "-X github.com/riyasyash/nishku/cmd.Version=$(VERSION) \
                     -X github.com/riyasyash/nishku/cmd.GitCommit=$(GIT_COMMIT) \
                     -X github.com/riyasyash/nishku/cmd.BuildDate=$(BUILD_DATE)"

# Build the binary
build:
	go build $(LDFLAGS) -buildvcs=false -o nishku .

# Install to /usr/local/bin
install: build
	sudo mv nishku /usr/local/bin/

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/nishku-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/nishku-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/nishku-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/nishku-windows-amd64.exe .
	@echo "Binaries built in dist/"

# Create release archives
release: build-all
	@echo "Creating release archives..."
	@mkdir -p dist/releases
	cd dist && tar -czf releases/nishku-$(VERSION)-darwin-amd64.tar.gz nishku-darwin-amd64 -C .. README.md
	cd dist && tar -czf releases/nishku-$(VERSION)-darwin-arm64.tar.gz nishku-darwin-arm64 -C .. README.md
	cd dist && tar -czf releases/nishku-$(VERSION)-linux-amd64.tar.gz nishku-linux-amd64 -C .. README.md
	cd dist && zip -q releases/nishku-$(VERSION)-windows-amd64.zip nishku-windows-amd64.exe -j ../README.md
	@echo "Release archives created in dist/releases/"
	@ls -lh dist/releases/

# Clean build artifacts
clean:
	rm -f nishku
	rm -rf dist/

# Run tests
test:
	go test ./...

# Run the app locally
run:
	go run . $(ARGS)

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Show version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(GIT_COMMIT)"
	@echo "Date:    $(BUILD_DATE)"
