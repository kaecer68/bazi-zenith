#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

TARGET_DATETIME="${TARGET_DATETIME:-1990-05-15 10:30}"
TARGET_GENDER="${TARGET_GENDER:-male}"
TARGET_YEAR="${TARGET_YEAR:-2026}"

cleanup() {
  bash "$REPO_ROOT/scripts/dev-clean.sh" >/dev/null 2>&1 || true
  rm -f "${TMP_GRPC_CLIENT:-}"
}
trap cleanup EXIT

echo "[smoke] sync contracts"
bash "$REPO_ROOT/scripts/sync-contracts.sh"

# shellcheck disable=SC1090
source "$REPO_ROOT/.env.ports"
export REST_PORT GRPC_PORT BAZI_REST_PORT BAZI_GRPC_PORT

bash "$REPO_ROOT/scripts/dev-clean.sh" >/dev/null 2>&1 || true

wait_rest_ready() {
  for _ in $(seq 1 120); do
    if curl -sS "http://127.0.0.1:${REST_PORT}/health" >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.25
  done
  echo "[smoke] REST did not become ready on :${REST_PORT}" >&2
  return 1
}

wait_grpc_ready() {
  for _ in $(seq 1 120); do
    if lsof -iTCP:"${GRPC_PORT}" -sTCP:LISTEN >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.25
  done
  echo "[smoke] gRPC did not become ready on :${GRPC_PORT}" >&2
  return 1
}

echo "[smoke] start REST on :${REST_PORT}"
go run ./cmd/bazi-server -port "${REST_PORT}" >/tmp/bazi-smoke-rest.log 2>&1 &
wait_rest_ready

REST_PROMPTS="$({
  curl -sS -X POST "http://127.0.0.1:${REST_PORT}/api/v1/chart" \
    -H 'Content-Type: application/json' \
    -d "{\"datetime\":\"${TARGET_DATETIME}\",\"gender\":\"${TARGET_GENDER}\",\"target_year\":${TARGET_YEAR}}" \
  | jq -c '.detail_chart.prompts'
} | tr -d '\n')"

echo "[smoke] REST prompts: ${REST_PROMPTS}"
bash "$REPO_ROOT/scripts/dev-clean.sh" >/dev/null 2>&1 || true

echo "[smoke] start gRPC on :${GRPC_PORT}"
go run ./cmd/bazi-grpc -port "${GRPC_PORT}" >/tmp/bazi-smoke-grpc.log 2>&1 &
wait_grpc_ready

TMP_GRPC_CLIENT="$(mktemp /tmp/bazi-smoke-grpc-client-XXXX.go)"
cat > "$TMP_GRPC_CLIENT" <<'EOF'
package main

import (
  "context"
  "fmt"
  "os"
  "strconv"
  "time"

  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"

  "github.com/kaecer68/bazi-zenith/gen/bazipb"
)

func main() {
  grpcPort := os.Getenv("GRPC_PORT")
  datetime := os.Getenv("TARGET_DATETIME")
  gender := os.Getenv("TARGET_GENDER")
  yearText := os.Getenv("TARGET_YEAR")

  year, err := strconv.Atoi(yearText)
  if err != nil {
    panic(err)
  }

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  conn, err := grpc.DialContext(ctx, "127.0.0.1:"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  client := bazipb.NewBaziServiceClient(conn)
  resp, err := client.GetChart(ctx, &bazipb.GetChartRequest{
    Datetime:   datetime,
    Gender:     gender,
    TargetYear: int32(year),
  })
  if err != nil {
    panic(err)
  }

  prompts := resp.GetDetailChart().GetPrompts()
  fmt.Printf("{\"tiangan\":%q,\"dizhi\":%q}", prompts.GetTiangan(), prompts.GetDizhi())
}
EOF

GRPC_PROMPTS="$({
  GRPC_PORT="${GRPC_PORT}" TARGET_DATETIME="${TARGET_DATETIME}" TARGET_GENDER="${TARGET_GENDER}" TARGET_YEAR="${TARGET_YEAR}" go run "$TMP_GRPC_CLIENT"
} | tr -d '\n')"

echo "[smoke] gRPC prompts: ${GRPC_PROMPTS}"

if [[ "$REST_PROMPTS" != "$GRPC_PROMPTS" ]]; then
  echo "[smoke] FAIL: REST and gRPC prompts mismatch" >&2
  exit 1
fi

echo "[smoke] PASS: REST and gRPC prompts are consistent"
