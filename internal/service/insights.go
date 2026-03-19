package service

import (
	"sort"

	"github.com/kaecer68/bazi-zenith/pkg/basis"
	"github.com/kaecer68/bazi-zenith/pkg/engine"
)

// ChartInsights 聚合命盤衍生資訊
type ChartInsights struct {
	FavorableElements   []string
	UnfavorableElements []string
	Directions          Directions
}

// Directions 代表四大領域吉方位
type Directions struct {
	Wealth       string
	Career       string
	Study        string
	Relationship string
}

var (
	produceMap = map[basis.Element]basis.Element{
		basis.Wood:  basis.Fire,
		basis.Fire:  basis.Earth,
		basis.Earth: basis.Metal,
		basis.Metal: basis.Water,
		basis.Water: basis.Wood,
	}
	controlMap = map[basis.Element]basis.Element{
		basis.Wood:  basis.Earth,
		basis.Fire:  basis.Metal,
		basis.Earth: basis.Water,
		basis.Metal: basis.Wood,
		basis.Water: basis.Fire,
	}
	elementDirections = map[basis.Element][]string{
		basis.Wood:  {"東", "東南"},
		basis.Fire:  {"南"},
		basis.Earth: {"西南", "東北"},
		basis.Metal: {"西", "西北"},
		basis.Water: {"北"},
	}
)

// BuildInsights 依據命盤強弱推導喜忌與吉方位
func BuildInsights(chart engine.BaziChart) ChartInsights {
	dayElement := chart.DayStem.Attr().Element
	fav, unfav := determineElements(dayElement, chart.Strength)

	dirs := determineDirections(dayElement)

	return ChartInsights{
		FavorableElements:   elementsToStrings(fav),
		UnfavorableElements: elementsToStrings(unfav),
		Directions:          dirs,
	}
}

func determineElements(dayElement basis.Element, strength engine.StrengthAnalysis) ([]basis.Element, []basis.Element) {
	resource := elementProducedBy(dayElement)
	output := elementProduces(dayElement)
	wealth := elementControls(dayElement)
	officer := elementControlledBy(dayElement)

	supporting := []basis.Element{dayElement, resource}
	reducing := []basis.Element{output, wealth, officer}

	switch strength.Status {
	case "身弱", "極弱":
		return uniqueElements(supporting), uniqueElements(reducing)
	case "身強", "極強":
		return uniqueElements(reducing), uniqueElements(supporting)
	default:
		neutralFav := []basis.Element{resource, output}
		neutralUnfav := []basis.Element{wealth}
		return uniqueElements(neutralFav), uniqueElements(neutralUnfav)
	}
}

func determineDirections(dayElement basis.Element) Directions {
	wealth := pickDirection(elementControls(dayElement), false)
	career := pickDirection(elementControlledBy(dayElement), false)
	study := pickDirection(elementProducedBy(dayElement), true)

	relationshipElement := dayElement
	relationship := pickDirection(relationshipElement, len(elementDirections[relationshipElement]) > 1)

	return Directions{
		Wealth:       wealth,
		Career:       career,
		Study:        study,
		Relationship: relationship,
	}
}

func elementProduces(e basis.Element) basis.Element {
	if val, ok := produceMap[e]; ok {
		return val
	}
	return basis.Wood
}

func elementProducedBy(e basis.Element) basis.Element {
	for k, v := range produceMap {
		if v == e {
			return k
		}
	}
	return basis.Water
}

func elementControls(e basis.Element) basis.Element {
	if val, ok := controlMap[e]; ok {
		return val
	}
	return basis.Earth
}

func elementControlledBy(e basis.Element) basis.Element {
	for k, v := range controlMap {
		if v == e {
			return k
		}
	}
	return basis.Metal
}

func pickDirection(el basis.Element, preferAlt bool) string {
	dirs := elementDirections[el]
	if len(dirs) == 0 {
		return "北"
	}
	if preferAlt && len(dirs) > 1 {
		return dirs[1]
	}
	return dirs[0]
}

func uniqueElements(elems []basis.Element) []basis.Element {
	seen := make(map[basis.Element]struct{})
	result := make([]basis.Element, 0, len(elems))
	for _, el := range elems {
		if el == "" {
			continue
		}
		if _, ok := seen[el]; ok {
			continue
		}
		seen[el] = struct{}{}
		result = append(result, el)
	}
	return result
}

func elementsToStrings(elems []basis.Element) []string {
	if len(elems) == 0 {
		return nil
	}
	out := make([]string, len(elems))
	for i, el := range elems {
		out[i] = string(el)
	}
	sort.Strings(out)
	return out
}
