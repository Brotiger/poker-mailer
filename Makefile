SHELL := /bin/bash

build:
	go build -o mailer ./cmd/main.go
.PHONY: build

up:
	./mailer
.PHONY: up

run:
	go run ./cmd
.PHONY: run