# pkg/engine - Bazi Calculation Engine

**Purpose:** Core `BaziEngine` and `BaziChart` types for chart generation and interpretation.

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| Chart calculation | `engine.go` | `BaziEngine.GetBaziChart()`, pillar computation |
| Strength analysis | `strength.go` | Day master strength scoring (身強/身弱) |
| Interpretations | `interpretation.go` | Yearly forecasts, advice generation |
| Data models | `model.go` | `BaziChart`, `PillarDetail`, `StrengthResult` |

## CONVENTIONS

- **Engine pattern**: `NewBaziEngine()` → `GetBaziChart(time, gender)` → `GenerateInterpretations(year)`
- **Return types**: Always return complete structs (`BaziChart`), not partial data
- **Time handling**: Accept `time.Time` with timezone — caller responsible for correct TZ

## ANTI-PATTERNS

- **DO NOT** duplicate domain logic from `pkg/basis` — import and use basis functions
- **DO NOT** add API-specific types here — use `pkg/api/v1` for JSON/gRPC conversion

## DEPENDENCIES

- `pkg/basis` — Domain constants and algorithms
- `lunar-zenith` — Astronomical calculations (solar terms, jieqi)

## TESTING

- Integration tests validate full chart generation
- Table-driven tests for edge cases (midnight births, leap months, etc.)
