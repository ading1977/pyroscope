# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

Pyroscope is a continuous profiling platform built with a microservices architecture that can run as a single binary or distributed system. Key architectural components:

- **Distributor**: Receives and forwards profiling data to ingesters
- **Ingester**: Stores profiling data in memory and eventually writes to long-term storage  
- **Store Gateway**: Serves queries against long-term storage blocks
- **Querier**: Coordinates queries across ingesters and store gateways
- **Query Frontend**: Provides query optimization and caching
- **Compactor**: Compacts and processes long-term storage blocks
- **Metastore**: Manages metadata with Raft consensus for distributed deployments

The system supports multiple profiling formats (pprof, JFR, OTLP) and integrates with object storage backends (S3, GCS, etc.) for long-term data retention.

## Technology Stack

- **Backend**: Go 1.23+ with gRPC/Connect for APIs
- **Frontend**: React/TypeScript with Webpack build system
- **Storage**: Parquet files for long-term storage, in-memory for recent data
- **Consensus**: HashiCorp Raft for metastore coordination
- **Serialization**: Protocol Buffers with vtproto for performance

## Build Commands

### Go Backend
```bash
# Build all binaries (requires frontend build)
make build

# Build without frontend (development)
make build-dev

# Build specific binaries
make go/bin-pyroscope
make go/bin-profilecli

# Debug builds
make go/bin-debug
```

### Frontend
```bash
# Production build (in Docker)
make frontend/build

# Development server
yarn dev

# Combined development (frontend + backend)
yarn backend:dev
```

### Testing
```bash
# Run all Go tests
make go/test

# Run specific package tests
go test ./pkg/distributor/...

# Run frontend tests
yarn test

# Run integration tests
make examples/test
```

### Linting & Formatting
```bash
# Lint all code
make lint

# Auto-fix formatting issues  
make fmt

# Individual linters
make go/lint
make buf/lint
```

## Development Workflow

### Code Generation
After modifying protobuf definitions or adding parquet tags:
```bash
make generate
```

### Running Locally
```bash
# Single binary mode
./pyroscope --config.file ./cmd/pyroscope/pyroscope.yaml

# With embedded Grafana
./pyroscope --target all,embedded-grafana

# Use make target with parameters
make run PARAMS="--config.file ./cmd/pyroscope/pyroscope.yaml"
```

### Docker Development
```bash
# Build Docker image
make docker-image/pyroscope/build

# Multi-arch builds
make docker-image/pyroscope/build-multiarch

# Debug image with delve
make docker-image/pyroscope/build-debug
```

### Package Management
```bash
# Update Go dependencies
make go/mod

# Update specific dependency
go get example.com/module@version
make go/mod
```

## Testing Guidelines

- Unit tests are co-located with source code (`*_test.go`)
- Integration tests are in `pkg/test/integration/`
- Example tests are in `examples/examples_test.go`
- Always include tests for new functionality and bug fixes
- Use table-driven tests for multiple test cases
- Use testify/assert for assertions

## Key Packages

- `pkg/pyroscope/`: Main server component and configuration
- `pkg/distributor/`: Ingestion and data distribution logic
- `pkg/ingester/`: In-memory storage and query handling
- `pkg/phlaredb/`: Long-term storage engine with Parquet
- `pkg/querier/`: Query coordination and processing
- `pkg/frontend/`: Query frontend and optimization
- `pkg/model/`: Core data structures and protobuf definitions
- `pkg/pprof/`: pprof format handling and utilities
- `pkg/objstore/`: Object storage abstractions and providers

## Configuration

Main config file: `cmd/pyroscope/pyroscope.yaml`
Components are configured via YAML with support for:
- Environment variable substitution
- Runtime configuration updates
- Per-tenant overrides
- Feature flags

## Frontend Development

Located in `public/app/` with:
- React 18 + TypeScript
- Redux Toolkit for state management
- Webpack for bundling
- Grafana UI components library

Development server runs on `:4041`, backend typically on `:4040`.

## Important Notes

- Always run `make generate` after protobuf changes
- The codebase uses workspaces - run `go work sync` after dependency updates
- Frontend requires Node v18 and Yarn v1.22
- Docker builds require the frontend to be built first
- Use `gotestsum` for test runs (included in make targets)
- Follow conventional commit message format for consistency