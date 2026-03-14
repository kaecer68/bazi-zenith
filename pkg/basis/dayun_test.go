package basis

import (
	"testing"
)

func TestGetDaYunSequence(t *testing.T) {
	// Example: Male, Yang Year (Jia), Month (Bing-Yin)
	// Expect clockwise: Ding-Mao, Wu-Chen, ...
	yearStem := Jia
	monthPillar := Pillar{Bing, YinB}
	gender := Male

	seq := GetDaYunSequence(yearStem, monthPillar, gender)

	if len(seq) != 10 {
		t.Errorf("Expected 10 pillars, got %d", len(seq))
	}

	if seq[0].Stem != Ding || seq[0].Branch != Mao {
		t.Errorf("First DaYun expected Ding-Mao, got %s-%s", seq[0].Stem, seq[0].Branch)
	}

	// Example: Female, Yang Year (Jia), Month (Bing-Yin)
	// Expect counter-clockwise: Yi-Chou, Jia-Zi, ...
	gender = Female
	seq = GetDaYunSequence(yearStem, monthPillar, gender)
	if seq[0].Stem != Yi || seq[0].Branch != Chou {
		t.Errorf("First DaYun expected Yi-Chou, got %s-%s", seq[0].Stem, seq[0].Branch)
	}
}

func TestCalculateDaYunAge(t *testing.T) {
	// 3 days = 1 year
	// 1 day = 24 * 3600 = 86400 seconds
	// 3 days = 259200 seconds
	y, m := CalculateDaYunAge(259200)
	if y != 1 || m != 0 {
		t.Errorf("Expected 1 year 0 months, got %d years %d months", y, m)
	}

	// 1.5 days = 0 years 6 months
	y, m = CalculateDaYunAge(129600)
	if y != 0 || m != 6 {
		t.Errorf("Expected 0 years 6 months, got %d years %d months", y, m)
	}
}
