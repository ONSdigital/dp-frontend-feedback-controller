package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-feedback-controller
type Config struct {
	APIRouterURL                string         `envconfig:"API_ROUTER_URL"`
	BindAddr                    string         `envconfig:"BIND_ADDR"`
	CacheUpdateInterval         *time.Duration `envconfig:"CACHE_UPDATE_INTERVAL"`
	CensusTopicID               string         `envconfig:"CENSUS_TOPIC_ID"`
	Debug                       bool           `envconfig:"DEBUG"`
	EnableCensusTopicSubsection bool           `envconfig:"ENABLE_CENSUS_TOPIC_SUBSECTION"`
	EnableNewNavBar             bool           `envconfig:"ENABLE_NEW_NAVBAR"`
	GracefulShutdownTimeout     time.Duration  `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval         time.Duration  `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout  time.Duration  `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	IsPublishing                bool           `envconfig:"IS_PUBLISHING"`
	PatternLibraryAssetsPath    string         `envconfig:"PATTERN_LIBRARY_ASSETS_PATH"`
	ServiceAuthToken            string         `envconfig:"SERVICE_AUTH_TOKEN"   json:"-"`
	SiteDomain                  string         `envconfig:"SITE_DOMAIN"`
	SupportedLanguages          []string       `envconfig:"SUPPORTED_LANGUAGES"`
	OTExporterOTLPEndpoint      string         `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTServiceName               string         `envconfig:"OTEL_SERVICE_NAME"`
	OTBatchTimeout              time.Duration  `envconfig:"OTEL_BATCH_TIMEOUT"`
	OtelEnabled                 bool           `envconfig:"OTEL_ENABLED"`
}

var cfg *Config
var RendererVersion = "v0.2.0" // default value, overridden at build

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	envCfg, err := get()
	if err != nil {
		return nil, err
	}

	if envCfg.Debug {
		envCfg.PatternLibraryAssetsPath = "http://localhost:9002/dist/assets"
	} else {
		envCfg.PatternLibraryAssetsPath = fmt.Sprintf("//cdn.ons.gov.uk/dis-design-system-go/%s", RendererVersion)
	}

	cfg = envCfg
	return cfg, nil
}

func get() (*Config, error) {
	cfg := &Config{
		APIRouterURL:                "http://localhost:23200/v1",
		BindAddr:                    ":25200",
		CensusTopicID:               "4445",
		Debug:                       false,
		EnableCensusTopicSubsection: false,
		EnableNewNavBar:             false,
		GracefulShutdownTimeout:     5 * time.Second,
		HealthCheckInterval:         30 * time.Second,
		HealthCheckCriticalTimeout:  90 * time.Second,
		IsPublishing:                false,
		ServiceAuthToken:            "",
		SiteDomain:                  "localhost",
		SupportedLanguages:          []string{"en", "cy"},
		OTExporterOTLPEndpoint:      "localhost:4317",
		OTServiceName:               "dp-frontend-feedback-controller",
		OTBatchTimeout:              5 * time.Second,
		OtelEnabled:                 false,
	}

	return cfg, envconfig.Process("", cfg)
}
