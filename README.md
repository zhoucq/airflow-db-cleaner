# Airflow DB Cleaner

一个用于清除Airflow数据库中历史数据的工具，以提高Airflow的性能。

## 功能

- 清理过期的DAG运行记录
- 清理过期的任务实例
- 清理过期的日志
- 支持自定义清理策略和保留期

## 安装

```bash
go get github.com/zhoucq/airflow-db-cleaner
```

## 使用方法

### 配置

在`config`目录中编辑配置文件，设置数据库连接和清理策略。

### 运行

```bash
# 使用默认配置运行
./bin/airflow-db-cleaner

# 指定配置文件运行
./bin/airflow-db-cleaner --config /path/to/config.yaml
```

## 构建

所有构建产物将输出到 `bin` 目录中：

```bash
# 构建当前平台版本
make build

# 构建 Linux x86_64 版本
make build-linux

# 同时构建多个平台版本
make build-all

# 带版本信息的构建
make build-release

# 清理所有构建产物
make clean
```

## 跨平台支持

本工具支持在不同平台上构建和运行：

- 可以在 Mac ARM (M系列芯片) 上开发
- 可以为 Linux x86_64 服务器构建二进制文件
- 使用 `make build-linux` 命令可以直接生成 Linux 版本
- 所有构建产物都会放在 `bin` 目录下

## 许可证

MIT 