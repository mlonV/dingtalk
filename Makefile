# Go 应用名称
APP_NAME := dingtalk

# 版本号
VERSION := 1.0.0

# 生成的可执行文件存储目录
BUILD_DIR := build

# 默认的平台，如果不指定则默认是当前系统的架构
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# 获取当前的操作系统和架构
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)
# make build GOOS=linux GOARCH=arm64
# make build GOOS=linux GOARCH=amd64



# 根据平台的不同生成适配的文件后缀（Windows需要.exe）
ifeq ($(GOOS), windows)
    EXT := .exe
else
    EXT :=
endif

# 默认目标：构建指定架构的二进制文件
build: clean
	@echo "Building $(APP_NAME) for GOOS=$(GOOS), GOARCH=$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(APP_NAME)-$(GOOS)-$(GOARCH)-$(VERSION)$(EXT) dingtalk.go
	@echo "Built $(APP_NAME) for $(GOOS)/$(GOARCH)"

# 清理构建文件
clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf $(BUILD_DIR)

.PHONY: build clean
