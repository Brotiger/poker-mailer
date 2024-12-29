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

nats-streams-add:
	nats -s "${MAILER_NATS_ADDR}" str add "${MAILER_NATS_STREAM}" --subjects "mailer" --ack --retention=work --max-age=1h --defaults;\
	nats -s "${MAILER_NATS_ADDR}" con add "${MAILER_NATS_STREAM}" "${MAILER_NATS_CONSUMER_NAME}" --pull --ack explicit --deliver all --max-deliver 1 --wait=3m --defaults;
.PHONY: nats-streams-add