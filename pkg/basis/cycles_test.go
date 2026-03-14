package basis

import "testing"

func TestGetYearPillar(t *testing.T) {
	tests := []struct {
		year int
		want Pillar
	}{
		{2024, Pillar{Jia, Chen}}, // 甲辰
		{2025, Pillar{Yi, Si}},    // 乙巳
		{1984, Pillar{Jia, Zi}},   // 甲子
	}

	for _, tt := range tests {
		got := GetYearPillar(tt.year)
		if got.Stem != tt.want.Stem || got.Branch != tt.want.Branch {
			t.Errorf("GetYearPillar(%d) = %s%s, want %s%s", tt.year, got.Stem, got.Branch, tt.want.Stem, tt.want.Branch)
		}
	}
}

func TestGetMonthPillar(t *testing.T) {
	// Jia Year (甲年) -> Yin month starts with Bing-Yin (丙寅)
	yearStem := Jia

	p1 := GetMonthPillar(yearStem, 1) // First month (Tiger)
	if p1.Stem != Bing || p1.Branch != YinB {
		t.Errorf("Jia Year Month 1 expected Bing-Yin, got %s%s", p1.Stem, p1.Branch)
	}

	p12 := GetMonthPillar(yearStem, 12) // Last month (Ox)
	if p12.Stem != Ding || p12.Branch != Chou {
		t.Errorf("Jia Year Month 12 expected Ding-Chou, got %s%s", p12.Stem, p12.Branch)
	}

	// Yi Year (乙年) -> Yin month starts with Wu-Yin (戊寅)
	yearStem = Yi
	p1 = GetMonthPillar(yearStem, 1)
	if p1.Stem != Wu || p1.Branch != YinB {
		t.Errorf("Yi Year Month 1 expected Wu-Yin, got %s%s", p1.Stem, p1.Branch)
	}
}
