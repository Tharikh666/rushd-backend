.PHONY: proto generate build run test tidy

proto:
	buf generate

generate: proto

tidy:
	go mod tidy

build:
	go build -o bin/rushd-server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

fmt:
	go fmt ./...