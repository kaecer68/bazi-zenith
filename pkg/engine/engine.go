package engine

import (
	"math"
	"time"

	"github.com/kaecer68/bazi-zenith/pkg/basis"
	"github.com/kaecer68/lunar-zenith/pkg/celestial"
)

// BaziEngine helps generate BaziChart from date/time.
type BaziEngine struct{}

// NewBaziEngine creates a new instance.
func NewBaziEngine() *BaziEngine {
	return &BaziEngine{}
}

// TimeToJD converts time.Time to Julian Day.
func TimeToJD(t time.Time) float64 {
	return float64(t.Unix())/86400.0 + 2440587.5
}

// GetBaziChart generates the complete BaziChart.
func (e *BaziEngine) GetBaziChart(birthTime time.Time, gender basis.Gender) BaziChart {
	jd := TimeToJD(birthTime)
	jde := jd + (69.0 / 86400.0) // Simple Delta-T approximation

	// 1. Calculate Year & Month Pillars based on Solar Terms
	// Bazi Year starts at "Li Chun" (立春, 315 deg)
	liChun := e.findNearestTerm(jde, 315.0)
	year := birthTime.Year()
	if jde < liChun {
		year--
	}
	yearPillar := basis.GetYearPillar(year)

	// Month starts at Section Terms (節氣)
	// Even indexes in the list below are the starts of 12 months.
	// Starts with Li Chun (立春) as month 1.

	monthIdx := 0
	lastTermJDE := 0.0
	nextTermJDE := 0.0

	// Find the section term immediately preceding jde
	// Since section terms are roughly 30 days apart, we can start from
	// the roughly estimated term and search backwards.
	currentLon := celestial.SolarLongitude(jde)
	// Normalization
	normLon := math.Mod(currentLon-315.0+360.0, 360.0)
	monthIdx = int(math.Floor(normLon/30.0)) + 1 // 1=Tiger, 2=Rabbit...

	lastTermLon := math.Mod(float64(monthIdx-1)*30.0+315.0, 360.0)
	nextTermLon := math.Mod(float64(monthIdx)*30.0+315.0, 360.0)

	lastTermJDE = e.findNearestTerm(jde, lastTermLon)
	// If the found term is actually in the future, go back 30 days
	if lastTermJDE > jde {
		lastTermJDE = e.findNearestTerm(jde-30, lastTermLon)
	}

	nextTermJDE = e.findNearestTerm(lastTermJDE+30, nextTermLon)
	// Ensure nextTermJDE is after jde
	if nextTermJDE <= jde {
		nextTermJDE = e.findNearestTerm(jde+30, nextTermLon)
	}

	monthPillar := basis.GetMonthPillar(yearPillar.Stem, monthIdx)

	// 2. Day Pillar
	// 2000-01-01 was Jia-Wu (index 30 of Jiazi)
	// JD of 2000-01-01 12:00 UTC is 2451545.0
	// For Local Time (UTC+8), 00:00 is 2451544.166...
	// Actually Jiazi changes at Midnight (00:00).
	// Formula: (JD + 0.5 + 8/24) offset from Jiazi.
	dayIdx := int(math.Floor(jd+0.5+8.0/24.0+49)) % 60
	dayPillar := basis.JiaziList[dayIdx]

	// 3. Hour Pillar
	// (DayStemIndex * 2 + HourIndex) % 10 -> Hour Stem
	// HourIndex (0-11): Zi (23-1), Chou (1-3)...
	hourIdx := (birthTime.Hour() + 1) / 2 % 12
	stemOrder := []basis.Stem{basis.Jia, basis.Yi, basis.Bing, basis.Ding, basis.Wu, basis.Ji, basis.Geng, basis.Xin, basis.Ren, basis.Gui}
	dayStemIdx := -1
	for i, s := range stemOrder {
		if s == dayPillar.Stem {
			dayStemIdx = i
			break
		}
	}
	hourStemIdx := (dayStemIdx*2 + hourIdx) % 10
	hourPillar := basis.Pillar{
		Stem:   stemOrder[hourStemIdx],
		Branch: basis.Zi, // Temporary, will be corrected below
	}
	// Correct Hour Pillar Branch
	hourBranches := []basis.Branch{basis.Zi, basis.Chou, basis.YinB, basis.Mao, basis.Chen, basis.Si, basis.WuB, basis.Wei, basis.Shen, basis.You, basis.Xu, basis.Hai}
	hourPillar.Branch = hourBranches[hourIdx]

	// 4. Da Yun Sequence & Age
	daYunSeq := basis.GetDaYunSequence(yearPillar.Stem, monthPillar, gender)

	// Calculate Age: birthTime to nearest term
	// If clockwise, to NEXT term. If counter, to PREVIOUS term.
	yearPolarity := yearPillar.Stem.Attr().Polarity
	isClockwise := (gender == basis.Male && yearPolarity == basis.Yang) || (gender == basis.Female && yearPolarity == basis.Yin)

	var diffSeconds int64
	if isClockwise {
		diffSeconds = int64((nextTermJDE - jde) * 86400)
	} else {
		diffSeconds = int64((jde - lastTermJDE) * 86400)
	}

	dyYears, dyMonths := basis.CalculateDaYunAge(diffSeconds)

	daYunList := make([]basis.DaYunInfo, len(daYunSeq))
	for i, p := range daYunSeq {
		daYunList[i] = basis.DaYunInfo{
			Pillar:   p,
			StartAge: dyYears + (i * 10), // Simplification: exactly 10 years per stage
		}
	}

	chart := BaziChart{
		Gender:      gender,
		DayStem:     dayPillar.Stem,
		YearPillar:  NewPillarDetail(yearPillar, dayPillar.Stem),
		MonthPillar: NewPillarDetail(monthPillar, dayPillar.Stem),
		DayPillar:   NewPillarDetail(dayPillar, dayPillar.Stem),
		HourPillar:  NewPillarDetail(hourPillar, dayPillar.Stem),
		DaYun:       daYunList,
		StartAgeY:   dyYears,
		StartAgeM:   dyMonths,
	}

	chart.PopulateShenSha()
	chart.Strength = chart.AnalyzeStrength()
	return chart
}

// PopulateShenSha scans the chart for symbolic stars.
func (c *BaziChart) PopulateShenSha() {
	pillars := []*PillarDetail{&c.YearPillar, &c.MonthPillar, &c.DayPillar, &c.HourPillar}

	for _, p := range pillars {
		// 1. Based on Day Stem
		if basis.GetTianYi(c.DayStem, p.Pillar.Branch) {
			p.ShenSha = append(p.ShenSha, basis.TianYi)
		}
		if basis.GetLuShen(c.DayStem, p.Pillar.Branch) {
			p.ShenSha = append(p.ShenSha, basis.LuShen)
		}
		if basis.GetYangRen(c.DayStem, p.Pillar.Branch) {
			p.ShenSha = append(p.ShenSha, basis.YangRen)
		}
		if basis.GetWenChang(c.DayStem, p.Pillar.Branch) {
			p.ShenSha = append(p.ShenSha, basis.WenChang)
		}

		// 2. Based on Year Branch
		if basis.GetYiMa(c.YearPillar.Pillar.Branch, p.Pillar.Branch) {
			p.ShenSha = append(p.ShenSha, basis.YiMa)
		}
		if basis.GetTaoHua(c.YearPillar.Pillar.Branch, p.Pillar.Branch) {
			p.ShenSha = append(p.ShenSha, basis.TaoHua)
		}
		if basis.GetHuaGai(c.YearPillar.Pillar.Branch, p.Pillar.Branch) {
			p.ShenSha = append(p.ShenSha, basis.HuaGai)
		}
		if basis.GetJiangXing(c.YearPillar.Pillar.Branch, p.Pillar.Branch) {
			p.ShenSha = append(p.ShenSha, basis.JiangXing)
		}
		if basis.GetHongLuan(c.YearPillar.Pillar.Branch, p.Pillar.Branch) {
			p.ShenSha = append(p.ShenSha, basis.HongLuan)
		}
		if basis.GetTianXi(c.YearPillar.Pillar.Branch, p.Pillar.Branch) {
			p.ShenSha = append(p.ShenSha, basis.TianXi)
		}

		// 3. Based on Day Branch (Cross-check)
		if p != &c.DayPillar {
			if basis.GetYiMa(c.DayPillar.Pillar.Branch, p.Pillar.Branch) {
				p.ShenSha = append(p.ShenSha, basis.YiMa)
			}
			if basis.GetTaoHua(c.DayPillar.Pillar.Branch, p.Pillar.Branch) {
				p.ShenSha = append(p.ShenSha, basis.TaoHua)
			}
		}
	}
}

func (e *BaziEngine) findNearestTerm(jde float64, targetLon float64) float64 {
	// Search within +/- 20 days
	return celestial.EstimateTermTime(targetLon, jde-20, jde+20)
}
