package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	cacheHelper "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/helper"
	"github.com/ONSdigital/dp-frontend-feedback-controller/assets"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/routes"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	dpotelgo "github.com/ONSdigital/dp-otel-go"
	render "github.com/ONSdigital/dp-renderer/v2"
	"github.com/ONSdigital/dp-renderer/v2/middleware/renderror"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	log.Namespace = "dp-frontend-feedback-controller"
	cfg, err := config.Get()
	ctx := context.Background()
	if err != nil {
		log.Error(ctx, "unable to retrieve service configuration", err)
		os.Exit(1)
	}

	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	if cfg.OtelEnabled {
		// Set up OpenTelemetry
		otelConfig := dpotelgo.Config{
			OtelServiceName:          cfg.OTServiceName,
			OtelExporterOtlpEndpoint: cfg.OTExporterOTLPEndpoint,
			OtelBatchTimeout:         cfg.OTBatchTimeout,
		}

		otelShutdown, err := dpotelgo.SetupOTelSDK(ctx, otelConfig)

		if err != nil {
			log.Error(ctx, "error setting up OpenTelemetry - hint: ensure OTEL_EXPORTER_OTLP_ENDPOINT is set", err)
		}

		// Handle shutdown properly so nothing leaks.
		defer func() {
			err = errors.Join(err, otelShutdown(context.Background()))
		}()
	}

	versionInfo, healthErr := health.NewVersionInfo(
		BuildTime,
		GitCommit,
		Version,
	)
	if healthErr != nil {
		log.Error(ctx, "failed to retrieve health check version", healthErr)
		return
	}

	r := mux.NewRouter()

	if cfg.OtelEnabled {
		r.Use(otelmux.Middleware(cfg.OTServiceName))
	}

	healthcheck := health.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)

	cacheConfig := cacheHelper.Config{
		APIRouterURL:                cfg.APIRouterURL,
		CacheUpdateInterval:         cfg.CacheUpdateInterval,
		EnableNewNavBar:             cfg.EnableNewNavBar,
		EnableCensusTopicSubsection: cfg.EnableCensusTopicSubsection,
		CensusTopicID:               cfg.CensusTopicID,
		IsPublishingMode:            cfg.IsPublishing,
		Languages:                   cfg.SupportedLanguages,
		ServiceAuthToken:            cfg.ServiceAuthToken,
	}

	svcErrors := make(chan error)
	cacheService, _ := cacheHelper.Init(ctx, cacheConfig)
	cacheService.RunUpdates(ctx, svcErrors)

	//nolint:typecheck // ignore typecheck as make command moves assets
	rend := render.NewWithDefaultClient(assets.Asset, assets.AssetNames, cfg.PatternLibraryAssetsPath, cfg.SiteDomain)

	middleware := []alice.Constructor{
		renderror.Handler(rend),
	}
	newAlice := alice.New(middleware...).Then(r)

	routes.Setup(ctx, r, cfg, rend, healthcheck, cacheService)

	healthcheck.Start(ctx)

	var s *dphttp.Server

	if cfg.OtelEnabled {
		otelHandler := otelhttp.NewHandler(newAlice, "/")
		s = dphttp.NewServer(cfg.BindAddr, otelHandler)
	} else {
		s = dphttp.NewServer(cfg.BindAddr, newAlice)
	}

	s.HandleOSSignals = false

	log.Info(ctx, "Starting server", log.Data{"config": cfg})

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Error(ctx, "failed to start http listen and serve", err)
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Info(ctx, "shutting service down gracefully")
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Error(ctx, "failed to shutdown http server", err)
	}
}
