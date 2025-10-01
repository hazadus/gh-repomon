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
    go fmt ./...

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
    GOOS=linux GOARCH=amd64 go build -o bin/gh-repomon-linux-amd64 ./cmd/repomon
    GOOS=linux GOARCH=arm64 go build -o bin/gh-repomon-linux-arm64 ./cmd/repomon
    GOOS=darwin GOARCH=amd64 go build -o bin/gh-repomon-darwin-amd64 ./cmd/repomon
    GOOS=darwin GOARCH=arm64 go build -o bin/gh-repomon-darwin-arm64 ./cmd/repomon
    GOOS=windows GOARCH=amd64 go build -o bin/gh-repomon-windows-amd64.exe ./cmd/repomon
    GOOS=windows GOARCH=arm64 go build -o bin/gh-repomon-windows-arm64.exe ./cmd/repomon
    @echo "Built binaries for all platforms in bin/"

# Generate commit message (see https://github.com/hazadus/gh-commitmsg)
commitmsg:
    gh commitmsg --language english --examples

# Generate code line count report
cloc:
    cloc --fullpath --exclude-list-file=.clocignore --md . > cloc.md
