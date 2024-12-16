package steps

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-frontend-feedback-controller/service"
	"github.com/ONSdigital/dp-frontend-feedback-controller/service/mocks"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/cucumber/godog"
)

// HealthCheckTest represents a test healthcheck struct that mimics the real healthcheck struct
type HealthCheckTest struct {
	Status    string                  `json:"status"`
	Version   healthcheck.VersionInfo `json:"version"`
	Uptime    time.Duration           `json:"uptime"`
	StartTime time.Time               `json:"start_time"`
	Checks    []*Check                `json:"checks"`
}

type Check struct {
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	StatusCode  int        `json:"status_code"`
	Message     string     `json:"message"`
	LastChecked *time.Time `json:"last_checked"`
	LastSuccess *time.Time `json:"last_success"`
	LastFailure *time.Time `json:"last_failure"`
}

func (c *FeedbackComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.uiFeature.RegisterSteps(ctx)

	ctx.Step(`^the feedback controller is running$`, c.theFeedbackControllerIsRunning)
	ctx.Step(`^I wait (\d+) seconds`, c.delayTimeBySeconds)
}

func (c *FeedbackComponent) theFeedbackControllerIsRunning() error {
	ctx := context.Background()

	initFunctions := &mocks.InitialiserMock{
		DoGetHTTPServerFunc:   c.getHTTPServer,
		DoGetHealthCheckFunc:  getHealthCheckOK,
		DoGetHealthClientFunc: c.getHealthClient,
	}

	serviceList := service.NewServiceList(initFunctions)

	c.svc = service.New()
	if err := c.svc.Init(ctx, c.Config, serviceList); err != nil {
		log.Error(err)
		return err
	}

	svcErrors := make(chan error, 1)

	c.StartTime = time.Now()
	c.svc.Run(ctx, c.svc.Config, svcErrors)
	c.ServiceRunning = true

	return nil
}

func (c *FeedbackComponent) delayTimeBySeconds(sec int) error {
	time.Sleep(time.Duration(int64(sec)) * time.Second)
	return nil
}