package steps

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/service"
	"github.com/ONSdigital/dp-frontend-feedback-controller/service/mocks"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	gitCommitHash = "19c075791fc985b380771404528b8dce28e65873"
	appVersion    = "v1.7.0"
)

type FeedbackComponent struct {
	Chrome         componenttest.Chrome
	Config         *config.Config
	errorChan      chan error
	ErrorFeature   componenttest.ErrorFeature
	FakeAPIRouter  *FakeAPI
	HTTPServer     *http.Server
	ServiceRunning bool
	StartTime      time.Time
	svc            *service.Service
	svcList        *service.ExternalServiceList
	uiFeature      componenttest.UIFeature
	WaitTimeOut    time.Duration
}

// Errorf implements assert.TestingT.
func (c *FeedbackComponent) Errorf(format string, args ...interface{}) {
	panic("unimplemented")
}

func NewFeedbackComponent() (*FeedbackComponent, error) {
	c := &FeedbackComponent{
		errorChan: make(chan error, 1),
		HTTPServer: &http.Server{
			ReadHeaderTimeout: 60 * time.Second,
		},
		ServiceRunning: false,
	}

	ctx := context.Background()

	log.Info(ctx, "configuration for component test", log.Data{"config": c.Config})

	var err error

	c.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	initMock := &mocks.InitialiserMock{
		DoGetHealthCheckFunc: c.DoGetHealthcheckOk,
		DoGetHTTPServerFunc:  c.DoGetHTTPServer,
	}

	c.svcList = service.NewServiceList(initMock)

	return c, nil
}

// Close server running component.
func (c *FeedbackComponent) Close() error {
	if c.ServiceRunning {
		err := c.close(context.Background())
		c.ServiceRunning = false
		return err
	}
	return nil
}

// Reset resets the component. Used to reset the component between tests.
func (c *FeedbackComponent) Reset() *FeedbackComponent {
	c.uiFeature.Reset()
	return c
}

func (c *FeedbackComponent) close(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.GracefulShutdownTimeout)
	hasShutdownError := false

	go func() {
		defer cancel()

		// stop any incoming requests
		if err := c.HTTPServer.Shutdown(ctx); err != nil {
			hasShutdownError = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		return err
	}

	return nil
}

func (c *FeedbackComponent) getHealthCheckOK(cfg *config.Config, _, _, _ string) (service.HealthChecker, error) {
	componentBuildTime := strconv.Itoa(int(time.Now().Unix()))
	versionInfo, err := healthcheck.NewVersionInfo(componentBuildTime, gitCommitHash, appVersion)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

func (c *FeedbackComponent) DoGetHealthcheckOk(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

func (c *FeedbackComponent) getHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}

func (c *FeedbackComponent) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
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
