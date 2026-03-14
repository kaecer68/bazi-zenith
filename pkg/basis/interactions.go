package basis

// InteractionType represents the type of relation between two pillars or branches.
type InteractionType string

const (
	Chong             InteractionType = "沖"
	He                InteractionType = "合"
	BranchHaiRelation InteractionType = "害"
)

// GetBranchChong returns true if two branches clash (六沖).
func GetBranchChong(b1, b2 Branch) bool {
	chongMap := map[Branch]Branch{
		Zi:   WuB,
		Chou: Wei,
		YinB: Shen,
		Mao:  You,
		Chen: Xu,
		Si:   Hai,
		WuB:  Zi,
		Wei:  Chou,
		Shen: YinB,
		You:  Mao,
		Xu:   Chen,
		Hai:  Si,
	}
	return chongMap[b1] == b2
}

// GetBranchHe returns true if two branches combine (六合).
func GetBranchHe(b1, b2 Branch) bool {
	heMap := map[Branch]Branch{
		Zi:   Chou,
		YinB: Hai,
		Mao:  Xu,
		Chen: You,
		Si:   Shen,
		WuB:  Wei,
		Chou: Zi,
		Hai:  YinB,
		Xu:   Mao,
		You:  Chen,
		Shen: Si,
		Wei:  WuB,
	}
	return heMap[b1] == b2
}

// GetBranchHai returns true if two branches harm each other (六害).
func GetBranchHai(b1, b2 Branch) bool {
	haiMap := map[Branch]Branch{
		Zi:   Wei,
		Chou: WuB,
		YinB: Si,
		Mao:  Chen,
		Chen: Mao,
		Si:   YinB,
		WuB:  Chou,
		Wei:  Zi,
		Shen: Hai,
		You:  Xu,
		Xu:   You,
		Hai:  Shen,
	}
	// Fixing the typo for Chen in the map if I used raw strings or variables incorrectly
	// Using Chen variable
	return haiMap[b1] == b2
}
