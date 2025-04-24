.PHONY: build clean run test build-linux release

# 二进制文件名
BINARY=airflow-db-cleaner
# 输出目录
OUTPUT_DIR=bin
# 发布包输出目录
RELEASE_DIR=release
# 配置文件路径
CONFIG_FILE=config/config.yaml

# 确保输出目录存在
ensure_dir:
	mkdir -p $(OUTPUT_DIR)

# 确保发布目录存在
ensure_release_dir:
	mkdir -p $(RELEASE_DIR)

# 构建目标
build: ensure_dir
	go build -o $(OUTPUT_DIR)/$(BINARY) .

# 支持的操作系统和架构
OSES=linux darwin windows
ARCHES=amd64 arm64 386

# 交叉编译Linux x86_64版本
build-linux: ensure_dir
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY)-linux-amd64 .

# 交叉编译指定OS和ARCH版本
build-os-arch: ensure_dir
	GOOS=$(OS) GOARCH=$(ARCH) go build -o $(OUTPUT_DIR)/$(BINARY)-$(OS)-$(ARCH) .

# 带版本信息的构建
build-release: ensure_dir
	go build -ldflags "-X main.version=$$(git describe --tags --always)" -o $(OUTPUT_DIR)/$(BINARY) .

# 带版本信息的Linux x86_64构建
build-release-linux: ensure_dir
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$$(git describe --tags --always)" -o $(OUTPUT_DIR)/$(BINARY)-linux-amd64 .

# 构建多平台版本
build-all: build build-linux

# 构建所有支持的OS和ARCH组合
build-all-os-arch: ensure_dir
	@for os in $(OSES); do \
		for arch in $(ARCHES); do \
			echo "Building for $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build -o $(OUTPUT_DIR)/$(BINARY)-$$os-$$arch .; \
		done; \
	done

# 创建单个发布包
release-os-arch: ensure_release_dir
	@echo "Creating release package for $(OS)/$(ARCH)..."
	@mkdir -p $(RELEASE_DIR)/$(BINARY)-$(OS)-$(ARCH)/config
	@GOOS=$(OS) GOARCH=$(ARCH) go build -o $(RELEASE_DIR)/$(BINARY)-$(OS)-$(ARCH)/$(BINARY) .
	@cp $(CONFIG_FILE) $(RELEASE_DIR)/$(BINARY)-$(OS)-$(ARCH)/config/
	@if [ "$(OS)" = "windows" ]; then \
		cd $(RELEASE_DIR) && zip -r $(BINARY)-$(OS)-$(ARCH).zip $(BINARY)-$(OS)-$(ARCH); \
	else \
		cd $(RELEASE_DIR) && tar -czf $(BINARY)-$(OS)-$(ARCH).tar.gz $(BINARY)-$(OS)-$(ARCH); \
	fi
	@echo "Release package created at $(RELEASE_DIR)/$(BINARY)-$(OS)-$(ARCH).tar.gz or .zip"

# 创建所有发布包
release-all: ensure_release_dir
	@for os in $(OSES); do \
		for arch in $(ARCHES); do \
			echo "Creating release package for $$os/$$arch..."; \
			mkdir -p $(RELEASE_DIR)/$(BINARY)-$$os-$$arch/config; \
			GOOS=$$os GOARCH=$$arch go build -o $(RELEASE_DIR)/$(BINARY)-$$os-$$arch/$(BINARY) .; \
			cp $(CONFIG_FILE) $(RELEASE_DIR)/$(BINARY)-$$os-$$arch/config/; \
			if [ "$$os" = "windows" ]; then \
				cd $(RELEASE_DIR) && zip -r $(BINARY)-$$os-$$arch.zip $(BINARY)-$$os-$$arch; \
			else \
				cd $(RELEASE_DIR) && tar -czf $(BINARY)-$$os-$$arch.tar.gz $(BINARY)-$$os-$$arch; \
			fi; \
			echo "Release package created for $$os/$$arch"; \
		done; \
	done
	@echo "All release packages created in $(RELEASE_DIR)"

# 清理构建产物
clean:
	rm -rf $(OUTPUT_DIR)
	rm -rf $(RELEASE_DIR)

# 运行程序（默认配置）
run: build
	$(OUTPUT_DIR)/$(BINARY)

# 运行程序（指定配置文件）
run-config: build
	$(OUTPUT_DIR)/$(BINARY) --config $(CONFIG)

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
	@echo "  make build-os-arch OS=linux ARCH=arm64 - 构建指定OS和架构的应用"
	@echo "  make build-all-os-arch   - 构建所有支持的OS和架构组合"
	@echo "  make build-release       - 构建带版本信息的当前平台应用"
	@echo "  make build-release-linux - 构建带版本信息的Linux x86_64平台应用"
	@echo "  make release-os-arch OS=linux ARCH=amd64 - 创建指定OS和架构的发布包"
	@echo "  make release-all         - 创建所有支持的OS和架构的发布包"
	@echo "  make clean               - 清理构建产物"
	@echo "  make run                 - 构建并运行（默认配置）"
	@echo "  make run-config CONFIG=path/to/config.yaml - 构建并运行（指定配置）"
	@echo "  make test                - 运行测试"
	@echo "  make deps                - 更新依赖" 
