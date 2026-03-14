# PRD: Bazi-Zenith (八字命盤引擎)

## 0. 專案願景 (Project Vision)
作為 **GoLuck** 生態系的第二基石，`Bazi-Zenith` 致力於提供最標準、最易擴展的子平八字排盤數據結構。它將繼承 `Lunar-Zenith` 的天文級精度，並服務於所有需要八字數據的上層應用。

## 1. 核心定位 (Positioning)
- **輸入**: 基於 `Lunar-Zenith` 的四柱干支與精確時空數據。
- **輸出**: 完整的八字命盤數據模型（十神、藏干、納音、長生運、大運、流年）。
- **目標**: 成就 Go 語言界最專業的開源八字建模庫。

## 2. 功能範圍 (Scope)
### Phase 1: 基礎數據模型 (Base Model)
- [x] **十神換算**: 根據日元 (Day Stem) 推算其餘七字之十神（正官、偏印等）。
- [x] **地支藏干**: 實現標準的「地支本氣、中氣、餘氣」映射。
- [x] **納音五行**: 計算六十甲子對應的納音（如：海中金、爐中火）。
- [x] **長生十二運**: 推算各柱地支相對於日元的能量狀態。

### Phase 2: 動態運程 (Dynamic Pillars)
- [x] **大運推算**: 根據性別與出生年干陰陽，確定起運時間與大運序列。
- [x] **流年/流月**: 動態生成指定年份的運程數據。
 
### Phase 3: 引擎整合與分析 (Engine Integration)
- [x] **核心排盤模組**: 封裝 `Lunar-Zenith` 調用，自動生成四柱干支。
- [x] **完整命盤數據**: 整合十神、藏干、納音、長生運於一個數據模型。
- [x] **神煞系統**: 實現基本的神煞邏輯（天乙、桃花、驛馬等）。
- [x] **底層修正**: 升級 `lunar-zenith` 至 v0.1.1，修復年份硬編碼 Bug。

### Phase 4: 高級斷語與動態互動 (Interpretations)
- [x] **動態神煞互動**: 計算大運/流年相對於命盤的神煞觸發。
- [x] **干支生剋互動**: 判定大運/流年與命盤的「沖、合、害」。
- [x] **簡易斷語生成**: 基於身強身弱與流年干支生成基礎命理建議。

### Phase 5: 應用適配與交互佈署 (Application & Interface)
- [x] **API 資料模型**: 在 `pkg/api` 定義標準 JSON 交換格式。
- [x] **CLI 終端工具**: 實現 `cmd/bazi-cli` 支持命令行排盤與精美終端輸出。

## 3. 技術規格 (Technical Specs)
- **Language**: Go 1.25+
- **Input Engine**: `github.com/kaecer68/lunar-zenith` (v0.1.1+, fixed 2024 hardcode bug)
- **Output Format**: JSON / Struct (支持 gRPC)

## 4. 視覺與體驗 (Visual Standards)
- 所有輸出的中文符號必須採用「繁體中文」。
- 術語遵循台灣主流子平八字習慣。

## 5. 知識產權 (License)
- **MIT License** (與生態系保持一致)。
