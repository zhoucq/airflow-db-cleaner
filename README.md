# Airflow DB Cleaner

[![CI](https://github.com/zhoucq/airflow-db-cleaner/actions/workflows/ci.yml/badge.svg)](https://github.com/zhoucq/airflow-db-cleaner/actions/workflows/ci.yml)
[![Release](https://github.com/zhoucq/airflow-db-cleaner/actions/workflows/release.yml/badge.svg)](https://github.com/zhoucq/airflow-db-cleaner/actions/workflows/release.yml)

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

## CI/CD

### Continuous Integration

This project uses GitHub Actions for continuous integration. The CI workflow runs on every push to the main branch and on pull requests:

- Runs all tests to ensure code quality
- Builds the project to verify it compiles correctly

### Releases

This project uses GitHub Actions to automatically create releases when a new tag is pushed.

To create a new release:

```bash
# Tag a new version
git tag v1.0.0

# Push the tag to GitHub
git push origin v1.0.0
```

This will trigger the release workflow which builds the project for all supported platforms and creates a GitHub release with the artifacts.

## Cross-platform Support

This tool supports building and running on different platforms:

- Can be developed on Mac ARM (M series chips)
- Can build binaries for Linux x86_64 servers
- Use `make build-linux` command to generate Linux version directly
- All build artifacts are placed in the `bin` directory

## License

MIT 
