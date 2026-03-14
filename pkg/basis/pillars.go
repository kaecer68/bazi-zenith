package basis

// Gender represents the sex of the individual.
type Gender string

const (
	Male   Gender = "乾"
	Female Gender = "坤"
)

// Pillar represents a pair of Stem and Branch.
type Pillar struct {
	Stem   Stem
	Branch Branch
}

// JiaziList is the 60 Jiazi sequence.
var JiaziList = []Pillar{
	{Jia, Zi}, {Yi, Chou}, {Bing, YinB}, {Ding, Mao}, {Wu, Chen}, {Ji, Si}, {Geng, WuB}, {Xin, Wei}, {Ren, Shen}, {Gui, You},
	{Jia, Xu}, {Yi, Hai}, {Bing, Zi}, {Ding, Chou}, {Wu, YinB}, {Ji, Mao}, {Geng, Chen}, {Xin, Si}, {Ren, WuB}, {Gui, Wei},
	{Jia, Shen}, {Yi, You}, {Bing, Xu}, {Ding, Hai}, {Wu, Zi}, {Ji, Chou}, {Geng, YinB}, {Xin, Mao}, {Ren, Chen}, {Gui, Si},
	{Jia, WuB}, {Yi, Wei}, {Bing, Shen}, {Ding, You}, {Wu, Xu}, {Ji, Hai}, {Geng, Zi}, {Xin, Chou}, {Ren, YinB}, {Gui, Mao},
	{Jia, Chen}, {Yi, Si}, {Bing, WuB}, {Ding, Wei}, {Wu, Shen}, {Ji, You}, {Geng, Xu}, {Xin, Hai}, {Ren, Zi}, {Gui, Chou},
	{Jia, YinB}, {Yi, Mao}, {Bing, Chen}, {Ding, Si}, {Wu, WuB}, {Ji, Wei}, {Geng, Shen}, {Xin, You}, {Ren, Xu}, {Gui, Hai},
}

// GetPillarIndex returns the index of a pillar in the 60 Jiazi sequence.
func GetPillarIndex(p Pillar) int {
	for i, j := range JiaziList {
		if j.Stem == p.Stem && j.Branch == p.Branch {
			return i
		}
	}
	return -1
}

// NextPillar returns the next pillar in the 60 Jiazi sequence.
func NextPillar(p Pillar) Pillar {
	idx := GetPillarIndex(p)
	if idx == -1 {
		return Pillar{}
	}
	return JiaziList[(idx+1)%60]
}

// PrevPillar returns the previous pillar in the 60 Jiazi sequence.
func PrevPillar(p Pillar) Pillar {
	idx := GetPillarIndex(p)
	if idx == -1 {
		return Pillar{}
	}
	return JiaziList[(idx+59)%60]
}
