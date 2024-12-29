package config

type SMTP struct {
	From     string `env:"MAILER_SMTP_FROM"`
	Host     string `env:"MAILER_SMTP_HOST"`
	Port     int    `env:"MAILER_SMTP_PORT" envDefault:"587"`
	User     string `env:"MAILER_SMTP_USER"`
	Password string `env:"MAILER_SMTP_PASSWORD"`
}
