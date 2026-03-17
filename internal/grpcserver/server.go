package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/kaecer68/bazi-zenith/gen/bazipb"
	"github.com/kaecer68/bazi-zenith/pkg/basis"
	"github.com/kaecer68/bazi-zenith/pkg/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server implements the BaziService gRPC server
type Server struct {
	bazipb.UnimplementedBaziServiceServer
	engine *engine.BaziEngine
}

// NewServer creates a new gRPC server instance
func NewServer() *Server {
	return &Server{
		engine: engine.NewBaziEngine(),
	}
}

// GetChart implements the BaziService.GetChart RPC method
func (s *Server) GetChart(ctx context.Context, req *bazipb.GetChartRequest) (*bazipb.GetChartResponse, error) {
	// Parse datetime
	loc, _ := time.LoadLocation("Asia/Taipei")
	birthTime, err := time.ParseInLocation("2006-01-02 15:04", req.Datetime, loc)
	if err != nil {
		// Try alternative format with T separator
		birthTime, err = time.Parse(time.RFC3339, req.Datetime)
		if err != nil {
			return nil, fmt.Errorf("invalid datetime format: %v", err)
		}
	}

	// Determine gender
	gender := basis.Male
	if req.Gender == "female" {
		gender = basis.Female
	}

	// Get target year
	targetYear := int(req.TargetYear)
	if targetYear == 0 {
		targetYear = time.Now().Year()
	}

	// Generate chart
	chart := s.engine.GetBaziChart(birthTime, gender)
	advice := chart.GenerateInterpretations(targetYear)

	// Convert to proto response
	return convertToProtoResponse(chart, advice, req.Gender), nil
}

// convertToProtoResponse converts internal BaziChart to proto response
func convertToProtoResponse(chart engine.BaziChart, advice []engine.Interpretation, gender string) *bazipb.GetChartResponse {
	resp := &bazipb.GetChartResponse{
		Gender:    gender,
		DayStem:   string(chart.DayStem),
		Pillars:   make(map[string]*bazipb.PillarData),
		StartAgeY: int32(chart.StartAgeY),
		StartAgeM: int32(chart.StartAgeM),
		Strength: &bazipb.StrengthAnalysis{
			Score:      chart.Strength.Score,
			Status:     chart.Strength.Status,
			IsDeLing:   chart.Strength.IsDeLing,
			IsDeDi:     chart.Strength.IsDeDi,
			IsDeZhu:    chart.Strength.IsDeZhu,
			Percentage: chart.Strength.Percentage,
		},
	}

	// Convert pillars
	mapPillar := func(name string, d engine.PillarDetail) {
		hStems := make([]string, len(d.HiddenStems))
		for i, h := range d.HiddenStems {
			hStems[i] = string(h.Stem)
		}
		tgHidden := make([]string, len(d.TenGodHidden))
		for i, t := range d.TenGodHidden {
			tgHidden[i] = string(t)
		}
		sSha := make([]string, len(d.ShenSha))
		for i, s := range d.ShenSha {
			sSha[i] = string(s)
		}

		resp.Pillars[name] = &bazipb.PillarData{
			Stem:         string(d.Pillar.Stem),
			Branch:       string(d.Pillar.Branch),
			TenGodStem:   string(d.TenGodStem),
			HiddenStems:  hStems,
			TenGodHidden: tgHidden,
			NaYin:        string(d.NaYin),
			LifeStage:    string(d.LifeStage),
			ShenSha:      sSha,
		}
	}

	mapPillar("year", chart.YearPillar)
	mapPillar("month", chart.MonthPillar)
	mapPillar("day", chart.DayPillar)
	mapPillar("hour", chart.HourPillar)

	// Convert DaYun
	for _, dy := range chart.DaYun {
		resp.DaYun = append(resp.DaYun, &bazipb.DaYunData{
			Pillar:   string(dy.Pillar.Stem) + string(dy.Pillar.Branch),
			StartAge: int32(dy.StartAge),
		})
	}

	// Convert advice
	for _, a := range advice {
		resp.Advice = append(resp.Advice, &bazipb.Interpretation{
			Title:   a.Title,
			Content: a.Content,
			Type:    a.Type,
		})
	}

	return resp
}

// Start starts the gRPC server
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
