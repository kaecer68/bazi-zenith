package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kaecer68/bazi-zenith/internal/grpcserver"
	v1 "github.com/kaecer68/bazi-zenith/pkg/api/v1"
	"github.com/kaecer68/bazi-zenith/pkg/basis"
	"github.com/kaecer68/bazi-zenith/pkg/engine"
)

// ChartRequest represents the JSON request body.
type ChartRequest struct {
	DateTime   string `json:"datetime"`
	Gender     string `json:"gender"`
	TargetYear int    `json:"target_year"`
}

func main() {
	port := flag.Int("port", 8080, "HTTP server port")
	flag.Parse()

	// Start gRPC server in a goroutine
	go grpcserver.Start()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/chart", handleChart)
	mux.HandleFunc("GET /health", handleHealth)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Bazi-Zenith REST server listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleChart(w http.ResponseWriter, r *http.Request) {
	var req ChartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if req.DateTime == "" {
		writeError(w, http.StatusBadRequest, "datetime is required (format: YYYY-MM-DD HH:mm)")
		return
	}

	loc, _ := time.LoadLocation("Asia/Taipei")
	birthTime, err := time.ParseInLocation("2006-01-02 15:04", req.DateTime, loc)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid datetime format: "+err.Error())
		return
	}

	gender := basis.Male
	if req.Gender == "female" {
		gender = basis.Female
	}

	targetYear := req.TargetYear
	if targetYear == 0 {
		targetYear = time.Now().Year()
	}

	e := engine.NewBaziEngine()
	chart := e.GetBaziChart(birthTime, gender)
	advice := chart.GenerateInterpretations(targetYear)
	resp := v1.FromChart(chart, advice)

	writeJSON(w, http.StatusOK, resp)
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "bazi-zenith",
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
