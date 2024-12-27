package steps

import (
	"context"
		"time"

	"github.com/ONSdigital/dp-frontend-feedback-controller/service"
	"github.com/ONSdigital/dp-frontend-feedback-controller/service/mocks"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/chromedp/chromedp"
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
	ctx.Step(`^the feedback controller is running$`, c.theFeedbackControllerIsRunning)
	ctx.Step(`^I wait (\d+) seconds`, c.delayTimeBySeconds)
	ctx.Step(`^element "([^"]*)" should be visible$`, c.elementShouldBeVisible)
	ctx.Step(`^I fill in input element "([^"]*)" with value "([^"]*)"$`, c.iFillInInputElementWithValue)
	ctx.Step(`^I click the "([^"]*)" element$`, c.iClickElement)
}

func (c *FeedbackComponent) theFeedbackControllerIsRunning() error {
	ctx := context.Background()

	initFunctions := &mocks.InitialiserMock{
		DoGetHTTPServerFunc:   c.getHTTPServer,
		DoGetHealthCheckFunc:  c.getHealthCheckOK,
		DoGetHealthClientFunc: c.getHealthClient,
	}

	serviceList := service.NewServiceList(initFunctions)

	c.svc = service.New()
	if err := c.svc.Init(ctx, c.Config, serviceList); err != nil {
		log.Error(ctx, "failed to init service", err)
		return err
	}

	svcErrors := make(chan error, 1)

	c.StartTime = time.Now()
	c.svc.Run(ctx, svcErrors)
	c.ServiceRunning = true

	return nil
}

func (c *FeedbackComponent) delayTimeBySeconds(sec int) error {
	time.Sleep(time.Duration(int64(sec)) * time.Second)
	return nil
}

func (c *FeedbackComponent) iFillInInputElementWithValue(fieldSelector, value string) error {
	jsScript := fmt.Sprintf(`document.querySelector('%s').value = '%s';`, fieldSelector, value)

	err := chromedp.Run(c.Chrome.Ctx,
		chromedp.Evaluate(jsScript, nil),
	)
	if err != nil {
		return err
	}

	return c.ErrorFeature.StepError()
}

func (c *FeedbackComponent) iClickElement(buttonSelector string) error {
	// if this doesn't work as expected, you might need a sleep after the click
	err := chromedp.Run(c.Chrome.Ctx,
		chromedp.Click(buttonSelector),
	)
	if err != nil {
		return err
	}

	return c.ErrorFeature.StepError()
}

func (c *FeedbackComponent) RunWithTimeOut(timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}

func (c *FeedbackComponent) elementShouldBeVisible(elementSelector string) error {
	err := chromedp.Run(c.Chrome.Ctx,
		c.RunWithTimeOut(c.WaitTimeOut, chromedp.Tasks{
			chromedp.WaitVisible(elementSelector),
		}),
	)
	assert.Nil(c, err)

	return c.ErrorFeature.StepError()
}
