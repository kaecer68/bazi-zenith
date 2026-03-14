package basis

// HiddenStem represents one of the stems hidden within an Earthly Branch.
type HiddenStem struct {
	Stem    Stem
	Type    string // 本氣, 中氣, 餘氣
	Percent int    // Roughly optional calculation
}

// BranchHiddenStems maps each Branch to its hidden Heavenly Stems.
var BranchHiddenStems = map[Branch][]HiddenStem{
	Zi: {
		{Stem: Gui, Type: "本氣"},
	},
	Chou: {
		{Stem: Ji, Type: "本氣"},
		{Stem: Gui, Type: "中氣"},
		{Stem: Xin, Type: "餘氣"},
	},
	YinB: {
		{Stem: Jia, Type: "本氣"},
		{Stem: Bing, Type: "中氣"},
		{Stem: Wu, Type: "餘氣"},
	},
	Mao: {
		{Stem: Yi, Type: "本氣"},
	},
	Chen: {
		{Stem: Wu, Type: "本氣"},
		{Stem: Yi, Type: "中氣"},
		{Stem: Gui, Type: "餘氣"},
	},
	Si: {
		{Stem: Bing, Type: "本氣"},
		{Stem: Geng, Type: "中氣"},
		{Stem: Wu, Type: "餘氣"},
	},
	WuB: {
		{Stem: Ding, Type: "本氣"},
		{Stem: Ji, Type: "中氣"},
	},
	Wei: {
		{Stem: Ji, Type: "本氣"},
		{Stem: Ding, Type: "中氣"},
		{Stem: Yi, Type: "餘氣"},
	},
	Shen: {
		{Stem: Geng, Type: "本氣"},
		{Stem: Ren, Type: "中氣"},
		{Stem: Wu, Type: "餘氣"},
	},
	You: {
		{Stem: Xin, Type: "本氣"},
	},
	Xu: {
		{Stem: Wu, Type: "本氣"},
		{Stem: Xin, Type: "中氣"},
		{Stem: Ding, Type: "餘氣"},
	},
	Hai: {
		{Stem: Ren, Type: "本氣"},
		{Stem: Jia, Type: "中氣"},
	},
}

// GetHiddenStems returns the hidden stems for a given branch.
func GetHiddenStems(b Branch) []HiddenStem {
	return BranchHiddenStems[b]
}
