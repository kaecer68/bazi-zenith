# pkg/basis - Domain Data & Algorithms

**Purpose:** Core Bazi domain constants, data models, and calculation algorithms.

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| Heavenly Stems (天干) | `cycles.go` | 10 stems, yin/yang, five elements |
| Earthly Branches (地支) | `cycles.go` | 12 branches, zodiac animals |
| Hidden Stems (藏干) | `hidden_stems.go` | Benqi, Zhongqi, Yuqi mappings |
| Ten Gods (十神) | `ten_gods.go` | Stem-to-stem relationships |
| NaYin (納音) | `na_yin.go` | 60 Jiazi sound elements |
| Life Cycle (長生十二運) | `life_cycle.go` | 12 life stages |
| ShenSha (神煞) | `shensha.go` | Lucky/unlucky stars |
| Interactions | `interactions.go` | Clashes, harms, combinations |

## CONVENTIONS

- **Traditional Chinese constants**: `JiaStem = "甲"`, `ZiBranch = "子"`
- **Type aliases**: `type Stem string`, `type Branch string`
- **Algorithm functions**: Pure functions, no state (e.g., `GetHiddenStems(branch Branch) []Stem`)

## ANTI-PATTERNS

- **DO NOT** mix domain data with business logic — keep calculations here, interpretation in `pkg/engine`
- **DO NOT** use English for domain constants — Traditional Chinese only

## TESTING

- Table-driven tests in `*_test.go` files
- Test naming: `TestFunctionName` (e.g., `TestGetHiddenStems`)
