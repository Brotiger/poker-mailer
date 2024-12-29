SHELL := /bin/bash

build:
	go build -o core_api ./cmd/main.go
.PHONY: build

up:
	./core_api
.PHONY: up

run:
	go run ./cmd
.PHONY: run