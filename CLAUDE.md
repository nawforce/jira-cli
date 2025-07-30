# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

JiraCLI is a feature-rich interactive command line tool for Atlassian Jira written in Go. It provides an alternative to the Jira UI for common operations like issue management, sprint planning, and project navigation.

## Development Commands

### Building and Installing
```bash
# Install dependencies and build
make deps install

# Build only
make build

# Install with ldflags (includes version info)
make install
```

### Testing and Quality
```bash
# Run all CI checks (lint + test)
make ci

# Run linter only
make lint

# Run tests only (clears test cache first)
make test

# Run single test
go test -run TestSpecificTest ./pkg/jira

# Run tests with race detection
CGO_ENABLED=1 go test -race ./...
```

### Development Environment
```bash
# Start local Jira server for testing
make jira.server

# Clean build artifacts
make clean

# Full clean including caches
make distclean
```

## Architecture Overview

### Core Structure

- **cmd/jira/main.go**: Entry point that delegates to root command
- **internal/cmd/**: Command implementations organized by feature (issue, epic, sprint, etc.)
- **internal/cmd/root/**: Root command setup, config loading, and auth validation
- **pkg/jira/**: Core Jira API client and domain models
- **api/client.go**: High-level API wrapper with version proxying
- **internal/view/**: Display and formatting logic for different data types
- **pkg/tui/**: Terminal UI components built on rivo/tview

### Key Architectural Patterns

#### Command Pattern
Each major feature (issue, epic, sprint) has its own package under `internal/cmd/` with subcommands organized in subdirectories. Commands follow the Cobra CLI pattern.

#### API Version Proxying
The tool supports both Jira Cloud (v3 API) and Server/DC (v2 API). The `api` package provides proxy functions that route to the appropriate API version based on configuration:
- `ProxyCreate()`, `ProxyGetIssue()`, `ProxySearch()` etc.
- Installation type (`Cloud` vs `Local`) determines routing

#### Configuration Management
- Uses Viper for configuration management
- Config hierarchy: CLI flags → env vars → config file
- Default config location: `~/.jira/.config.yml`
- Supports `JIRA_CONFIG_FILE` env var for custom config paths

#### Authentication
- Supports multiple auth types: `basic`, `bearer` (PAT), and `mtls`
- Token sources: `JIRA_API_TOKEN` env var, `.netrc` file, or system keychain
- Auth validation occurs in root command's `PersistentPreRun`

#### View Layer
- Separation between data fetching (`pkg/jira`) and presentation (`internal/view`)
- Interactive TUI using tview library
- Plain text, CSV, and JSON output modes
- Markdown rendering for issue descriptions using Glamour

### Testing Approach

- HTTP client testing uses `httptest.NewServer()` for mocking Jira API
- Table-driven tests for complex scenarios
- Testdata files in JSON format for API responses
- Race detection enabled for concurrent operations
- Test files follow `*_test.go` naming convention

### Package Organization

- **internal/**: Application-specific code (commands, views, utilities)
- **pkg/**: Reusable packages (jira client, TUI components, parsers)
- **api/**: High-level API facade
- **cmd/**: Application entry points

## Key Configuration

### Environment Variables
- `JIRA_API_TOKEN`: Authentication token
- `JIRA_AUTH_TYPE`: Auth method (basic/bearer/mtls)
- `JIRA_CONFIG_FILE`: Custom config file path
- `CF_ACCESS_TOKEN`: Cloudflare Access token (optional)

### Project Structure
- Go modules with vendoring (`make deps` creates vendor/)
- Uses golangci-lint with specific rule set in `.golangci.yml`
- CI/CD via GitHub Actions with Go 1.24+
- Docker support with multi-stage builds

## Development Notes

### Adding New Commands
1. Create package under `internal/cmd/[feature]/`
2. Implement Cobra command structure
3. Add to root command in `internal/cmd/root/root.go`
4. Add corresponding API methods in `pkg/jira/` if needed
5. Add view formatting in `internal/view/` if needed

### API Client Extensions
- New API endpoints go in `pkg/jira/`
- Add version-specific implementations (v2/v3)
- Create proxy functions in `api/client.go`
- Follow existing patterns for error handling and HTTP client usage

### Testing New Features
- Mock HTTP responses with `httptest`
- Use `testdata/` directories for complex JSON fixtures
- Follow table-driven test patterns
- Include both positive and error cases

### Code Quality
- All code must pass golangci-lint rules
- Tests must pass with race detection
- Follow Go naming conventions
- Use structured logging where appropriate