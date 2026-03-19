package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/kaecer68/bazi-zenith/gen/bazipb"
	v1 "github.com/kaecer68/bazi-zenith/pkg/api/v1"
	"github.com/kaecer68/bazi-zenith/pkg/basis"
	"github.com/kaecer68/bazi-zenith/pkg/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	bazipb.UnimplementedBaziServiceServer
	engine *engine.BaziEngine
}

func NewServer() *Server {
	return &Server{
		engine: engine.NewBaziEngine(),
	}
}

func (s *Server) GetChart(ctx context.Context, req *bazipb.GetChartRequest) (*bazipb.GetChartResponse, error) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	birthTime, err := time.ParseInLocation("2006-01-02 15:04", req.Datetime, loc)
	if err != nil {
		birthTime, err = time.Parse(time.RFC3339, req.Datetime)
		if err != nil {
			return nil, fmt.Errorf("invalid datetime format: %v", err)
		}
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
	apiResp := v1.FromChart(chart, advice)

	return convertToProtoResponse(apiResp), nil
}

func convertToProtoResponse(r v1.BaziResponse) *bazipb.GetChartResponse {
	resp := &bazipb.GetChartResponse{
		Gender:    r.Gender,
		DayStem:   r.DayStem,
		Pillars:   make(map[string]*bazipb.PillarData),
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

	return resp
}

func Start() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	bazipb.RegisterBaziServiceServer(s, NewServer())
	reflection.Register(s)

	log.Printf("Bazi-Zenith gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
