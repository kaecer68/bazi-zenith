# cmd/ - Application Entry Points

**Purpose:** Three standalone binary entry points for CLI, REST API, and gRPC services.

## WHERE TO LOOK

| Binary | File | Purpose | Default Port |
|--------|------|---------|--------------|
| CLI tool | `bazi-cli/main.go` | Terminal-based chart generation | N/A |
| REST server | `bazi-server/main.go` | HTTP API + embedded gRPC | 8080 |
| gRPC server | `bazi-grpc/main.go` | Standalone gRPC service | 50051 |

## STRUCTURE

```
cmd/
├── bazi-cli/       # Terminal排盤工具
│   └── main.go     # Flag parsing, chart printing
├── bazi-server/    # REST API服務
│   └── main.go     # HTTP handlers, CORS, routes
└── bazi-grpc/      # gRPC服務
    └── main.go     # gRPC server, proto conversion
```

## CONVENTIONS

- **Timezone**: Hardcoded `Asia/Taipei` for all date parsing
- **Date format**: `2006-01-02 15:04` (Go reference time)
- **Gender enum**: `"male"` / `"female"` strings → `basis.Male` / `basis.Female`
- **Error handling**: HTTP 400 for bad input, 500 for internal errors
- **gRPC reflection**: Enabled for grpcurl/Postman exploration

## ANTI-PATTERNS

- **DO NOT** add business logic here — delegate to `pkg/engine`
- **DO NOT** parse dates manually — always use `time.ParseInLocation`
- **NEVER** mutate request data — pass to engine as-is

## BUILD COMMANDS

```bash
# Individual binaries
make build-cli     # → bin/bazi-cli
make build-rest    # → bin/bazi-server
make build-grpc    # → bin/bazi-grpc

# All binaries
make build-all

# Run with defaults
make run-rest      # Port 8080
make run-grpc      # Port 50051
```

## DEPENDENCIES

- `pkg/engine` — Chart calculation
- `pkg/basis` — Gender enum, timezone
- `pkg/api/v1` — JSON response types (REST only)
- `internal/grpcserver` — Shared gRPC implementation (bazi-server only)
- `gen/bazipb` — Proto-generated types (gRPC only)
