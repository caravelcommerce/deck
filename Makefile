.PHONY: build install clean test

# Build the binary
build:
	go build -o bin/deck main.go

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o bin/deck-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/deck-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/deck-linux-amd64 main.go

# Install to /usr/local/bin
install: build
	sudo cp bin/deck /usr/local/bin/deck
	sudo chmod +x /usr/local/bin/deck
	@echo "✅ Deck installed successfully to /usr/local/bin/deck"

# Uninstall from /usr/local/bin
uninstall:
	sudo rm -f /usr/local/bin/deck
	@echo "✅ Deck uninstalled successfully"

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Download dependencies
deps:
	go mod download
	go mod tidy
