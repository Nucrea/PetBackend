all: release

release:
	go build -o ./.build/release/backend main.go