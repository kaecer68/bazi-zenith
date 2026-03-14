# 🔱 Bazi-Zenith (八字命盤引擎)

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Lunar-Zenith](https://img.shields.io/badge/Engine-Lunar--Zenith-blueviolet?style=flat-square)](https://github.com/kaecer68/lunar-zenith)

**高精度子平八字排盤引擎** — 基於天文級黃曆引擎 [Lunar-Zenith](https://github.com/kaecer68/lunar-zenith) 構建，提供專業的四柱八字計算、命理分析與流年斷語生成。

> 🌏 本專案以 **繁體中文** 為主要語言，遵循「德凱 GoLuck」系列標準。

---

## ✨ 功能特色

### 🧮 精確排盤
- **四柱八字**：自動計算年、月、日、時柱（以立春為歲首、節氣定月令）
- **地支藏干**：完整的本氣、中氣、餘氣映射
- **十神關係**：天干十神 + 藏干十神全面解析
- **納音五行**：六十甲子納音對照
- **長生十二運**：日元在各地支的生旺死絕狀態

### 🔮 命理分析
- **身強身弱**：基於得令、得地、得助的量化評分模型
- **大運排序**：根據性別與年干陰陽自動推算順逆大運
- **起運歲數**：精確計算至月（基於出生時刻與最近節氣的時間差）

### ⭐ 神煞系統
| 類別 | 星煞 |
| :--- | :--- |
| **日干系** | 天乙貴人、祿神、羊刃、文昌貴人 |
| **支對支系** | 驛馬、桃花、華蓋、將星 |
| **歲氣系** | 紅鸞、天喜 |

### 📜 流年斷語
- **值太歲 / 沖太歲**：自動檢測流年與命盤的沖合關係
- **動態神煞互動**：流年支觸發貴人、桃花、驛馬等
- **喜忌分析**：根據身強身弱與流年五行生剋，生成運勢建議

---

## 🚀 快速開始

### 安裝

```bash
go get github.com/kaecer68/bazi-zenith
```

### 作為 Go 庫使用

```go
package main

import (
    "fmt"
    "time"

    "github.com/kaecer68/bazi-zenith/pkg/basis"
    "github.com/kaecer68/bazi-zenith/pkg/engine"
)

func main() {
    e := engine.NewBaziEngine()
    loc, _ := time.LoadLocation("Asia/Taipei")
    birth := time.Date(1990, 5, 15, 14, 30, 0, 0, loc)

    chart := e.GetBaziChart(birth, basis.Male)

    fmt.Printf("日元: %s\n", chart.DayStem)
    fmt.Printf("年柱: %s%s (%s)\n",
        chart.YearPillar.Pillar.Stem,
        chart.YearPillar.Pillar.Branch,
        chart.YearPillar.NaYin)
    fmt.Printf("身強分析: %s (%.1f)\n",
        chart.Strength.Status, chart.Strength.Score)

    // 流年分析
    advice := chart.GenerateInterpretations(2025)
    for _, a := range advice {
        fmt.Printf("[%s] %s: %s\n", a.Type, a.Title, a.Content)
    }
}
```

### CLI 工具

```bash
# 編譯
go build -o bazi-cli ./cmd/bazi-cli

# 排盤 (公曆日期 + 時間)
./bazi-cli -dt "1990-05-15 14:30" -g male -y 2025
```

**輸出範例：**

```text
==========================================
   Bazi-Zenith (八字命盤引擎) - 乾造
==========================================
      【年柱】  【月柱】  【日柱】  【時柱】
十神:  正印        比肩        日元        正財
天干:    甲       丁       丁       庚
地支:    辰       卯       丑       戌
藏干:  戊乙癸       乙         己癸辛       戊辛丁
納音:  佛燈火       爐中火       澗下水       釵釧金
------------------------------------------
身強分析: 身強 (總分: 60.0) | 起運歲數: 6 歲 11 個月
------------------------------------------
大運: [6]戊辰 [16]己巳 [26]庚午 [36]辛未 [46]壬申 [56]癸酉
------------------------------------------
★ 2025 (乙巳) 流年斷語:
 ● [平] 流年與身強關係: ...
==========================================
```

### JSON API 對接

```go
import v1 "github.com/kaecer68/bazi-zenith/pkg/api/v1"

// 將 BaziChart 轉換為標準 JSON 結構
response := v1.FromChart(chart, advice)
// 可直接 json.Marshal(response) 回傳給前端
```

---

## 📁 專案結構

```
bazi-zenith/
├── cmd/
│   └── bazi-cli/          # CLI 終端排盤工具
│       └── main.go
├── pkg/
│   ├── api/
│   │   └── v1/            # JSON API 資料交換模型
│   │       └── types.go
│   ├── basis/             # 基礎數據模型與算法
│   │   ├── definitions.go # 天干、地支、五行、陰陽定義
│   │   ├── hidden_stems.go# 地支藏干
│   │   ├── ten_gods.go    # 十神關係
│   │   ├── na_yin.go      # 納音五行
│   │   ├── life_cycle.go  # 長生十二運
│   │   ├── pillars.go     # 六十甲子序列
│   │   ├── dayun.go       # 大運算法
│   │   ├── cycles.go      # 流年流月算法
│   │   ├── shensha.go     # 神煞系統
│   │   └── interactions.go# 沖合害
│   └── engine/            # 核心排盤引擎
│       ├── engine.go      # BaziEngine 主邏輯
│       ├── model.go       # BaziChart 數據模型
│       ├── strength.go    # 身強身弱分析
│       └── interpretation.go # 斷語生成
├── PRD.md                 # 產品需求文檔
├── LICENSE                # MIT License
└── go.mod
```

---

## 🔗 依賴

| 套件 | 用途 |
| :--- | :--- |
| [lunar-zenith](https://github.com/kaecer68/lunar-zenith) | 天文級高精度太陽黃經計算、節氣定位 |

---

## 🛣️ 路線圖

- [x] Phase 1: 基礎數據模型（十神、藏干、納音、長生十二運）
- [x] Phase 2: 動態演算法（大運、流年流月、五虎遁）
- [x] Phase 3: 引擎整合（Lunar-Zenith 對接、神煞系統）
- [x] Phase 4: 斷語生成（身強身弱、沖合害、流年互動）
- [x] Phase 5: 應用介面（JSON API 模型、CLI 排盤工具）
- [ ] Phase 6: 進階分析（格局判定、用神取用、深度斷語）

---

## 📄 授權

本專案採用 [MIT License](LICENSE) 開源授權。

---

## 🙏 致謝

本專案為「德凱 GoLuck」系列的一部分，基於傳統子平八字命理學，結合現代天文算法與軟體工程實踐而成。

**Made with ❤️ by [Kaecer68](https://github.com/kaecer68)**
