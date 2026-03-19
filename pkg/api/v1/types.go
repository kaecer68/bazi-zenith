package v1

import (
	"github.com/kaecer68/bazi-zenith/internal/service"
	"github.com/kaecer68/bazi-zenith/pkg/engine"
)

// BaziResponse represents the standard JSON API response structure.
type BaziResponse struct {
	Gender              string                  `json:"gender"`
	DayStem             string                  `json:"day_stem"`
	Pillars             map[string]PillarData   `json:"pillars"` // year, month, day, hour
	DaYun               []DaYunData             `json:"da_yun"`
	StartAgeY           int                     `json:"start_age_y"`
	StartAgeM           int                     `json:"start_age_m"`
	Strength            engine.StrengthAnalysis `json:"strength"`
	Advice              []engine.Interpretation `json:"advice"`
	FavorableElements   []string                `json:"favorable_elements"`
	UnfavorableElements []string                `json:"unfavorable_elements"`
	Directions          BaziDirections          `json:"directions"`
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
	Pillar   string `json:"pillar"`
	StartAge int    `json:"start_age"`
}

// BaziDirections represents auspicious directions per life area.
type BaziDirections struct {
	Wealth       string `json:"wealth"`
	Career       string `json:"career"`
	Study        string `json:"study"`
	Relationship string `json:"relationship"`
}

// FromChart converts internal BaziChart to API Response.
func FromChart(c engine.BaziChart, advice []engine.Interpretation) BaziResponse {
	resp := BaziResponse{
		Gender:    string(c.Gender),
		DayStem:   string(c.DayStem),
		StartAgeY: c.StartAgeY,
		StartAgeM: c.StartAgeM,
		Strength:  c.Strength,
		Advice:    advice,
		Pillars:   make(map[string]PillarData),
	}

	insights := service.BuildInsights(c)
	resp.FavorableElements = insights.FavorableElements
	resp.UnfavorableElements = insights.UnfavorableElements
	resp.Directions = BaziDirections{
		Wealth:       insights.Directions.Wealth,
		Career:       insights.Directions.Career,
		Study:        insights.Directions.Study,
		Relationship: insights.Directions.Relationship,
	}

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

	for _, dy := range c.DaYun {
		resp.DaYun = append(resp.DaYun, DaYunData{
			Pillar:   string(dy.Pillar.Stem) + string(dy.Pillar.Branch),
			StartAge: dy.StartAge,
		})
	}

	return resp
}
