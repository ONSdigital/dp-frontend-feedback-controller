package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-feedback-controller
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	HelloWorldEmphasise        bool          `envconfig:"HELLO_WORLD_EMPHASISE"`
	MailHost                   string        `envconfig:"MAIL_HOST"`
	MailUser                   string        `envconfig:"MAIL_USER"`
	MailPassword               string        `envconfig:"MAIL_PASSWORD" json:"-"`
	MailPort                   string        `envconfig:"MAIL_PORT"`
	FeedbackTo                 string        `envconfig:"FEEDBACK_TO"`
	FeedbackFrom               string        `envconfig:"FEEDBACK_FROM"`
	RendererURL                string        `envconfig:"RENDERER_URL"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg := &Config{
		BindAddr:                   "localhost:25200",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		HelloWorldEmphasise:        true,
		MailHost:                   "localhost",
		MailPort:                   "1025",
		MailUser:                   "",
		MailPassword:               "",
		FeedbackTo:                 "",
		FeedbackFrom:               "",
	}

	return cfg, envconfig.Process("", cfg)
}
