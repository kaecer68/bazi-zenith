package basis

// LifeStage represents one of the twelve stages of life (長生十二運).
type LifeStage string

const (
	ChangSheng LifeStage = "長生"
	MuYu       LifeStage = "沐浴"
	GuanDai    LifeStage = "冠帶"
	LinGuan    LifeStage = "臨官"
	DiWang     LifeStage = "帝旺"
	Shuai      LifeStage = "衰"
	BingS      LifeStage = "病"
	SiS        LifeStage = "死"
	Mu         LifeStage = "墓"
	Jue        LifeStage = "絕"
	Tai        LifeStage = "胎"
	YangS      LifeStage = "養"
)

var stages = []LifeStage{
	ChangSheng, MuYu, GuanDai, LinGuan, DiWang, Shuai,
	BingS, SiS, Mu, Jue, Tai, YangS,
}

var branchOrder = []Branch{
	Zi, Chou, YinB, Mao, Chen, Si, WuB, Wei, Shen, You, Xu, Hai,
}

// GetLifeStage calculates the life stage of a Stem relative to a Branch.
func GetLifeStage(s Stem, b Branch) LifeStage {
	// Starting branch for each Stem's "Chang Sheng" (長生) stage
	startBranchMap := map[Stem]Branch{
		Jia:  Hai,
		Bing: YinB,
		Wu:   YinB,
		Geng: Si,
		Ren:  Shen,
		Yi:   WuB,
		Ding: You,
		Ji:   You,
		Xin:  Zi,
		Gui:  Mao,
	}

	startBranch := startBranchMap[s]
	polarity := s.Attr().Polarity

	// Find index of startBranch and b
	startIndex := -1
	targetIndex := -1
	for i, br := range branchOrder {
		if br == startBranch {
			startIndex = i
		}
		if br == b {
			targetIndex = i
		}
	}

	if startIndex == -1 || targetIndex == -1 {
		return ""
	}

	var diff int
	if polarity == Yang {
		// Clockwise: (target - start + 12) % 12
		diff = (targetIndex - startIndex + 12) % 12
	} else {
		// Counter-clockwise: (start - target + 12) % 12
		diff = (startIndex - targetIndex + 12) % 12
	}

	return stages[diff]
}
