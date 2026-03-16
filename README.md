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

### REST API 服務

```bash
# 編譯 & 啟動
go build -o bazi-server ./cmd/bazi-server
./bazi-server -port 8080
```

**呼叫範例：**

```bash
curl -X POST http://localhost:8080/api/v1/chart \
  -H "Content-Type: application/json" \
  -d '{"datetime": "1990-05-15 14:30", "gender": "male", "target_year": 2025}'
```

| 端點 | 方法 | 說明 |
| :--- | :--- | :--- |
| `/api/v1/chart` | POST | 生成八字命盤與流年斷語 |
| `/health` | GET | 健康檢查 |

### gRPC 服務

```bash
# 編譯 & 啟動
go build -o bazi-grpc ./cmd/bazi-grpc
./bazi-grpc -port 50051
```

**Proto 定義**：`api/proto/bazi/v1/bazi.proto`

```protobuf
service BaziService {
  rpc GetChart (GetChartRequest) returns (GetChartResponse);
}
```

**grpcurl 呼叫範例：**

```bash
grpcurl -plaintext -d '{"datetime": "1990-05-15 14:30", "gender": "male", "target_year": 2025}' \
  localhost:50051 bazi.v1.BaziService/GetChart
```

> 已啟用 gRPC Server Reflection，可使用 grpcurl 或 Postman 直接探索 API。

### JSON API 對接（程式碼內嵌）

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
├── api/
│   └── proto/bazi/v1/     # gRPC Proto 定義
│       └── bazi.proto
├── cmd/
│   ├── bazi-cli/          # CLI 終端排盤工具
│   ├── bazi-server/       # REST API 服務 (HTTP)
│   └── bazi-grpc/         # gRPC 服務
├── gen/
│   └── bazipb/            # Proto 自動生成的 Go 代碼
├── pkg/
│   ├── api/v1/            # JSON API 資料交換模型
│   ├── basis/             # 基礎數據模型與算法
│   └── engine/            # 核心排盤引擎
├── Makefile               # 建構與 Proto 生成腳本
├── PRD.md                 # 產品需求文檔
├── LICENSE                # MIT License
└── go.mod
```

---

## 🔗 依賴

| 套件 | 用途 |
| :--- | :--- |
| [lunar-zenith](https://github.com/kaecer68/lunar-zenith) | 天文級高精度太陽黃經計算、節氣定位 |
| [google.golang.org/grpc](https://grpc.io/) | gRPC 框架 |
| [google.golang.org/protobuf](https://protobuf.dev/) | Protocol Buffers |

---

## 🛣️ 路線圖

- [x] Phase 1: 基礎數據模型（十神、藏干、納音、長生十二運）
- [x] Phase 2: 動態演算法（大運、流年流月、五虎遁）
- [x] Phase 3: 引擎整合（Lunar-Zenith 對接、神煞系統）
- [x] Phase 4: 斷語生成（身強身弱、沖合害、流年互動）
- [x] Phase 5: 應用介面（JSON API 模型、CLI 排盤工具）
- [x] Phase 6: 服務化（REST API、gRPC 服務、Proto 定義）
- [ ] Phase 7: 進階分析（格局判定、用神取用、深度斷語）

---

## � 關於作者

**德凱/KAECER** — 對傳統文化數位化有興趣的前端工程師

- **Blog**: https://goluck.im/
- **Twitter**: [@kaecer](https://twitter.com/kaecer)

### 🌟 相關專案

| 專案名稱 | 說明 |
| :--- | :--- |
| [lunar-zenith](https://github.com/kaecer68/lunar-zenith) | 黃曆大全 — 天文級高精度節氣與農曆計算引擎 |
| [ziwei-zenith](https://github.com/kaecer68/ziwei-zenith) | 紫微斗數 — 專業紫微斗數排盤與命理分析引擎 |
| [bazi-zenith](https://github.com/kaecer68/bazi-zenith) | 八字直斷 — 子平八字排盤與流年斷語引擎（本專案） |

---

## �📄 授權

本專案採用 [MIT License](LICENSE) 開源授權。

---

## 🙏 致謝

本專案為「德凱 GoLuck」系列的一部分，基於傳統子平八字命理學，結合現代天文算法與軟體工程實踐而成。
德凱/KAECER
[@kaecer](https://twitter.com/kaecer)
https://goluck.im/
**Made with ❤️ by [Kaecer68](https://github.com/kaecer68)**
