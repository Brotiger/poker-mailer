package config

import (
	"log"

	"github.com/caarlos0/env/v7"
)

type Config struct {
	Nats *Nats
}

var Cfg *Config

func init() {
	natsConfig := &Nats{}

	if err := env.Parse(natsConfig); err != nil {
		log.Fatalf("failed to parse nats config, error: %v", err)
	}

	Cfg = &Config{
		Nats: natsConfig,
	}
}
