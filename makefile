all: release

release:
	GOEXPERIMENT=boringcrypto  go build -ldflags "-s -w" -o ./.build/release/backend main.go

run: release
	./.build/release/backend -c ./misc/config.yaml -o ./.run/log.txt -p ./.run/cpu.pprof

venv:
	python3 -m pip install --user virtualenv

locust:
	pip3 install locust