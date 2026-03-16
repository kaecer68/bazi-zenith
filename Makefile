.PHONY: proto build-rest build-grpc build-all clean

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

build-all: build-cli build-rest build-grpc

# Run targets
run-rest:
	go run ./cmd/bazi-server -port 8080

run-grpc:
	go run ./cmd/bazi-grpc -port 50051

# Clean
clean:
	rm -rf bin/

# Test
test:
	go test ./...
