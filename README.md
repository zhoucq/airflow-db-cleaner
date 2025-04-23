# Airflow DB Cleaner

A tool for cleaning historical data in Airflow database to improve Airflow performance.

## Features

- Clean expired DAG run records
- Clean expired task instances
- Clean expired logs
- Support custom cleaning strategies and retention periods

## Installation

```bash
go get github.com/zhoucq/airflow-db-cleaner
```

## Usage

### Configuration

Edit the configuration file in the `config` directory to set database connections and cleaning strategies.

### Running

```bash
# Run with default configuration
./bin/airflow-db-cleaner

# Run with specified configuration file
./bin/airflow-db-cleaner --config /path/to/config.yaml
```

## Build

All build artifacts will be output to the `bin` directory:

```bash
# Build for current platform
make build

# Build for Linux x86_64
make build-linux

# Build for multiple platforms simultaneously
make build-all

# Build with version information
make build-release

# Clean all build artifacts
make clean
```

## Cross-platform Support

This tool supports building and running on different platforms:

- Can be developed on Mac ARM (M series chips)
- Can build binaries for Linux x86_64 servers
- Use `make build-linux` command to generate Linux version directly
- All build artifacts are placed in the `bin` directory

## License

MIT 