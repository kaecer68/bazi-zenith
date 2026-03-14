package engine

import (
	"github.com/kaecer68/bazi-zenith/pkg/basis"
)

// PillarDetail contains all calculated attributes for a single pillar.
type PillarDetail struct {
	Pillar       basis.Pillar
	TenGodStem   basis.TenGod
	HiddenStems  []basis.HiddenStem
	TenGodHidden []basis.TenGod
	NaYin        basis.NaYin
	LifeStage    basis.LifeStage
	ShenSha      []basis.ShenSha
}

// BaziChart represents a complete Bazi profile.
type BaziChart struct {
	Gender      basis.Gender
	YearPillar  PillarDetail
	MonthPillar PillarDetail
	DayPillar   PillarDetail
	HourPillar  PillarDetail

	DayStem basis.Stem

	DaYun     []basis.DaYunInfo
	StartAgeY int
	StartAgeM int

	Strength StrengthAnalysis
}

// NewPillarDetail creates a detail object from a pillar and a day stem.
func NewPillarDetail(p basis.Pillar, dayStem basis.Stem) PillarDetail {
	hidden := basis.GetHiddenStems(p.Branch)
	tgHidden := make([]basis.TenGod, len(hidden))
	for i, h := range hidden {
		tgHidden[i] = basis.GetTenGod(dayStem, h.Stem)
	}

	return PillarDetail{
		Pillar:       p,
		TenGodStem:   basis.GetTenGod(dayStem, p.Stem),
		HiddenStems:  hidden,
		TenGodHidden: tgHidden,
		NaYin:        basis.GetNaYin(p.Stem, p.Branch),
		LifeStage:    basis.GetLifeStage(dayStem, p.Branch),
		ShenSha:      []basis.ShenSha{}, // Calculated during full chart generation
	}
}
