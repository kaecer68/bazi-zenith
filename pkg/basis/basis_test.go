package basis

import (
	"testing"
)

func TestGetHiddenStems(t *testing.T) {
	tests := []struct {
		branch Branch
		want   []Stem
	}{
		{Zi, []Stem{Gui}},
		{Chou, []Stem{Ji, Gui, Xin}},
		{YinB, []Stem{Jia, Bing, Wu}},
	}

	for _, tt := range tests {
		stems := GetHiddenStems(tt.branch)
		if len(stems) != len(tt.want) {
			t.Errorf("GetHiddenStems(%s) got %d stems, want %d", tt.branch, len(stems), len(tt.want))
			continue
		}
		for i, s := range stems {
			if s.Stem != tt.want[i] {
				t.Errorf("GetHiddenStems(%s)[%d] = %s, want %s", tt.branch, i, s.Stem, tt.want[i])
			}
		}
	}
}

func TestGetTenGod(t *testing.T) {
	tests := []struct {
		me    Stem
		other Stem
		want  TenGod
	}{
		{Jia, Jia, BiJian},     // Same Element, Same Polarity
		{Jia, Yi, JieCai},      // Same Element, Diff Polarity
		{Jia, Bing, ShiShen},   // Me (Wood) produces Other (Fire), Same Polarity
		{Jia, Ding, ShangGuan}, // Me (Wood) produces Other (Fire), Diff Polarity
		{Jia, Wu, PianCai},     // Me (Wood) controls Other (Earth), Same Polarity
		{Jia, Ji, ZhengCai},    // Me (Wood) controls Other (Earth), Diff Polarity
		{Jia, Geng, QiSha},     // Other (Metal) controls Me (Wood), Same Polarity
		{Jia, Xin, ZhengGuan},  // Other (Metal) controls Me (Wood), Diff Polarity
		{Jia, Ren, PianYin},    // Other (Water) produces Me (Wood), Same Polarity
		{Jia, Gui, ZhengYin},   // Other (Water) produces Me (Wood), Diff Polarity

		{Bing, Ren, QiSha},     // Other (Water) controls Me (Fire), Same Polarity
		{Ding, Gui, QiSha},     // Other (Water) controls Me (Fire), Same Polarity
		{Xin, Bing, ZhengGuan}, // Other (Fire) controls Me (Metal), Diff Polarity
	}

	for _, tt := range tests {
		got := GetTenGod(tt.me, tt.other)
		if got != tt.want {
			t.Errorf("GetTenGod(me:%s, other:%s) = %s, want %s", tt.me, tt.other, got, tt.want)
		}
	}
}

func TestGetNaYin(t *testing.T) {
	tests := []struct {
		s    Stem
		b    Branch
		want NaYin
	}{
		{Jia, Zi, "海中金"},
		{Bing, YinB, "爐中火"},
		{Gui, Hai, "大海水"},
	}

	for _, tt := range tests {
		got := GetNaYin(tt.s, tt.b)
		if got != tt.want {
			t.Errorf("GetNaYin(%s, %s) = %s, want %s", tt.s, tt.b, got, tt.want)
		}
	}
}

func TestGetLifeStage(t *testing.T) {
	tests := []struct {
		s    Stem
		b    Branch
		want LifeStage
	}{
		{Jia, Hai, ChangSheng}, // Jia (Yang Wood) starts at Hai
		{Jia, WuB, SiS},        // Jia (Yang Wood)死 at Wu
		{Yi, WuB, ChangSheng},  // Yi (Yin Wood) starts at Wu
		{Yi, Hai, SiS},         // Yi (Yin Wood)死 at Hai
		{Bing, YinB, ChangSheng},
		{Ren, Shen, ChangSheng},
	}

	for _, tt := range tests {
		got := GetLifeStage(tt.s, tt.b)
		if got != tt.want {
			t.Errorf("GetLifeStage(%s, %s) = %s, want %s", tt.s, tt.b, got, tt.want)
		}
	}
}
