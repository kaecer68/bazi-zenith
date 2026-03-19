# internal/ - Service Layer

**Purpose:** Internal service implementations and shared server infrastructure.

## WHERE TO LOOK

| Component | Location | Purpose |
|-----------|----------|---------|
| Business insights | `service/insights.go` | 喜忌分析, 吉方位推導 |
| gRPC server | `grpcserver/server.go` | Shared gRPC implementation |

## STRUCTURE

```
internal/
├── service/
│   └── insights.go      # ChartInsights, Directions, element analysis
└── grpcserver/
    └── server.go        # Server struct, GetChart RPC handler
```

## CONVENTIONS

- **Service pattern**: `BuildInsights(chart)` → returns derived analysis
- **Element mappings**: Use `produceMap`, `controlMap` for 五行生剋
- **Directions**: Map elements to 東/南/西/北/東南/西南/東北/西北
- **Strength-based logic**: 身弱→supporting, 身強→reducing

## ANTI-PATTERNS

- **DO NOT** add API routes here — that's in `cmd/bazi-server/`
- **DO NOT** duplicate domain logic from `pkg/basis`
- **NEVER** hardcode business rules as magic numbers — use `pkg/basis` constants

## KEY TYPES

```go
// service/insights.go
ChartInsights struct {
    FavorableElements   []string    // 喜用神
    UnfavorableElements []string    // 忌神
    Directions          Directions  // 吉方位
}

Directions struct {
    Wealth       string  // 財位
    Career       string  // 事業位
    Study        string  // 文昌位
    Relationship string  // 桃花位
}
```

## DEPENDENCIES

- `pkg/basis` — Elements, stems, branches
- `pkg/engine` — BaziChart, StrengthAnalysis
- `gen/bazipb` — Proto message types (grpcserver only)
