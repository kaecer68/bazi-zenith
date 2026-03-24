package v1

import (
	"fmt"
	"strings"

	"github.com/kaecer68/bazi-zenith/internal/service"
	"github.com/kaecer68/bazi-zenith/pkg/basis"
	"github.com/kaecer68/bazi-zenith/pkg/engine"
)

// BaziResponse represents the standard JSON API response structure.
type BaziResponse struct {
	Gender              string                `json:"gender"`
	DayStem             string                `json:"day_stem"`
	Pillars             map[string]PillarData `json:"pillars"` // year, month, day, hour
	DaYun               []DaYunData           `json:"da_yun"`
	StartAgeY           int                   `json:"start_age_y"`
	StartAgeM           int                   `json:"start_age_m"`
	Strength            StrengthData          `json:"strength"`
	Advice              []InterpretationData  `json:"advice"`
	FavorableElements   []string              `json:"favorable_elements"`
	UnfavorableElements []string              `json:"unfavorable_elements"`
	Directions          BaziDirections        `json:"directions"`
	DetailChart         DetailChart           `json:"detail_chart"`

	YearPillar       string `json:"year_pillar,omitempty"`
	MonthPillar      string `json:"month_pillar,omitempty"`
	DayPillar        string `json:"day_pillar,omitempty"`
	HourPillar       string `json:"hour_pillar,omitempty"`
	DayMaster        string `json:"day_master,omitempty"`
	DayMasterElement string `json:"day_master_element,omitempty"`
}

type PillarData struct {
	Stem         string   `json:"stem"`
	Branch       string   `json:"branch"`
	TenGodStem   string   `json:"ten_god_stem"`
	HiddenStems  []string `json:"hidden_stems"`
	TenGodHidden []string `json:"ten_god_hidden"`
	NaYin        string   `json:"na_yin"`
	LifeStage    string   `json:"life_stage"`
	ShenSha      []string `json:"shen_sha"`
}

type DaYunData struct {
	Pillar    string `json:"pillar"`
	StartAge  int    `json:"start_age"`
	StartYear int    `json:"start_year,omitempty"`
}

// BaziDirections represents auspicious directions per life area.
type BaziDirections struct {
	Wealth       string `json:"wealth"`
	Career       string `json:"career"`
	Study        string `json:"study"`
	Relationship string `json:"relationship"`
}

type StrengthData struct {
	Score      float64 `json:"score"`
	Status     string  `json:"status"`
	IsDeLing   bool    `json:"is_de_ling"`
	IsDeDi     bool    `json:"is_de_di"`
	IsDeZhu    bool    `json:"is_de_zhu"`
	Percentage float64 `json:"percentage"`

	ScoreLegacy      float64 `json:"Score,omitempty"`
	StatusLegacy     string  `json:"Status,omitempty"`
	IsDeLingLegacy   bool    `json:"IsDeLing,omitempty"`
	IsDeDiLegacy     bool    `json:"IsDeDi,omitempty"`
	IsDeZhuLegacy    bool    `json:"IsDeZhu,omitempty"`
	PercentageLegacy float64 `json:"Percentage,omitempty"`
}

type InterpretationData struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type"`

	TitleLegacy   string `json:"Title,omitempty"`
	ContentLegacy string `json:"Content,omitempty"`
	TypeLegacy    string `json:"Type,omitempty"`
}

type DetailChart struct {
	Natal            NatalMatrix      `json:"natal"`
	DayunBoard       []BoardItem      `json:"dayun_board"`
	LiunianBoard     []BoardItem      `json:"liunian_board"`
	LiuyueBoard      []MonthBoardItem `json:"liuyue_board"`
	FiveElementState []string         `json:"five_element_state"`
	Prompts          DetailPrompts    `json:"prompts"`
}

type NatalMatrix struct {
	TenGodStem FourPillarsText `json:"ten_god_stem"`
	TianGan    FourPillarsText `json:"tian_gan"`
	DiZhi      FourPillarsText `json:"di_zhi"`
	CangGan    FourPillarsText `json:"cang_gan"`
	NaYin      FourPillarsText `json:"na_yin"`
	XingYun    FourPillarsText `json:"xing_yun"`
	ZiZuo      FourPillarsText `json:"zi_zuo"`
	KongWang   FourPillarsText `json:"kong_wang"`
}

type FourPillarsText struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
	Hour  string `json:"hour"`
}

type BoardItem struct {
	Index        int    `json:"index,omitempty"`
	Year         int    `json:"year,omitempty"`
	StartAge     int    `json:"start_age,omitempty"`
	StartYear    int    `json:"start_year,omitempty"`
	Pillar       string `json:"pillar"`
	TenGodStem   string `json:"ten_god_stem"`
	TenGodBranch string `json:"ten_god_branch,omitempty"`
}

type MonthBoardItem struct {
	Month        int    `json:"month"`
	Pillar       string `json:"pillar"`
	TenGodStem   string `json:"ten_god_stem"`
	TenGodBranch string `json:"ten_god_branch,omitempty"`
}

type DetailPrompts struct {
	Tiangan string `json:"tiangan"`
	Dizhi   string `json:"dizhi"`
}

// FromChart converts internal BaziChart to API Response.
func FromChart(c engine.BaziChart, advice []engine.Interpretation, targetYear int, birthYear int) BaziResponse {
	resp := BaziResponse{
		Gender:    string(c.Gender),
		DayStem:   string(c.DayStem),
		StartAgeY: c.StartAgeY,
		StartAgeM: c.StartAgeM,
		Strength:  buildStrength(c.Strength),
		Pillars:   make(map[string]PillarData),
	}

	resp.DayMaster = string(c.DayStem)
	resp.DayMasterElement = string(c.DayStem.Attr().Element)

	insights := service.BuildInsights(c)
	resp.FavorableElements = insights.FavorableElements
	resp.UnfavorableElements = insights.UnfavorableElements
	resp.Directions = BaziDirections{
		Wealth:       insights.Directions.Wealth,
		Career:       insights.Directions.Career,
		Study:        insights.Directions.Study,
		Relationship: insights.Directions.Relationship,
	}
	resp.Advice = buildAdvice(advice)

	mapPillar := func(name string, d engine.PillarDetail) {
		hStems := make([]string, len(d.HiddenStems))
		for i, h := range d.HiddenStems {
			hStems[i] = string(h.Stem)
		}
		tgHidden := make([]string, len(d.TenGodHidden))
		for i, t := range d.TenGodHidden {
			tgHidden[i] = string(t)
		}
		sSha := make([]string, len(d.ShenSha))
		for i, s := range d.ShenSha {
			sSha[i] = string(s)
		}

		resp.Pillars[name] = PillarData{
			Stem:         string(d.Pillar.Stem),
			Branch:       string(d.Pillar.Branch),
			TenGodStem:   string(d.TenGodStem),
			HiddenStems:  hStems,
			TenGodHidden: tgHidden,
			NaYin:        string(d.NaYin),
			LifeStage:    string(d.LifeStage),
			ShenSha:      sSha,
		}
	}

	mapPillar("year", c.YearPillar)
	mapPillar("month", c.MonthPillar)
	mapPillar("day", c.DayPillar)
	mapPillar("hour", c.HourPillar)

	resp.YearPillar = resp.Pillars["year"].Stem + resp.Pillars["year"].Branch
	resp.MonthPillar = resp.Pillars["month"].Stem + resp.Pillars["month"].Branch
	resp.DayPillar = resp.Pillars["day"].Stem + resp.Pillars["day"].Branch
	resp.HourPillar = resp.Pillars["hour"].Stem + resp.Pillars["hour"].Branch

	for _, dy := range c.DaYun {
		resp.DaYun = append(resp.DaYun, DaYunData{
			Pillar:    string(dy.Pillar.Stem) + string(dy.Pillar.Branch),
			StartAge:  dy.StartAge,
			StartYear: birthYear + dy.StartAge,
		})
	}

	resp.DetailChart = buildDetailChart(c, targetYear, birthYear)

	return resp
}

func buildStrength(s engine.StrengthAnalysis) StrengthData {
	return StrengthData{
		Score:            s.Score,
		Status:           s.Status,
		IsDeLing:         s.IsDeLing,
		IsDeDi:           s.IsDeDi,
		IsDeZhu:          s.IsDeZhu,
		Percentage:       s.Percentage,
		ScoreLegacy:      s.Score,
		StatusLegacy:     s.Status,
		IsDeLingLegacy:   s.IsDeLing,
		IsDeDiLegacy:     s.IsDeDi,
		IsDeZhuLegacy:    s.IsDeZhu,
		PercentageLegacy: s.Percentage,
	}
}

func buildAdvice(items []engine.Interpretation) []InterpretationData {
	resp := make([]InterpretationData, 0, len(items))
	for _, item := range items {
		resp = append(resp, InterpretationData{
			Title:         item.Title,
			Content:       item.Content,
			Type:          item.Type,
			TitleLegacy:   item.Title,
			ContentLegacy: item.Content,
			TypeLegacy:    item.Type,
		})
	}
	return resp
}

func buildDetailChart(c engine.BaziChart, targetYear int, birthYear int) DetailChart {
	natal := buildNatalMatrix(c)
	return DetailChart{
		Natal:            natal,
		DayunBoard:       buildDayunBoard(c, birthYear),
		LiunianBoard:     buildLiunianBoard(c, targetYear),
		LiuyueBoard:      buildLiuyueBoard(c, targetYear),
		FiveElementState: buildElementState(c),
		Prompts: DetailPrompts{
			Tiangan: buildTianganPrompt(c, targetYear, birthYear),
			Dizhi:   buildDizhiPrompt(c, targetYear, birthYear),
		},
	}
}

func buildNatalMatrix(c engine.BaziChart) NatalMatrix {
	year := c.YearPillar
	month := c.MonthPillar
	day := c.DayPillar
	hour := c.HourPillar

	toHidden := func(hs []basis.HiddenStem) string {
		parts := make([]string, 0, len(hs))
		for _, h := range hs {
			parts = append(parts, string(h.Stem))
		}
		return strings.Join(parts, "")
	}

	toSelfSeat := func(d engine.PillarDetail) string {
		if len(d.TenGodHidden) == 0 {
			return string(d.LifeStage)
		}
		return string(d.TenGodHidden[0])
	}

	toKongWang := func(d engine.PillarDetail) string {
		kw := basis.GetKongWang(d.Pillar)
		if len(kw) < 2 {
			return ""
		}
		return string(kw[0]) + string(kw[1])
	}

	return NatalMatrix{
		TenGodStem: FourPillarsText{Year: string(year.TenGodStem), Month: string(month.TenGodStem), Day: string(day.TenGodStem), Hour: string(hour.TenGodStem)},
		TianGan:    FourPillarsText{Year: string(year.Pillar.Stem), Month: string(month.Pillar.Stem), Day: string(day.Pillar.Stem), Hour: string(hour.Pillar.Stem)},
		DiZhi:      FourPillarsText{Year: string(year.Pillar.Branch), Month: string(month.Pillar.Branch), Day: string(day.Pillar.Branch), Hour: string(hour.Pillar.Branch)},
		CangGan:    FourPillarsText{Year: toHidden(year.HiddenStems), Month: toHidden(month.HiddenStems), Day: toHidden(day.HiddenStems), Hour: toHidden(hour.HiddenStems)},
		NaYin:      FourPillarsText{Year: string(year.NaYin), Month: string(month.NaYin), Day: string(day.NaYin), Hour: string(hour.NaYin)},
		XingYun:    FourPillarsText{Year: string(year.LifeStage), Month: string(month.LifeStage), Day: string(day.LifeStage), Hour: string(hour.LifeStage)},
		ZiZuo:      FourPillarsText{Year: toSelfSeat(year), Month: toSelfSeat(month), Day: toSelfSeat(day), Hour: toSelfSeat(hour)},
		KongWang:   FourPillarsText{Year: toKongWang(year), Month: toKongWang(month), Day: toKongWang(day), Hour: toKongWang(hour)},
	}
}

func buildDayunBoard(c engine.BaziChart, birthYear int) []BoardItem {
	items := make([]BoardItem, 0, len(c.DaYun))
	for idx, dy := range c.DaYun {
		items = append(items, BoardItem{
			Index:        idx + 1,
			StartAge:     dy.StartAge,
			StartYear:    birthYear + dy.StartAge,
			Pillar:       string(dy.Pillar.Stem) + string(dy.Pillar.Branch),
			TenGodStem:   string(basis.GetTenGod(c.DayStem, dy.Pillar.Stem)),
			TenGodBranch: tenGodFromBranch(c.DayStem, dy.Pillar.Branch),
		})
	}
	return items
}

func buildLiunianBoard(c engine.BaziChart, targetYear int) []BoardItem {
	items := make([]BoardItem, 0, 10)
	for i := 0; i < 10; i++ {
		y := targetYear + i
		p := basis.GetYearPillar(y)
		items = append(items, BoardItem{
			Year:         y,
			Pillar:       string(p.Stem) + string(p.Branch),
			TenGodStem:   string(basis.GetTenGod(c.DayStem, p.Stem)),
			TenGodBranch: tenGodFromBranch(c.DayStem, p.Branch),
		})
	}
	return items
}

func buildLiuyueBoard(c engine.BaziChart, targetYear int) []MonthBoardItem {
	items := make([]MonthBoardItem, 0, 12)
	yearPillar := basis.GetYearPillar(targetYear)
	for month := 1; month <= 12; month++ {
		p := basis.GetMonthPillar(yearPillar.Stem, month)
		items = append(items, MonthBoardItem{
			Month:        month,
			Pillar:       string(p.Stem) + string(p.Branch),
			TenGodStem:   string(basis.GetTenGod(c.DayStem, p.Stem)),
			TenGodBranch: tenGodFromBranch(c.DayStem, p.Branch),
		})
	}
	return items
}

func buildElementState(c engine.BaziChart) []string {
	monthBranch := c.MonthPillar.Pillar.Branch
	wang := seasonElementByBranch(monthBranch)
	ordered := []basis.Element{basis.Fire, basis.Earth, basis.Wood, basis.Water, basis.Metal}

	resp := make([]string, 0, len(ordered))
	for _, e := range ordered {
		state := "死"
		switch {
		case e == wang:
			state = "旺"
		case produces(wang, e):
			state = "相"
		case produces(e, wang):
			state = "休"
		case controls(wang, e):
			state = "囚"
		default:
			state = "死"
		}
		resp = append(resp, fmt.Sprintf("%s%s", string(e), state))
	}

	return resp
}

func buildTianganPrompt(c engine.BaziChart, targetYear int, birthYear int) string {
	liuNianStem := basis.GetYearPillar(targetYear).Stem
	items := []struct {
		Name string
		Stem basis.Stem
	}{
		{Name: "年干", Stem: c.YearPillar.Pillar.Stem},
		{Name: "月干", Stem: c.MonthPillar.Pillar.Stem},
		{Name: "日干", Stem: c.DayPillar.Pillar.Stem},
		{Name: "時干", Stem: c.HourPillar.Pillar.Stem},
	}

	if dayunPillar, ok := getActiveDayunPillar(c.DaYun, targetYear, birthYear); ok {
		items = append(items, struct {
			Name string
			Stem basis.Stem
		}{Name: "大運干", Stem: dayunPillar.Stem})
	}

	items = append(items, struct {
		Name string
		Stem basis.Stem
	}{
		Name: "流年干",
		Stem: liuNianStem,
	})

	prompts := make([]string, 0)
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			tags := make([]string, 0, 2)
			if element, ok := stemHeElement(items[i].Stem, items[j].Stem); ok {
				tags = append(tags, fmt.Sprintf("合%s", element))
			}
			if isStemChong(items[i].Stem, items[j].Stem) {
				tags = append(tags, "沖")
			}

			if len(tags) > 0 {
				prompts = append(prompts, fmt.Sprintf("%s%s與%s%s：%s", items[i].Name, items[i].Stem, items[j].Name, items[j].Stem, strings.Join(tags, "/")))
			}
		}
	}

	if len(prompts) == 0 {
		return "天干以平和為主。"
	}

	if len(prompts) > 8 {
		prompts = prompts[:8]
	}

	return strings.Join(prompts, "；") + "。"
}

func buildDizhiPrompt(c engine.BaziChart, targetYear int, birthYear int) string {
	liuNianBranch := basis.GetYearPillar(targetYear).Branch
	items := []struct {
		Name   string
		Branch basis.Branch
	}{
		{Name: "年支", Branch: c.YearPillar.Pillar.Branch},
		{Name: "月支", Branch: c.MonthPillar.Pillar.Branch},
		{Name: "日支", Branch: c.DayPillar.Pillar.Branch},
		{Name: "時支", Branch: c.HourPillar.Pillar.Branch},
	}

	if dayunPillar, ok := getActiveDayunPillar(c.DaYun, targetYear, birthYear); ok {
		items = append(items, struct {
			Name   string
			Branch basis.Branch
		}{Name: "大運支", Branch: dayunPillar.Branch})
	}

	items = append(items, struct {
		Name   string
		Branch basis.Branch
	}{
		Name:   "流年支",
		Branch: liuNianBranch,
	})

	prompts := make([]string, 0)
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			tags := make([]string, 0, 4)
			if basis.GetBranchHe(items[i].Branch, items[j].Branch) {
				tags = append(tags, "合")
			}
			if basis.GetBranchChong(items[i].Branch, items[j].Branch) {
				tags = append(tags, "沖")
			}
			if basis.GetBranchXing(items[i].Branch, items[j].Branch) {
				tags = append(tags, "刑")
			}
			if basis.GetBranchHai(items[i].Branch, items[j].Branch) {
				tags = append(tags, "害")
			}

			if len(tags) > 0 {
				prompts = append(prompts, fmt.Sprintf("%s%s與%s%s：%s", items[i].Name, items[i].Branch, items[j].Name, items[j].Branch, strings.Join(tags, "/")))
			}
		}
	}

	if len(prompts) == 0 {
		return "地支以平和為主。"
	}

	if len(prompts) > 8 {
		prompts = prompts[:8]
	}

	return strings.Join(prompts, "；") + "。"
}

func getActiveDayunPillar(dayun []basis.DaYunInfo, targetYear int, birthYear int) (basis.Pillar, bool) {
	if len(dayun) == 0 {
		return basis.Pillar{}, false
	}

	for i := 0; i < len(dayun); i++ {
		startYear := birthYear + dayun[i].StartAge
		nextStartYear := startYear + 10
		if i+1 < len(dayun) {
			nextStartYear = birthYear + dayun[i+1].StartAge
		}
		if targetYear >= startYear && targetYear < nextStartYear {
			return dayun[i].Pillar, true
		}
	}

	if targetYear < birthYear+dayun[0].StartAge {
		return dayun[0].Pillar, true
	}

	return dayun[len(dayun)-1].Pillar, true
}

func tenGodFromBranch(dayStem basis.Stem, branch basis.Branch) string {
	hidden := basis.GetHiddenStems(branch)
	if len(hidden) == 0 {
		return ""
	}
	return string(basis.GetTenGod(dayStem, hidden[0].Stem))
}

func seasonElementByBranch(b basis.Branch) basis.Element {
	switch b {
	case basis.YinB, basis.Mao, basis.Chen:
		return basis.Wood
	case basis.Si, basis.WuB, basis.Wei:
		return basis.Fire
	case basis.Shen, basis.You, basis.Xu:
		return basis.Metal
	default:
		return basis.Water
	}
}

func produces(from basis.Element, to basis.Element) bool {
	switch from {
	case basis.Wood:
		return to == basis.Fire
	case basis.Fire:
		return to == basis.Earth
	case basis.Earth:
		return to == basis.Metal
	case basis.Metal:
		return to == basis.Water
	case basis.Water:
		return to == basis.Wood
	}
	return false
}

func controls(from basis.Element, to basis.Element) bool {
	switch from {
	case basis.Wood:
		return to == basis.Earth
	case basis.Fire:
		return to == basis.Metal
	case basis.Earth:
		return to == basis.Water
	case basis.Metal:
		return to == basis.Wood
	case basis.Water:
		return to == basis.Fire
	}
	return false
}

func stemHeElement(a, b basis.Stem) (string, bool) {
	switch {
	case (a == basis.Jia && b == basis.Ji) || (a == basis.Ji && b == basis.Jia):
		return "土", true
	case (a == basis.Yi && b == basis.Geng) || (a == basis.Geng && b == basis.Yi):
		return "金", true
	case (a == basis.Bing && b == basis.Xin) || (a == basis.Xin && b == basis.Bing):
		return "水", true
	case (a == basis.Ding && b == basis.Ren) || (a == basis.Ren && b == basis.Ding):
		return "木", true
	case (a == basis.Wu && b == basis.Gui) || (a == basis.Gui && b == basis.Wu):
		return "火", true
	default:
		return "", false
	}
}

func isStemChong(a, b basis.Stem) bool {
	switch {
	case (a == basis.Jia && b == basis.Geng) || (a == basis.Geng && b == basis.Jia):
		return true
	case (a == basis.Yi && b == basis.Xin) || (a == basis.Xin && b == basis.Yi):
		return true
	case (a == basis.Bing && b == basis.Ren) || (a == basis.Ren && b == basis.Bing):
		return true
	case (a == basis.Ding && b == basis.Gui) || (a == basis.Gui && b == basis.Ding):
		return true
	default:
		return false
	}
}
