package steps

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	componentTest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/service"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/maxcnunes/httpfake"
)

const (
	gitCommitHash = "19c075791fc985b380771404528b8dce28e65873"
	appVersion    = "v1.7.0"
)

type FeedbackComponent struct {
	APIFeature     *componentTest.APIFeature
	Config         *config.Config
	ErrorFeature   componentTest.ErrorFeature
	FakeAPIRouter  *FakeAPI
	HTTPServer     *http.Server
	ServiceRunning bool
	svc            *service.Service
	svcErrors      chan error
	StartTime      time.Time
}

func NewFeedbackComponent() (c *FeedbackComponent, err error) {
	c = &FeedbackComponent{
		HTTPServer: &http.Server{
			ReadHeaderTimeout: 5 * time.Second,
		},
		svcErrors: make(chan error),
	}

	ctx := context.Background()

	c.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	log.Info(ctx, "configuration for component test", log.Data{"config": c.Config})

	c.FakeAPIRouter = NewFakeAPI()
	c.Config.APIRouterURL = c.FakeAPIRouter.fakeHTTP.ResolveURL("")

	c.Config.HealthCheckInterval = 1 * time.Second
	c.Config.HealthCheckCriticalTimeout = 3 * time.Second

	c.FakeAPIRouter.healthRequest = c.FakeAPIRouter.fakeHTTP.NewHandler().Get("/health")
	c.FakeAPIRouter.healthRequest.CustomHandle = healthCheckStatusHandle(200)

	c.FakeAPIRouter.feedbackRequest = c.FakeAPIRouter.fakeHTTP.NewHandler().Get("/feedback")

	return c, nil
}

// InitAPIFeature initialises the ApiFeature
func (c *FeedbackComponent) InitAPIFeature() *componentTest.APIFeature {
	c.APIFeature = componentTest.NewAPIFeature(c.InitialiseService)

	return c.APIFeature
}

// Close server running component.
func (c *FeedbackComponent) Close() error {
	if c.svc != nil && c.ServiceRunning {
		c.svc.Close(context.Background())
		c.ServiceRunning = false
	}

	c.FakeAPIRouter.Close()

	return nil
}

// InitialiseService returns the http.Handler that's contained within the component.
func (c *FeedbackComponent) InitialiseService() (http.Handler, error) {
	return c.HTTPServer.Handler, nil
}

func getHealthCheckOK(cfg *config.Config, _, _, _ string) (service.HealthChecker, error) {
	componentBuildTime := strconv.Itoa(int(time.Now().Unix()))
	versionInfo, err := healthcheck.NewVersionInfo(componentBuildTime, gitCommitHash, appVersion)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

func (c *FeedbackComponent) getHealthClient(name, url string) *health.Client {
	return &health.Client{
		URL:    url,
		Name:   name,
		Client: c.FakeAPIRouter.getMockAPIHTTPClient(),
	}
}

// newMock mocks HTTP Client
func (f *FakeAPI) getMockAPIHTTPClient() *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(_ []string) {},
		GetPathsWithNoRetriesFunc: func() []string { return []string{} },
		DoFunc: func(_ context.Context, req *http.Request) (*http.Response, error) {
			return f.fakeHTTP.Server.Client().Do(req)
		},
	}
}

func (c *FeedbackComponent) getHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}

func healthCheckStatusHandle(status int) httpfake.Responder {
	return func(w http.ResponseWriter, _ *http.Request, rh *httpfake.Request) {
		rh.Lock()
		defer rh.Unlock()
		w.WriteHeader(status)
	}
}
