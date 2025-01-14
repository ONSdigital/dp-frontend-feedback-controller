package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out mocks/initialiser.go -pkg mocks . Initialiser
//go:generate moq -out mocks/healthcheck.go -pkg mocks . HealthChecker
//go:generate moq -out mocks/server.go -pkg mocks . HTTPServer

// Initialiser defines the methods to initialise external services
type Initialiser interface {
	DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer
	DoGetHealthClient(name, url string) *health.Client
	DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error)
}

// HealthChecker defines the required methods from Healthcheck
type HealthChecker interface {
	AddCheck(name string, checker healthcheck.Checker) (err error)
	Handler(w http.ResponseWriter, req *http.Request)
	Start(ctx context.Context)
	Stop()
}

// HTTPServer defines the required methods from the HTTP server
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
