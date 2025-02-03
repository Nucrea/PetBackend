all: install release

release:
	GOEXPERIMENT=boringcrypto  go build -ldflags "-s -w" -o ./.build/release/backend main.go

install:
	go install

grpc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

	# protoc --go_out=. --go_opt=paths=source_relative \
    # --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    # helloworld/helloworld.proto

run: install release
	mkdir -p ./.run 
	./.build/release/backend -c ./misc/config.yaml -o ./.run/log.txt -p ./.run/cpu.pprof