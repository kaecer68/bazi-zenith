package engine

import (
	"github.com/kaecer68/bazi-zenith/pkg/basis"
)

// StrengthAnalysis contains the result of day stem strength evaluation.
type StrengthAnalysis struct {
	Score      float64 // Total score (usually out of 100 or relative)
	Status     string  // 身強, 身弱, 中和, 極強, 極弱
	IsDeLing   bool    // 得令 (Born in supportive month)
	IsDeDi     bool    // 得地 (Supportive branches)
	IsDeZhu    bool    // 得助 (Supportive stems)
	Percentage float64 // Support percentage
}

// AnalyzeStrength performs a basic quantitative analysis of the day stem's strength.
// This is a simplified model (Percentage-based) for MVP.
func (c *BaziChart) AnalyzeStrength() StrengthAnalysis {
	me := c.DayStem
	meAttr := me.Attr()

	totalSupport := 0.0

	// 1. Check Month Branch (De Ling - 得令) - Weight: 40%
	monthBranch := c.MonthPillar.Pillar.Branch
	if isSupportive(meAttr.Element, monthBranch) {
		totalSupport += 40.0
	}

	// 2. Check Other Branches (De Di - 得地) - Weight: 10% each
	branches := []basis.Branch{c.YearPillar.Pillar.Branch, c.DayPillar.Pillar.Branch, c.HourPillar.Pillar.Branch}
	for _, b := range branches {
		if isSupportive(meAttr.Element, b) {
			totalSupport += 10.0
		}
	}

	// 3. Check Other Stems (De Zhu - 得助) - Weight: 10% each
	stems := []basis.Stem{c.YearPillar.Pillar.Stem, c.MonthPillar.Pillar.Stem, c.HourPillar.Pillar.Stem}
	for _, s := range stems {
		if isStemSupportive(meAttr.Element, s) {
			totalSupport += 10.0
		}
	}

	analysis := StrengthAnalysis{
		Score:      totalSupport,
		Percentage: totalSupport,
		IsDeLing:   isSupportive(meAttr.Element, monthBranch),
	}

	// Status determination
	switch {
	case totalSupport >= 80:
		analysis.Status = "極強"
	case totalSupport >= 60:
		analysis.Status = "身強"
	case totalSupport >= 40:
		analysis.Status = "中和"
	case totalSupport >= 20:
		analysis.Status = "身弱"
	default:
		analysis.Status = "極弱"
	}

	return analysis
}

// isSupportive checks if a branch's main element supports or is the same as the stem element.
func isSupportive(meElement basis.Element, b basis.Branch) bool {
	// Simple mapping of branch to its primary element
	mainElement := map[basis.Branch]basis.Element{
		basis.Zi:   basis.Water,
		basis.Chou: basis.Earth,
		basis.YinB: basis.Wood,
		basis.Mao:  basis.Wood,
		basis.Chen: basis.Earth,
		basis.Si:   basis.Fire,
		basis.WuB:  basis.Fire,
		basis.Wei:  basis.Earth,
		basis.Shen: basis.Metal,
		basis.You:  basis.Metal,
		basis.Xu:   basis.Earth,
		basis.Hai:  basis.Water,
	}[b]

	return mainElement == meElement || isProducing(mainElement, meElement)
}

func isStemSupportive(meElement basis.Element, s basis.Stem) bool {
	sElement := s.Attr().Element
	return sElement == meElement || isProducing(sElement, meElement)
}

// Helper (replicated from basis to avoid circular or for internal engine use if needed,
// but since basis is accessible, we can use it or define local helpers)
func isProducing(from, to basis.Element) bool {
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
