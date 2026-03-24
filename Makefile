.PHONY: proto build-rest build-grpc build-cli build-all build-web clean dev sync verify sync-contracts verify-contracts dev-clean smoke-prompts

# Proto code generation
proto:
	protoc \
		--proto_path=api/proto \
		--go_out=. --go_opt=module=github.com/kaecer68/bazi-zenith \
		--go-grpc_out=. --go-grpc_opt=module=github.com/kaecer68/bazi-zenith \
		api/proto/bazi/v1/bazi.proto
	@echo "Proto files generated."

# Install protoc plugins (run once)
proto-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Build targets
build-cli:
	go build -o bin/bazi-cli ./cmd/bazi-cli

build-rest:
	go build -o bin/bazi-server ./cmd/bazi-server

build-grpc:
	go build -o bin/bazi-grpc ./cmd/bazi-grpc

build-web:
	cd web && npm run build

build-all: build-cli build-rest build-grpc build-web

# Run targets
run-rest:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh
	bash -c 'set -a; . ./.env.ports; set +a; go run ./cmd/bazi-server -port "$$REST_PORT"'

run-grpc:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh
	bash -c 'set -a; . ./.env.ports; set +a; go run ./cmd/bazi-grpc -port "$$GRPC_PORT"'

dev-web:
	cd web && npm run dev

# Development mode (backend + frontend)
dev:
	@echo "Starting development servers..."
	@chmod +x scripts/sync-contracts.sh && bash scripts/sync-contracts.sh
	@bash -c 'set -a; . ./.env.ports; set +a; echo "Backend: http://localhost:$$REST_PORT"; echo "Frontend: http://localhost:5173"'
	@make run-rest &
	@sleep 2
	@make dev-web

sync-contracts:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh

sync: sync-contracts

verify-contracts:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh --check

verify: verify-contracts

dev-clean:
	@chmod +x scripts/dev-clean.sh
	bash scripts/dev-clean.sh

smoke-prompts:
	@chmod +x scripts/smoke-prompts-consistency.sh
	bash scripts/smoke-prompts-consistency.sh

# Clean
clean:
	rm -rf bin/
	rm -rf web/dist/

# Test
test:
	go test ./...
