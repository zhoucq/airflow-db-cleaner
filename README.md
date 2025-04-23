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
./airflow-db-cleaner

# 指定配置文件运行
./airflow-db-cleaner --config /path/to/config.yaml
```

## 构建

```bash
go build -o airflow-db-cleaner ./cmd/airflow-db-cleaner
```

## 许可证

MIT 