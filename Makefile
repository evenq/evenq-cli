.PHONY: build

build:
	go build -ldflags="-s -w" -o dist/evenq-cli src/main.go
