package basis

// GetYearPillar returns the sexagenary pillar for a given Gregorian year.
// Note: In Bazi, the year pillar changes at "Li Chun" (立春), not Jan 1st.
// This function returns the base pillar for the year number.
func GetYearPillar(year int) Pillar {
	// 4 AD was a Jia-Zi year (index 0).
	// (year - 4) % 60 gives the index.
	// For years BC or near epoch, handles negative.
	idx := (year - 4) % 60
	if idx < 0 {
		idx += 60
	}
	return JiaziList[idx]
}

// GetMonthStems returns the 12 month stems for a given year stem using "Wu Hu Xun" (五虎遁).
// Bazi months always start from the Tiger month (YinB).
func GetMonthStems(yearStem Stem) []Stem {
	var startStem Stem
	switch yearStem {
	case Jia, Ji:
		startStem = Bing
	case Yi, Geng:
		startStem = Wu
	case Bing, Xin:
		startStem = Geng
	case Ding, Ren:
		startStem = Ren
	case Wu, Gui:
		startStem = Jia
	}

	stems := make([]Stem, 12)
	// Find index of startStem in Stem sequence
	stemOrder := []Stem{Jia, Yi, Bing, Ding, Wu, Ji, Geng, Xin, Ren, Gui}
	startIdx := -1
	for i, s := range stemOrder {
		if s == startStem {
			startIdx = i
			break
		}
	}

	for i := 0; i < 12; i++ {
		stems[i] = stemOrder[(startIdx+i)%10]
	}
	return stems
}

// GetMonthPillar returns the pillar for a specific month (1-12, where 1 is the Yin/Tiger month).
func GetMonthPillar(yearStem Stem, month int) Pillar {
	if month < 1 || month > 12 {
		return Pillar{}
	}
	stems := GetMonthStems(yearStem)
	return Pillar{
		Stem:   stems[month-1],
		Branch: branchOrder[(month+1)%12], // month 1 -> YinB (idx 2), 2 -> Mao (idx 3)...
	}
}
