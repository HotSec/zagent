# 定义GOPATH和项目目录
GOPATH ?= $(shell go env GOPATH)
PROJECT_DIR ?= $(shell pwd)

# 定义输出目录
BIN_DIR := $(PROJECT_DIR)/build

VERSION = 0.1.0
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)
gitBranch = $(shell git rev-parse --abbrev-ref HEAD) 

ldflags="-w -X main.Version=${VERSION} -X main.gitBranch=$(gitBranch) -X main.gitTag=${gitTag} -X main.buildDate=${buildDate} -X main.gitCommit=${gitCommit} -X main.gitTreeState=${gitTreeState}"


# 定义编译目标
all: build

# 编译目标
build:
	@mkdir -p $(BIN_DIR)
	go build -v -ldflags ${ldflags} -o $(BIN_DIR)/zagent cmd/zagent/zagent.go

# 测试目标
test:
	go test -v ./

# 格式化代码目标
fmt:
	gofmt -s -w ./

# 清理目标
clean:
	rm -rf $(BIN_DIR)

# 帮助信息
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build - Build the application"
	@echo "  test  - Run tests"
	@echo "  fmt   - Format Go source code"
	@echo "  clean - Clean up build artifacts"

.PHONY: build test fmt clean help
