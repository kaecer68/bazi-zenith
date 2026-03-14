package engine

import (
	"fmt"
	"github.com/kaecer68/bazi-zenith/pkg/basis"
)

// Interpretation represents a single piece of astrological advice.
type Interpretation struct {
	Title   string
	Content string
	Type    string // 平, 吉, 凶
}

// GenerateInterpretations generates a list of advices based on the chart and a target year.
func (c *BaziChart) GenerateInterpretations(targetYear int) []Interpretation {
	results := []Interpretation{}

	yearPillar := basis.GetYearPillar(targetYear)

	// 1. General Strength Advice
	results = append(results, c.getStrengthAdvice(yearPillar)...)

	// 2. Interaction Advice (Fan Tai Sui, Clashes)
	results = append(results, c.getInteractionAdvice(yearPillar)...)

	return results
}

func (c *BaziChart) getStrengthAdvice(yearPillar basis.Pillar) []Interpretation {
	me := c.DayStem
	meAttr := me.Attr()
	yearStemAttr := yearPillar.Stem.Attr()

	adv := Interpretation{
		Title: "流年與身強關係",
		Type:  "平",
	}

	// Simplified logic: If weak, like supportive elements. If strong, like controlling/draining elements.
	isSupportingYear := (yearStemAttr.Element == meAttr.Element) || isProducing(yearStemAttr.Element, meAttr.Element)

	if c.Strength.Status == "身弱" || c.Strength.Status == "極弱" {
		if isSupportingYear {
			adv.Content = "今年流年對您有生扶之功，運勢較佳，適合穩步發展。"
			adv.Type = "吉"
		} else {
			adv.Content = "今年運勢較平，日元能量較低，宜保守行事，多注意休息。"
		}
	} else if c.Strength.Status == "身強" || c.Strength.Status == "極強" {
		if isSupportingYear {
			adv.Content = "今年流年助力過強，可能會有壓力或閉塞感，注意身心調適。"
			adv.Type = "平"
		} else {
			adv.Content = "今年流年有洩秀/制約之功，才華得以施展，財運、事業運佳。"
			adv.Type = "吉"
		}
	} else {
		adv.Content = "今年運勢平衡穩定，適合落實既定計畫。"
	}

	return []Interpretation{adv}
}

func (c *BaziChart) getInteractionAdvice(yearPillar basis.Pillar) []Interpretation {
	advs := []Interpretation{}

	// Bazi Year vs Liu Nian Branch (Fan Tai Sui)
	if yearPillar.Branch == c.YearPillar.Pillar.Branch {
		advs = append(advs, Interpretation{
			Title:   "值太歲",
			Content: "今年為您的本命年（值太歲），正所謂「太歲當頭坐，無喜必有禍」，心態宜平和。",
			Type:    "平",
		})
	}

	// Clashes with Pillars
	pillars := []struct {
		Name   string
		Pillar basis.Pillar
	}{
		{"年柱", c.YearPillar.Pillar},
		{"月柱", c.MonthPillar.Pillar},
		{"日柱", c.DayPillar.Pillar},
		{"時柱", c.HourPillar.Pillar},
	}

	for _, p := range pillars {
		if basis.GetBranchChong(yearPillar.Branch, p.Pillar.Branch) {
			advs = append(advs, Interpretation{
				Title:   fmt.Sprintf("沖%s", p.Name),
				Content: fmt.Sprintf("流年與您的%s地支發生相沖（六沖），%s波動較大，需注意變動。", p.Name, p.Name),
				Type:    "凶",
			})
		}
	}

	// Dynamic ShenSha: Is current year branch a star for the natal chart?
	if basis.GetTianYi(c.DayStem, yearPillar.Branch) {
		advs = append(advs, Interpretation{
			Title:   "天乙貴人入流年",
			Content: "今年流年支為您的天乙貴人，主有貴人相助，逢凶化吉。",
			Type:    "吉",
		})
	}
	if basis.GetTaoHua(c.YearPillar.Pillar.Branch, yearPillar.Branch) || basis.GetTaoHua(c.DayPillar.Pillar.Branch, yearPillar.Branch) {
		advs = append(advs, Interpretation{
			Title:   "桃花入流年",
			Content: "今年流年支為您的桃花星，主異性緣佳，社交生活豐富。",
			Type:    "吉",
		})
	}
	if basis.GetYiMa(c.YearPillar.Pillar.Branch, yearPillar.Branch) || basis.GetYiMa(c.DayPillar.Pillar.Branch, yearPillar.Branch) {
		advs = append(advs, Interpretation{
			Title:   "驛馬入流年",
			Content: "今年流年支為您的驛馬星，主遠行、搬遷或事業變動頻繁。",
			Type:    "平",
		})
	}

	return advs
}
