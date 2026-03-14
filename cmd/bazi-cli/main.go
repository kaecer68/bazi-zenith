package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kaecer68/bazi-zenith/pkg/basis"
	"github.com/kaecer68/bazi-zenith/pkg/engine"
)

func main() {
	dateTimeStr := flag.String("dt", "", "Date and time in format '2006-01-02 15:04'")
	genderStr := flag.String("g", "male", "Gender: male (乾) or female (坤)")
	targetYear := flag.Int("y", time.Now().Year(), "Target year for interpretation (Liu Nian)")
	flag.Parse()

	if *dateTimeStr == "" {
		fmt.Println("用法: bazi-cli -dt 'YYYY-MM-DD HH:mm' [-g male/female] [-y 2024]")
		os.Exit(1)
	}

	loc, _ := time.LoadLocation("Asia/Taipei")
	birthTime, err := time.ParseInLocation("2006-01-02 15:04", *dateTimeStr, loc)
	if err != nil {
		fmt.Printf("日期格式錯誤: %v\n", err)
		os.Exit(1)
	}

	gender := basis.Male
	if *genderStr == "female" {
		gender = basis.Female
	}

	baziEngine := engine.NewBaziEngine()
	chart := baziEngine.GetBaziChart(birthTime, gender)

	printChart(chart, *targetYear)
}

func printChart(c engine.BaziChart, targetYear int) {
	fmt.Println("\n==========================================")
	fmt.Printf("   Bazi-Zenith (八字命盤引擎) - %s造\n", c.Gender)
	fmt.Println("==========================================")

	// Display Pillars vertically (Simplified structure for terminal)
	fmt.Printf("      【年柱】  【月柱】  【日柱】  【時柱】\n")
	fmt.Printf("十神:  %-8s  %-8s  %-8s  %-8s\n",
		c.YearPillar.TenGodStem, c.MonthPillar.TenGodStem, "日元", c.HourPillar.TenGodStem)

	fmt.Printf("天干:    %-2s      %-2s      %-2s      %-2s\n",
		c.YearPillar.Pillar.Stem, c.MonthPillar.Pillar.Stem, c.DayPillar.Pillar.Stem, c.HourPillar.Pillar.Stem)

	fmt.Printf("地支:    %-2s      %-2s      %-2s      %-2s\n",
		c.YearPillar.Pillar.Branch, c.MonthPillar.Pillar.Branch, c.DayPillar.Pillar.Branch, c.HourPillar.Pillar.Branch)

	fmt.Printf("藏干:  %-8s  %-8s  %-8s  %-8s\n",
		formatHidden(c.YearPillar), formatHidden(c.MonthPillar), formatHidden(c.DayPillar), formatHidden(c.HourPillar))

	fmt.Printf("納音:  %-8s  %-8s  %-8s  %-8s\n",
		c.YearPillar.NaYin, c.MonthPillar.NaYin, c.DayPillar.NaYin, c.HourPillar.NaYin)

	fmt.Printf("神煞:  %-8v  %-8v  %-8v  %-8v\n",
		c.YearPillar.ShenSha, c.MonthPillar.ShenSha, c.DayPillar.ShenSha, c.HourPillar.ShenSha)

	fmt.Println("------------------------------------------")
	fmt.Printf("身強分析: %s (總分: %.1f) | 起運歲數: %d 歲 %d 個月\n",
		c.Strength.Status, c.Strength.Score, c.StartAgeY, c.StartAgeM)

	fmt.Println("------------------------------------------")
	fmt.Printf("大運: ")
	for i, dy := range c.DaYun {
		if i > 5 {
			break
		} // Show first 6
		fmt.Printf("[%d]%s%s ", dy.StartAge, dy.Pillar.Stem, dy.Pillar.Branch)
	}
	fmt.Println("\n------------------------------------------")

	advice := c.GenerateInterpretations(targetYear)
	fmt.Printf("★ %d (%s%s) 流年斷語:\n", targetYear, string(basis.GetYearPillar(targetYear).Stem), string(basis.GetYearPillar(targetYear).Branch))
	for _, it := range advice {
		fmt.Printf(" ● [%s] %s: %s\n", it.Type, it.Title, it.Content)
	}
	fmt.Println("==========================================")
}

func formatHidden(d engine.PillarDetail) string {
	res := ""
	for _, h := range d.HiddenStems {
		res += string(h.Stem)
	}
	return res
}
