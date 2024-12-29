package config

type Nats struct {
	Addr                 string `env:"MAILER_NATS_ADDR" envDefault:"localhost:4222"`
	ClientCert           string `env:"MAILER_NATS_CLIENT_CERT"`
	ClientKey            string `env:"MAILER_NATS_CLIENT_KEY"`
	RootCA               string `env:"MAILER_NATS_ROOT_CA"`
	ConsumerName         string `env:"MAILER_NATS_CONSUMER_NAME" envDefault:"mailer"`
	Stream               string `env:"MAILER_NATS_STREAM" envDefault:"mailer"`
	BatchSize            int    `env:"MAILER_NATS_BATCH_SIZE" envDefault:"1"`
	ReconnectWait        int    `env:"MAILER_NATS_RECONNECT_WAIT" envDefault:"10000"`
	PingInterval         int    `env:"MAILER_NATS_PING_INTERVAL" envDefault:"20000"`
	MaxReconnects        int    `env:"MAILER_NATS_MAX_RECONNECTS" envDefault:"10"`
	RetryOnFailedConnect bool   `env:"MAILER_NATS_RETRY_ON_FAILED_CONNECT" envDefault:"true"`
	MaxPingOutstanding   int    `env:"MAILER_NATS_MAX_PING_OUTSTANDING" envDefault:"5"`
}
