package config

type App struct {
	GracefulShutdownTimeoutMS int `env:"MAILER_APP_GRACEFUL_SHUTDOWN_TIMEOUT_MS" envDefault:"5000"`
}
