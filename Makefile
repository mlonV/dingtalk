.PHONY: all install build build-arm build-amd download gotool clean test lint help

BINARY = bin/dingtalk
TARGET = dingtalk.go

all: gotool build

install: build
	@cp $(BINARY) /usr/local/bin/

build: $(TARGET)
	@CGO_ENABLED=0 go build -ldflags "-s -w" -o $(BINARY) $(TARGET)

build-arm: $(TARGET)
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $(BINARY) $(TARGET)

build-amd: $(TARGET)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(BINARY) $(TARGET)

gotool:
	@go get -u golang.org/x/lint/golint

clean:
	@rm -f $(BINARY)

test:
	@go test -v ./...

lint:
	@golint -set_exit_status ./...

help:
	@echo "Available targets:"
	@echo "  all       - Build the project"
	@echo "  install   - Install the binary to /usr/local/bin/"
	@echo "  build     - Build the project for the current platform"
	@echo "  build-arm - Build the project for ARM64 Linux"
	@echo "  build-amd - Build the project for AMD64 Linux"
	@echo "  gotool    - Install necessary Go tools"
	@echo "  clean     - Clean up build artifacts"
	@echo "  test      - Run tests"
	@echo "  lint      - Run linter"
	@echo "  help      - Show this help message"