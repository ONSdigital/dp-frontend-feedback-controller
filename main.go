package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	cacheHelper "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/helper"
	"github.com/ONSdigital/dp-frontend-feedback-controller/assets"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/routes"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	render "github.com/ONSdigital/dp-renderer/v2"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
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

	healthcheck := health.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	if err = registerCheckers(ctx, &healthcheck); err != nil {
		os.Exit(1)
	}

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

	svcErrors := make(chan error, 0)
	cacheService, err := cacheHelper.Init(ctx, cacheConfig)
	cacheService.RunUpdates(ctx, svcErrors)

	//nolint:typecheck
	rend := render.NewWithDefaultClient(assets.Asset, assets.AssetNames, cfg.PatternLibraryAssetsPath, cfg.SiteDomain)

	routes.Setup(ctx, r, cfg, rend, healthcheck, cacheService)

	healthcheck.Start(ctx)

	s := server.New(cfg.BindAddr, r)
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

func registerCheckers(ctx context.Context, h *health.HealthCheck) (err error) {
	// TODO ADD HEALTH CHECKS HERE
	return
}
