.PHONY: proto build-rest build-grpc build-cli build-all build-web clean dev

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
	go run ./cmd/bazi-server -port 8080

run-grpc:
	go run ./cmd/bazi-grpc -port 50051

dev-web:
	cd web && npm run dev

# Development mode (backend + frontend)
dev:
	@echo "Starting development servers..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"
	@make run-rest &
	@sleep 2
	@make dev-web

# Clean
clean:
	rm -rf bin/
	rm -rf web/dist/

# Test
test:
	go test ./...
