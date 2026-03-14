package basis

// DaYunInfo represents a single stage of the Big Cycle (大運).
type DaYunInfo struct {
	Pillar   Pillar
	StartAge int // Starting age (years old)
}

// GetDaYunSequence generates the seqence of 10-year Big Cycles.
// yearStem: The stem of the birth year (to check polarity).
// monthPillar: The pillar of the birth month.
// gender: Male or Female.
func GetDaYunSequence(yearStem Stem, monthPillar Pillar, gender Gender) []Pillar {
	yearPolarity := yearStem.Attr().Polarity
	isClockwise := false

	if gender == Male {
		if yearPolarity == Yang {
			isClockwise = true
		}
	} else {
		if yearPolarity == Yin {
			isClockwise = true
		}
	}

	sequence := make([]Pillar, 0, 8)
	current := monthPillar

	for i := 0; i < 10; i++ {
		if isClockwise {
			current = NextPillar(current)
		} else {
			current = PrevPillar(current)
		}
		sequence = append(sequence, current)
	}

	return sequence
}

// CalculateDaYunAge returns the starting age for Da Yun.
// diffSeconds: Total seconds between birth and the nearest "Jie Qi" (Section Term).
// In Bazi, 3 days = 1 year, 1 day = 4 months, 1 hour (2h) = 10 days.
func CalculateDaYunAge(diffSeconds int64) (years int, months int) {
	// 1 year = 3 days = 3 * 24 * 3600 seconds = 259200 seconds
	years = int(diffSeconds / 259200)
	remaining := diffSeconds % 259200

	// 1 month = 1 day / 4 = 24 * 3600 / 4 = 21600 seconds (No, wait)
	// 1 year = 3 days. 1 month = 1/4 day = 6 hours.
	// 3 days = 72 hours. 72h / 12 months = 6 hours per month.
	// 6 hours = 6 * 3600 = 21600 seconds.
	months = int(remaining / 21600)

	return years, months
}
