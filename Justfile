# Default recipe - show available commands
default:
    @just --list

# Run the application
run *ARGS:
    go run ./cmd/repomon {{ARGS}}

# Build the binary
build:
    go build -o bin/gh-repomon ./cmd/repomon

# Format code
format:
    gofmt -s -w .

# Run linter
lint:
    golangci-lint run ./...

# Run tests
test:
    go test -v -race ./...

# Run tests with coverage
test-coverage:
    go test -v -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

# Run integration tests only
test-integration:
    go test -v -race -tags=integration ./test/integration/...

# Run all tests (unit + integration)
test-all:
    go test -v -race ./...
    go test -v -race -tags=integration ./test/integration/...

# Clean build artifacts
clean:
    rm -rf bin/
    rm -f coverage.out coverage.html

# Install/update dependencies
deps:
    go mod download
    go mod tidy

# Build for all platforms
build-all:
    @echo "Building binaries for all platforms..."
    @mkdir -p bin
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/gh-repomon-linux-amd64 ./cmd/repomon
    GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/gh-repomon-linux-arm64 ./cmd/repomon
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/gh-repomon-darwin-amd64 ./cmd/repomon
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/gh-repomon-darwin-arm64 ./cmd/repomon
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/gh-repomon-windows-amd64.exe ./cmd/repomon
    GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o bin/gh-repomon-windows-arm64.exe ./cmd/repomon
    @echo "✓ Built binaries for all platforms in bin/"

# Create a local release build with checksums
release:
    @echo "Creating local release build..."
    @just clean
    @just test
    @just build-all
    @cd bin && sha256sum * > checksums.txt
    @echo "✓ Release build created in bin/ with checksums"
    @echo ""
    @echo "Files created:"
    @ls -lh bin/

# Test goreleaser configuration without publishing
release-test:
    @echo "Testing goreleaser configuration..."
    goreleaser release --snapshot --clean
    @echo "✓ Test release created in dist/"

# Generate commit message (see https://github.com/hazadus/gh-commitmsg)
commitmsg:
    gh commitmsg --language english --examples

# Generate code line count report
cloc:
    cloc --fullpath --exclude-list-file=.clocignore --md . > cloc.md
