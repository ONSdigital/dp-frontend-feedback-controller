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
	MailHost                   string        `envconfig:"MAIL_HOST"`
	MailUser                   string        `envconfig:"MAIL_USER"`
	MailPassword               string        `envconfig:"MAIL_PASSWORD" json:"-"`
	MailPort                   string        `envconfig:"MAIL_PORT"`
	FeedbackTo                 string        `envconfig:"FEEDBACK_TO"`
	FeedbackFrom               string        `envconfig:"FEEDBACK_FROM"`
	PatternLibraryAssetsPath   string        `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	SiteDomain                 string        `envconfig:"SITE_DOMAIN"`
	Debug                      bool          `envconfig:"DEBUG"`
	SupportedLanguages         []string      `envconfig:"SUPPORTED_LANGUAGES"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	envCfg, err := get()
	if err != nil {
		return nil, err
	}

	if envCfg.Debug {
		envCfg.PatternLibraryAssetsPath = "http://localhost:9000/dist"
	} else {
		envCfg.PatternLibraryAssetsPath = "//cdn.ons.gov.uk/sixteens/67f6982"
	}
	return envCfg, nil
}

func get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg := &Config{
		BindAddr:                   "localhost:25200",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		MailHost:                   "localhost",
		MailPort:                   "1025",
		MailUser:                   "",
		MailPassword:               "",
		FeedbackTo:                 "to@gmail.com",
		FeedbackFrom:               "from@gmail.com",
		SiteDomain:                 "localhost",
		Debug:                      false,
		SupportedLanguages:         []string{"en", "cy"},
	}

	return cfg, envconfig.Process("", cfg)
}
