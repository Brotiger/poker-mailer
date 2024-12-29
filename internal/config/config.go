package config

import (
	"log"

	"github.com/caarlos0/env/v7"
)

type Config struct {
	App  *App
	Nats *Nats
	SMTP *SMTP
}

var Cfg *Config

func init() {
	appConfig := &App{}
	natsConfig := &Nats{}
	SMTPConfig := &SMTP{}

	if err := env.Parse(appConfig); err != nil {
		log.Fatalf("failed to parse app config, error: %v", err)
	}

	if err := env.Parse(natsConfig); err != nil {
		log.Fatalf("failed to parse nats config, error: %v", err)
	}

	if err := env.Parse(SMTPConfig); err != nil {
		log.Fatalf("failed to parse smtp config, error: %v", err)
	}

	Cfg = &Config{
		App:  appConfig,
		Nats: natsConfig,
		SMTP: SMTPConfig,
	}
}
