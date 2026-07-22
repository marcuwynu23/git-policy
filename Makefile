BIN_DIR=dist
.PHONY: all build clean test lint vet fmt run install uninstall doctor validate help install-binary dev link

BINARY_NAME=git-policy
BINARY_EXT=
ifdef ComSpec
  BINARY_EXT=.exe
endif
GO_BUILD=CGO_ENABLED=0 go build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GO_FLAGS=-ldflags="-w -X github.com/marcuwynu23/git-policy/cmd.version=$(VERSION)"

all: lint vet test build

build:
	mkdir -p dist
	$(GO_BUILD) $(GO_FLAGS) -o dist/$(BINARY_NAME)$(BINARY_EXT) .

clean:
	rm -f $(BINARY_NAME)$(BINARY_EXT)
	rm -rf dist/

test:
	go test -v -count=1 ./...

test-short:
	go test -short -count=1 ./...

test-cover:
	go test -v -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	}
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy
	go mod verify

run:
	go run . run

install:
	go run . install

uninstall:
	go run . uninstall

doctor:
	go run . doctor

validate:
	go run . validate

version:
	go run . version

.PHONY: dist dist-windows dist-linux dist-darwin
dist:
	$(GO_BUILD) $(GO_FLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=linux GOARCH=amd64 $(GO_BUILD) $(GO_FLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO_BUILD) $(GO_FLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) $(GO_FLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO_BUILD) $(GO_FLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .

install-binary: build
	cp -f dist/$(BINARY_NAME)$(BINARY_EXT) $(BIN_DIR)/$(BINARY_NAME)$(BINARY_EXT)

link: build
	@echo "Creating symbolic link..."
	@# Check if we're on Windows or Unix-like
	@if [ -n "$(ComSpec)" ]; then \
		powershell -Command "if (Test-Path 'C:\Bin\tools\$(BINARY_NAME)$(BINARY_EXT)') { Remove-Item -Force 'C:\Bin\tools\$(BINARY_NAME)$(BINARY_EXT)' }"; \
		powershell -Command "New-Item -ItemType SymbolicLink -Path 'C:\Bin\tools' -Name '$(BINARY_NAME)$(BINARY_EXT)' -Target '$(CURDIR)\dist\$(BINARY_NAME)$(BINARY_EXT)' -Force"; \
		echo "Symbolic link created: C:\Bin\tools\$(BINARY_NAME)$(BINARY_EXT) -> $(CURDIR)\dist\$(BINARY_NAME)$(BINARY_EXT)"; \
	else \
		ln -sf $(CURDIR)/dist/$(BINARY_NAME)$(BINARY_EXT) /usr/local/bin/$(BINARY_NAME)$(BINARY_EXT); \
		echo "Symbolic link created: /usr/local/bin/$(BINARY_NAME)$(BINARY_EXT) -> $(CURDIR)/dist/$(BINARY_NAME)$(BINARY_EXT)"; \
	fi

dev: build link install
	@echo "git-policy built, linked to PATH, and hooks installed."

help:
	@echo "Usage:"
	@echo "  make build         - Build the binary"
	@echo "  make install-binary - Build + copy binary to $(BIN_DIR)"
	@echo "  make link          - Build + create symlink in $(BIN_DIR) (live testing)"
	@echo "  make dev           - Build + symlink + install hooks (full dev setup)"
	@echo "  make test          - Run all tests"
	@echo "  make test-cover    - Run tests with coverage"
	@echo "  make lint          - Run linter"
	@echo "  make vet           - Run go vet"
	@echo "  make fmt           - Format code"
	@echo "  make tidy          - Tidy modules"
	@echo "  make run           - Run git-policy"
	@echo "  make install       - Install hooks"
	@echo "  make dist          - Build all platforms"
	@echo "  make clean         - Clean build artifacts"
