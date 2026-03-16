# AGENTS.md - Bazi-Zenith Development Guide

## Project Overview
- **Project**: Bazi-Zenith (八字命盤引擎)
- **Type**: Go library and CLI tool for Chinese Bazi/Fortune-telling calculations
- **Go Version**: 1.25.6
- **Primary Language**: Traditional Chinese (繁體中文) for constants and documentation

## Build, Lint & Test Commands

### Running Tests
```bash
# Run all tests
go test ./...

# Run a single test file
go test ./pkg/basis/ -v

# Run a specific test
go test ./pkg/basis/ -run TestDayun -v
```

### Building
```bash
# Build CLI tool
go build -o bazi-cli ./cmd/bazi-cli

# Run CLI
./bazi-cli -dt "1990-05-15 14:30" -g male -y 2025
```

### Code Quality
```bash
# Format code
go fmt ./...

# Vet (static analysis)
go vet ./...

# Run all checks
go fmt ./... && go vet ./... && go test ./...
```

---

## Code Style Guidelines

### Naming Conventions
- **Types**: PascalCase (e.g., `BaziChart`, `PillarDetail`, `StemAttr`)
- **Constants**: PascalCase with Chinese characters (e.g., `Jia Stem = "甲"`)
- **Functions**: PascalCase (e.g., `GetTenGod`, `NewPillarDetail`)
- **Variables**: Mixed - use descriptive English names for logic, Chinese for domain terms
- **Packages**: Lowercase, single word (e.g., `basis`, `engine`, `v1`)

### Import Organization
Standard Go import grouping:
```go
import (
    "flag"
    "fmt"
    "os"
    "time"

    "github.com/kaecer68/bazi-zenith/pkg/basis"
    "github.com/kaecer68/bazi-zenith/pkg/engine"
)
```

### Type Definitions
- Use `type X string` for domain-specific string types
- Use structs for composite data (follow existing patterns in `model.go`)
- Embed basic types for clarity (e.g., `type Stem string`)

### Error Handling
- Use standard Go error patterns (return `error` from functions)
- No custom error types yet - prefer simple error checking
- Use `fmt.Errorf` for formatted error messages when needed

### JSON API Types
- Use `json:"field_name"` tags for serialization
- Export all API response fields (capitalized)
- Follow pattern in `pkg/api/v1/types.go` for conversion functions

### Testing
- Test files: `*_test.go` in same package as implementation
- Use table-driven tests for multiple test cases
- Test naming: `TestFunctionName`

### Code Comments
- Comment exported types and functions with English or Chinese
- Use doc comments (leading `//` on line before declaration)
- Example: `// Stem represents a Heavenly Stem (天干)`

### Project Structure
```
bazi-zenith/
├── cmd/bazi-cli/      # CLI entry point
│   └── main.go
├── pkg/
│   ├── api/v1/        # JSON API types
│   ├── basis/         # Core data definitions & algorithms
│   └── engine/        # Bazi calculation engine
└── go.mod
```

### Key Packages
- **pkg/basis**: Contains all domain data (Stems, Branches, Elements, TenGods, etc.)
- **pkg/engine**: Core `BaziEngine` and `BaziChart` types
- **pkg/api/v1**: JSON response conversion

### Common Operations
```go
// Create engine and generate chart
e := engine.NewBaziEngine()
chart := e.GetBaziChart(birthTime, gender)

// Get interpretations
advice := chart.GenerateInterpretations(2025)

// Convert to API response
resp := v1.FromChart(chart, advice)
```

---

## Dependencies
- `github.com/kaecer68/lunar-zenith` - Astronomical calendar calculations

## Notes
- This is a Chinese domain-specific project using Traditional Chinese characters
- Constants, error messages, and output may be in Chinese
- Follow existing patterns in the codebase for consistency