package basis

// Stem represents a Heavenly Stem (天干)
type Stem string

const (
	Jia  Stem = "甲"
	Yi   Stem = "乙"
	Bing Stem = "丙"
	Ding Stem = "丁"
	Wu   Stem = "戊"
	Ji   Stem = "己"
	Geng Stem = "庚"
	Xin  Stem = "辛"
	Ren  Stem = "壬"
	Gui  Stem = "癸"
)

// Branch represents an Earthly Branch (地支)
type Branch string

const (
	Zi   Branch = "子"
	Chou Branch = "丑"
	YinB Branch = "寅"
	Mao  Branch = "卯"
	Chen Branch = "辰"
	Si   Branch = "巳"
	WuB  Branch = "午"
	Wei  Branch = "未"
	Shen Branch = "申"
	You  Branch = "酉"
	Xu   Branch = "戌"
	Hai  Branch = "亥"
)

// Element represents the Five Elements (五行)
type Element string

const (
	Wood  Element = "木"
	Fire  Element = "火"
	Earth Element = "土"
	Metal Element = "金"
	Water Element = "水"
)

// Polarity represents Yin or Yang (陰陽)
type Polarity string

const (
	Yang Polarity = "陽"
	Yin  Polarity = "陰"
)

// StemAttr stores the properties of a Heavenly Stem
type StemAttr struct {
	Element  Element
	Polarity Polarity
}

var StemAttributes = map[Stem]StemAttr{
	Jia:  {Element: Wood, Polarity: Yang},
	Yi:   {Element: Wood, Polarity: Yin},
	Bing: {Element: Fire, Polarity: Yang},
	Ding: {Element: Fire, Polarity: Yin},
	Wu:   {Element: Earth, Polarity: Yang},
	Ji:   {Element: Earth, Polarity: Yin},
	Geng: {Element: Metal, Polarity: Yang},
	Xin:  {Element: Metal, Polarity: Yin},
	Ren:  {Element: Water, Polarity: Yang},
	Gui:  {Element: Water, Polarity: Yin},
}

// Attr returns the attributes of the stem
func (s Stem) Attr() StemAttr {
	return StemAttributes[s]
}
