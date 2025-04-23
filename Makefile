.PHONY: build clean run test build-linux

# 二进制文件名
BINARY=airflow-db-cleaner

# 构建目标
build:
	go build -o $(BINARY) ./cmd/airflow-db-cleaner

# 交叉编译Linux x86_64版本
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux-amd64 ./cmd/airflow-db-cleaner

# 带版本信息的构建
build-release:
	go build -ldflags "-X main.version=$$(git describe --tags --always)" -o $(BINARY) ./cmd/airflow-db-cleaner

# 带版本信息的Linux x86_64构建
build-release-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$$(git describe --tags --always)" -o $(BINARY)-linux-amd64 ./cmd/airflow-db-cleaner

# 构建多平台版本
build-all: build build-linux

# 清理构建产物
clean:
	rm -f $(BINARY) $(BINARY)-linux-amd64

# 运行程序（默认配置）
run: build
	./$(BINARY)

# 运行程序（指定配置文件）
run-config: build
	./$(BINARY) --config $(CONFIG)

# 测试
test:
	go test -v ./...

# 安装依赖
deps:
	go mod tidy

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  make build               - 构建当前平台应用"
	@echo "  make build-linux         - 构建Linux x86_64平台应用"
	@echo "  make build-all           - 构建所有平台应用"
	@echo "  make build-release       - 构建带版本信息的当前平台应用"
	@echo "  make build-release-linux - 构建带版本信息的Linux x86_64平台应用"
	@echo "  make clean               - 清理构建产物"
	@echo "  make run                 - 构建并运行（默认配置）"
	@echo "  make run-config CONFIG=path/to/config.yaml - 构建并运行（指定配置）"
	@echo "  make test                - 运行测试"
	@echo "  make deps                - 更新依赖" 