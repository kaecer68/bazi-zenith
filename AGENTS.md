# AGENTS.md - Bazi-Zenith Development Guide

**Generated:** 2026-03-18
**Module:** github.com/kaecer68/bazi-zenith
**Go Version:** 1.25.6

## OVERVIEW

八字命盤引擎 — Go library + CLI + gRPC/REST services for Chinese Bazi fortune-telling calculations. Built on lunar-zenith astronomical engine.

**Stack:** Go 1.25.6, gRPC, Protocol Buffers, Traditional Chinese domain constants

---

## 📋 契約優先開發流程 (Contract-First)

### 契約是唯一真相

**所有 API 變更必須遵循**:
```
destiny-contracts/openapi/bazi-zenith.yaml
```

**規則**:
- ✅ 先更新契約，再修改代碼
- ✅ 不允許添加契約未定義的欄位
- ✅ 不允許修改契約定義的類型
- ✅ 如發現契約有問題，先更新契約

### 契約文件位置

```
bazi-zenith/
├── contracts/              # ← symlink 指向 destiny-contracts
│   ├── openapi/
│   │   └── bazi-zenith.yaml   # ← 契約源
│   ├── TASK-BOARD.md       # 跨服務任務看板
│   └── HANDOFF.md          # AI 交接報告
```

### AI 任務執行流程

1. **檢查 TASK-BOARD.md** → 了解當前任務
2. **讀取契約文件** → 確認欄位定義
3. **生成代碼** → `make generate`
4. **實現業務邏輯** → `internal/service/`
5. **驗證** → `openapi-generator validate`
6. **填寫 HANDOFF.md** → 回報結果

### 完成檢查清單

```markdown
- [ ] 已讀取最新契約文件
- [ ] 已運行 make generate
- [ ] 已運行 openapi-generator validate
- [ ] 單元測試通過
- [ ] 新增欄位已出現在契約中
- [ ] API 響應範例與契約一致
- [ ] 已更新 HANDOFF.md
```

---

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Domain data (Stems, Branches, TenGods) | `pkg/basis/` | Core constants, algorithms |
| Bazi calculation engine | `pkg/engine/` | `BaziEngine`, `BaziChart` |
| REST API service | `cmd/bazi-server/`, `internal/grpcserver/` | HTTP server |
| gRPC service | `cmd/bazi-grpc/` | gRPC server |
| CLI tool | `cmd/bazi-cli/` | Terminal排盤 |
| Proto definitions | `api/proto/bazi/v1/` | Source of truth for RPC |
| Generated proto code | `gen/bazipb/` | **DO NOT EDIT** |

---

## ANTI-PATTERNS (THIS PROJECT)

- **NEVER edit** `gen/bazipb/*.go` — regenerated from proto
- **NEVER add API fields** without updating `contracts/openapi/bazi-zenith.yaml` first
- **AVOID custom error types** — use simple `error` returns (project convention)

---

## UNIQUE STYLES

- **Traditional Chinese** for domain constants (e.g., `JiaStem = "甲"`)
- **Mixed naming**: English for logic, Chinese for domain terms
- **Table-driven tests** in `*_test.go` files

---

## COMMANDS

```bash
# Generate proto code
make proto

# Build all binaries
make build-all

# Test
make test
# Or: go test ./...

# Format + vet + test
go fmt ./... && go vet ./... && go test ./...
```

---

## NOTES

- **Contract symlink**: `contracts/` → `destiny-contracts/` (external repo)
- **Dependencies**: `lunar-zenith` for astronomical calculations
- **Age**: Small project (~3k LOC), disciplined Go patterns
