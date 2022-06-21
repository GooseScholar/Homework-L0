.PHONY: build
build:
		go build -v ./cmd/pub/main.go

.DEFAULT_GOAL := build