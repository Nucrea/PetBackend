all: install release

release:
	GOEXPERIMENT=boringcrypto  go build -ldflags "-s -w" -o ./.build/release/backend main.go

install:
	go install

run: install release
	mkdir -p ./.run 
	./.build/release/backend -c ./misc/config.yaml -o ./.run/log.txt -p ./.run/cpu.pprof