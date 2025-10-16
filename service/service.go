package service

import (
	"context"
	"errors"

	render "github.com/ONSdigital/dis-design-system-go"
	"github.com/ONSdigital/dis-design-system-go/middleware/renderror"
	feedbackAPI "github.com/ONSdigital/dp-feedback-api/sdk"
	cacheHelper "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/helper"
	"github.com/ONSdigital/dp-frontend-feedback-controller/assets"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/routes"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

// Service contains the healthcheck, server and serviceList for the controller
type Service struct {
	Config      *config.Config
	HealthCheck HealthChecker
	Server      HTTPServer
	ServiceList *ExternalServiceList
}

// New creates a new service
func New() *Service {
	return &Service{}
}

// Init initialises all the service dependencies, including healthcheck with checkers, api and middleware
func (svc *Service) Init(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList) (err error) {
	log.Info(ctx, "initialising service")

	svc.Config = cfg
	svc.ServiceList = serviceList

	// Get health client for api router
	routerHealthClient := serviceList.GetHealthClient("api-router", cfg.APIRouterURL)

	// Initialise clients
	clients := routes.Clients{
		Renderer:    render.NewWithDefaultClient(assets.Asset, assets.AssetNames, cfg.PatternLibraryAssetsPath, cfg.SiteDomain),
		FeedbackAPI: feedbackAPI.NewWithHealthClient(routerHealthClient),
	}

	// Get healthcheck with checkers
	svc.HealthCheck, err = serviceList.GetHealthCheck(cfg, BuildTime, GitCommit, Version)
	if err != nil {
		log.Fatal(ctx, "failed to create health check", err)
		return err
	}
	if err = svc.registerCheckers(ctx, clients); err != nil {
		log.Error(ctx, "failed to register checkers", err)
		return err
	}
	clients.HealthCheckHandler = svc.HealthCheck.Handler

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

	// Create service initialiser and an error channel for service errors
	svcErrors := make(chan error)

	cacheService, _ := cacheHelper.Init(ctx, cacheConfig)
	cacheService.RunUpdates(ctx, svcErrors)

	// Initialise router
	r := mux.NewRouter()
	if cfg.OtelEnabled {
		r.Use(otelmux.Middleware(cfg.OTServiceName))
	}
	middleware := []alice.Constructor{
		renderror.Handler(clients.Renderer),
	}
	newAlice := alice.New(middleware...).Then(r)
	routes.Setup(ctx, r, cfg, clients, cacheService)
	svc.Server = serviceList.GetHTTPServer(cfg.BindAddr, newAlice)

	return nil
}

// Run starts an initialised service
func (svc *Service) Run(ctx context.Context, svcErrors chan error) {
	log.Info(ctx, "Starting service", log.Data{"config": svc.Config})

	// Start healthcheck
	svc.HealthCheck.Start(ctx)

	// Start HTTP server
	log.Info(ctx, "Starting server")
	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			log.Fatal(ctx, "failed to start http listen and serve", err)
			svcErrors <- err
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	log.Info(ctx, "commencing graceful shutdown")
	ctx, cancel := context.WithTimeout(ctx, svc.Config.GracefulShutdownTimeout)
	hasShutdownError := false

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		log.Info(ctx, "stop health checkers")
		svc.HealthCheck.Stop()

		// stop any incoming requests
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func (svc *Service) registerCheckers(ctx context.Context, c routes.Clients) (err error) {
	hasErrors := false

	if err = svc.HealthCheck.AddCheck("Feedback API", c.FeedbackAPI.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "failed to add feedback API checker", err)
	}

	if hasErrors {
		return errors.New("error(s) registering checkers for healthcheck")
	}

	return nil
}
