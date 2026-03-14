package basis

// TenGod represents the relationship between the Day Stem and other Stems.
type TenGod string

const (
	BiJian    TenGod = "比肩"
	JieCai    TenGod = "劫財"
	ShiShen   TenGod = "食神"
	ShangGuan TenGod = "傷官"
	PianCai   TenGod = "偏財"
	ZhengCai  TenGod = "正財"
	QiSha     TenGod = "七殺"
	ZhengGuan TenGod = "正官"
	PianYin   TenGod = "偏印"
	ZhengYin  TenGod = "正印"
)

// GetTenGod calculates the Ten God relationship relative to the Day Stem.
func GetTenGod(me Stem, other Stem) TenGod {
	meAttr := me.Attr()
	otherAttr := other.Attr()

	samePolarity := (meAttr.Polarity == otherAttr.Polarity)

	// Same Element
	if meAttr.Element == otherAttr.Element {
		if samePolarity {
			return BiJian
		}
		return JieCai
	}

	// Me Produces Other (Output)
	if isProducing(meAttr.Element, otherAttr.Element) {
		if samePolarity {
			return ShiShen
		}
		return ShangGuan
	}

	// Other Produces Me (Input)
	if isProducing(otherAttr.Element, meAttr.Element) {
		if samePolarity {
			return PianYin
		}
		return ZhengYin
	}

	// Me Controls Other (Wealth)
	if isControlling(meAttr.Element, otherAttr.Element) {
		if samePolarity {
			return PianCai
		}
		return ZhengCai
	}

	// Other Controls Me (Officer)
	if isControlling(otherAttr.Element, meAttr.Element) {
		if samePolarity {
			return QiSha
		}
		return ZhengGuan
	}

	return ""
}

func isProducing(from, to Element) bool {
	switch from {
	case Wood:
		return to == Fire
	case Fire:
		return to == Earth
	case Earth:
		return to == Metal
	case Metal:
		return to == Water
	case Water:
		return to == Wood
	}
	return false
}

func isControlling(from, to Element) bool {
	switch from {
	case Wood:
		return to == Earth
	case Fire:
		return to == Metal
	case Earth:
		return to == Water
	case Metal:
		return to == Wood
	case Water:
		return to == Fire
	}
	return false
}
