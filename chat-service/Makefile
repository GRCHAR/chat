# Makefile for Chat Service

.PHONY: build run clean test deps migrate help

# 变量定义
APP_NAME := chat-service
BUILD_DIR := bin
MAIN_DIR := cmd/server

# Go相关参数
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# 构建目标
build: deps
	@echo "构建应用..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) ./$(MAIN_DIR)
	@echo "构建完成: $(BUILD_DIR)/$(APP_NAME)"

# 运行应用
run: deps
	@echo "启动应用..."
	$(GOCMD) run ./$(MAIN_DIR)/main.go

# 清理构建文件
clean:
	@echo "清理构建文件..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 安装依赖
deps:
	@echo "下载依赖..."
	$(GOMOD) download
	$(GOMOD) tidy

# 运行测试
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 格式化代码
fmt:
	@echo "格式化代码..."
	$(GOCMD) fmt ./...

# 代码检查
lint:
	@echo "运行代码检查..."
	golangci-lint run

# 数据库迁移
migrate-up:
	@echo "执行数据库迁移..."
	$(GOCMD) run $(MAIN_DIR)/migrate.go up

migrate-down:
	@echo "回滚数据库迁移..."
	$(GOCMD) run $(MAIN_DIR)/migrate.go down

# Docker相关命令
docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(APP_NAME):latest .

docker-run:
	@echo "运行Docker容器..."
	docker run -p 8080:8080 $(APP_NAME):latest

# 生产环境构建
build-prod:
	@echo "生产环境构建..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./$(MAIN_DIR)
	@echo "生产环境构建完成: $(BUILD_DIR)/$(APP_NAME)-linux-amd64"

# 开发环境启动
dev: deps
	@echo "启动开发环境..."
	$(GOCMD) run ./$(MAIN_DIR)/main.go

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  build         - 构建应用"
	@echo "  run           - 运行应用"
	@echo "  clean         - 清理构建文件"
	@echo "  deps          - 下载依赖"
	@echo "  test          - 运行测试"
	@echo "  test-coverage - 生成测试覆盖率报告"
	@echo "  fmt           - 格式化代码"
	@echo "  lint          - 代码检查"
	@echo "  migrate-up    - 执行数据库迁移"
	@echo "  migrate-down  - 回滚数据库迁移"
	@echo "  docker-build  - 构建Docker镜像"
	@echo "  docker-run    - 运行Docker容器"
	@echo "  build-prod    - 生产环境构建"
	@echo "  dev           - 开发环境启动"
	@echo "  help          - 显示此帮助信息"
