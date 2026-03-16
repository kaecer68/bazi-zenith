package engine

import (
	"fmt"
	"testing"
	"time"

	"github.com/kaecer68/bazi-zenith/pkg/basis"
)

func TestGetBaziChart(t *testing.T) {
	engine := NewBaziEngine()
	// 2024-03-14 19:16:00 (UTC+8)
	loc, _ := time.LoadLocation("Asia/Taipei")
	birthTime := time.Date(2024, 3, 14, 19, 16, 0, 0, loc)

	chart := engine.GetBaziChart(birthTime, basis.Male)

	fmt.Printf("Year: %s%s (%s) %v\n", chart.YearPillar.Pillar.Stem, chart.YearPillar.Pillar.Branch, chart.YearPillar.NaYin, chart.YearPillar.ShenSha)
	fmt.Printf("Month: %s%s (%s) %v\n", chart.MonthPillar.Pillar.Stem, chart.MonthPillar.Pillar.Branch, chart.MonthPillar.NaYin, chart.MonthPillar.ShenSha)
	fmt.Printf("Day: %s%s (%s) %v\n", chart.DayPillar.Pillar.Stem, chart.DayPillar.Pillar.Branch, chart.DayPillar.NaYin, chart.DayPillar.ShenSha)
	fmt.Printf("Hour: %s%s (%s) %v\n", chart.HourPillar.Pillar.Stem, chart.HourPillar.Pillar.Branch, chart.HourPillar.NaYin, chart.HourPillar.ShenSha)
	fmt.Printf("Start Age: %d years %d months\n", chart.StartAgeY, chart.StartAgeM)
	fmt.Printf("Strength: %s (Score: %.1f)\n", chart.Strength.Status, chart.Strength.Score)

	interpretations2024 := chart.GenerateInterpretations(2024)
	fmt.Println("\n--- 命理斷語 (2024 甲辰) ---")
	for _, it := range interpretations2024 {
		fmt.Printf("[%s] %s: %s\n", it.Type, it.Title, it.Content)
	}

	interpretations2025 := chart.GenerateInterpretations(2025)
	fmt.Println("\n--- 命理斷語 (2025 乙巳) ---")
	for _, it := range interpretations2025 {
		fmt.Printf("[%s] %s: %s\n", it.Type, it.Title, it.Content)
	}

	if chart.YearPillar.Pillar.Stem != basis.Jia || chart.YearPillar.Pillar.Branch != basis.Chen {
		t.Errorf("2024 should be Jia-Chen, got %s%s", chart.YearPillar.Pillar.Stem, chart.YearPillar.Pillar.Branch)
	}
}

func TestVerifyYearBug1972(t *testing.T) {
	engine := NewBaziEngine()
	loc, _ := time.LoadLocation("Asia/Taipei")
	birthTime := time.Date(1972, 6, 15, 12, 0, 0, 0, loc)

	chart := engine.GetBaziChart(birthTime, basis.Male)

	yearPillar := string(chart.YearPillar.Pillar.Stem) + string(chart.YearPillar.Pillar.Branch)
	expected := "壬子"
	if yearPillar != expected {
		t.Errorf("1972 Year Pillar should be %s, but got %s", expected, yearPillar)
	}
}

func TestVerifyYearBug1990(t *testing.T) {
	engine := NewBaziEngine()
	loc, _ := time.LoadLocation("Asia/Taipei")
	birthTime := time.Date(1990, 6, 15, 12, 0, 0, 0, loc)

	chart := engine.GetBaziChart(birthTime, basis.Male)

	yearPillar := string(chart.YearPillar.Pillar.Stem) + string(chart.YearPillar.Pillar.Branch)
	expected := "庚午"
	if yearPillar != expected {
		t.Errorf("1990 Year Pillar should be %s, but got %s", expected, yearPillar)
	}
}
