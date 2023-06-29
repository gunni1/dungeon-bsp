.PHONY: build clean

server:
	go build ./...
	go build -o bin/server cmd/server/main.go

cli:
	go build ./...
	go build -o bin/cli cmd/cli/main.go

