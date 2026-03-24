package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/kaecer68/bazi-zenith/gen/bazipb"
	v1 "github.com/kaecer68/bazi-zenith/pkg/api/v1"
	"github.com/kaecer68/bazi-zenith/pkg/basis"
	"github.com/kaecer68/bazi-zenith/pkg/engine"
)

type baziServer struct {
	bazipb.UnimplementedBaziServiceServer
	engine *engine.BaziEngine
}

func (s *baziServer) GetChart(_ context.Context, req *bazipb.GetChartRequest) (*bazipb.GetChartResponse, error) {
	if req.Datetime == "" {
		return nil, status.Error(codes.InvalidArgument, "datetime is required (format: YYYY-MM-DD HH:mm)")
	}

	loc, _ := time.LoadLocation("Asia/Taipei")
	birthTime, err := time.ParseInLocation("2006-01-02 15:04", req.Datetime, loc)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid datetime format: %v", err)
	}

	gender := basis.Male
	if req.Gender == "female" {
		gender = basis.Female
	}

	targetYear := int(req.TargetYear)
	if targetYear == 0 {
		targetYear = time.Now().Year()
	}

	chart := s.engine.GetBaziChart(birthTime, gender)
	advice := chart.GenerateInterpretations(targetYear)
	apiResp := v1.FromChart(chart, advice, targetYear, birthTime.Year())

	return toProto(apiResp), nil
}

func toProto(r v1.BaziResponse) *bazipb.GetChartResponse {
	resp := &bazipb.GetChartResponse{
		Gender:    r.Gender,
		DayStem:   r.DayStem,
		StartAgeY: int32(r.StartAgeY),
		StartAgeM: int32(r.StartAgeM),
		Strength: &bazipb.StrengthAnalysis{
			Score:      r.Strength.Score,
			Status:     r.Strength.Status,
			IsDeLing:   r.Strength.IsDeLing,
			IsDeDi:     r.Strength.IsDeDi,
			IsDeZhu:    r.Strength.IsDeZhu,
			Percentage: r.Strength.Percentage,
		},
		Pillars:             make(map[string]*bazipb.PillarData),
		FavorableElements:   r.FavorableElements,
		UnfavorableElements: r.UnfavorableElements,
		Directions: &bazipb.Directions{
			Wealth:       r.Directions.Wealth,
			Career:       r.Directions.Career,
			Study:        r.Directions.Study,
			Relationship: r.Directions.Relationship,
		},
	}

	for name, p := range r.Pillars {
		resp.Pillars[name] = &bazipb.PillarData{
			Stem:         p.Stem,
			Branch:       p.Branch,
			TenGodStem:   p.TenGodStem,
			HiddenStems:  p.HiddenStems,
			TenGodHidden: p.TenGodHidden,
			NaYin:        p.NaYin,
			LifeStage:    p.LifeStage,
			ShenSha:      p.ShenSha,
		}
	}

	for _, dy := range r.DaYun {
		resp.DaYun = append(resp.DaYun, &bazipb.DaYunData{
			Pillar:   dy.Pillar,
			StartAge: int32(dy.StartAge),
		})
	}

	for _, a := range r.Advice {
		resp.Advice = append(resp.Advice, &bazipb.Interpretation{
			Title:   a.Title,
			Content: a.Content,
			Type:    a.Type,
		})
	}

	resp.DetailChart = convertDetailChartToProto(r.DetailChart)

	return resp
}

func convertDetailChartToProto(dc v1.DetailChart) *bazipb.DetailChart {
	protoDC := &bazipb.DetailChart{
		Natal:            convertNatalMatrixToProto(dc.Natal),
		FiveElementState: dc.FiveElementState,
		Prompts: &bazipb.DetailPrompts{
			Tiangan: dc.Prompts.Tiangan,
			Dizhi:   dc.Prompts.Dizhi,
		},
	}

	for _, item := range dc.DayunBoard {
		protoDC.DayunBoard = append(protoDC.DayunBoard, &bazipb.BoardItem{
			Index:        int32(item.Index),
			Year:         int32(item.Year),
			StartAge:     int32(item.StartAge),
			StartYear:    int32(item.StartYear),
			Pillar:       item.Pillar,
			TenGodStem:   item.TenGodStem,
			TenGodBranch: item.TenGodBranch,
		})
	}

	for _, item := range dc.LiunianBoard {
		protoDC.LiunianBoard = append(protoDC.LiunianBoard, &bazipb.BoardItem{
			Year:         int32(item.Year),
			Pillar:       item.Pillar,
			TenGodStem:   item.TenGodStem,
			TenGodBranch: item.TenGodBranch,
		})
	}

	for _, item := range dc.LiuyueBoard {
		protoDC.LiuyueBoard = append(protoDC.LiuyueBoard, &bazipb.MonthBoardItem{
			Month:        int32(item.Month),
			Pillar:       item.Pillar,
			TenGodStem:   item.TenGodStem,
			TenGodBranch: item.TenGodBranch,
		})
	}

	return protoDC
}

func convertNatalMatrixToProto(nm v1.NatalMatrix) *bazipb.NatalMatrix {
	return &bazipb.NatalMatrix{
		TenGodStem: &bazipb.FourPillarsText{Year: nm.TenGodStem.Year, Month: nm.TenGodStem.Month, Day: nm.TenGodStem.Day, Hour: nm.TenGodStem.Hour},
		TianGan:    &bazipb.FourPillarsText{Year: nm.TianGan.Year, Month: nm.TianGan.Month, Day: nm.TianGan.Day, Hour: nm.TianGan.Hour},
		DiZhi:      &bazipb.FourPillarsText{Year: nm.DiZhi.Year, Month: nm.DiZhi.Month, Day: nm.DiZhi.Day, Hour: nm.DiZhi.Hour},
		CangGan:    &bazipb.FourPillarsText{Year: nm.CangGan.Year, Month: nm.CangGan.Month, Day: nm.CangGan.Day, Hour: nm.CangGan.Hour},
		NaYin:      &bazipb.FourPillarsText{Year: nm.NaYin.Year, Month: nm.NaYin.Month, Day: nm.NaYin.Day, Hour: nm.NaYin.Hour},
		XingYun:    &bazipb.FourPillarsText{Year: nm.XingYun.Year, Month: nm.XingYun.Month, Day: nm.XingYun.Day, Hour: nm.XingYun.Hour},
		ZiZuo:      &bazipb.FourPillarsText{Year: nm.ZiZuo.Year, Month: nm.ZiZuo.Month, Day: nm.ZiZuo.Day, Hour: nm.ZiZuo.Hour},
		KongWang:   &bazipb.FourPillarsText{Year: nm.KongWang.Year, Month: nm.KongWang.Month, Day: nm.KongWang.Day, Hour: nm.KongWang.Hour},
	}
}

func main() {
	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		log.Fatal("GRPC_PORT environment variable is required. Please ensure .env.ports is loaded or set GRPC_PORT directly.")
	}

	defaultPort, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid GRPC_PORT value: %v", err)
	}

	port := flag.Int("port", defaultPort, "gRPC server port (overridden by GRPC_PORT env)")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	bazipb.RegisterBaziServiceServer(s, &baziServer{
		engine: engine.NewBaziEngine(),
	})
	reflection.Register(s)

	log.Printf("Bazi-Zenith gRPC server listening on :%d", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
