package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	golog "log"
	"os"
	"testing"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-frontend-feedback-controller/features/steps"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var componentFlag = flag.Bool("component", false, "perform component tests")

func InitializeScenario(ctx *godog.ScenarioContext) {
	component, err := steps.NewFeedbackComponent()
	if err != nil {
		fmt.Printf("failed to create feedback controller component - error: %v", err)
		os.Exit(1)
	}

	url := fmt.Sprintf("http://%s%s", component.Config.SiteDomain, component.Config.BindAddr)

	uiFeature := componenttest.NewUIFeature(url)

	uiFeature.RegisterSteps(ctx)

	component.RegisterSteps(ctx)

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		uiFeature.Reset()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		uiFeature.Close()
		err = component.Close()
		return ctx, err
	})
}

func TestComponent(t *testing.T) {
	if *componentFlag {
		log.SetDestination(io.Discard, io.Discard)
		golog.SetOutput(io.Discard)
		defer func() {
			log.SetDestination(os.Stdout, os.Stderr)
			golog.SetOutput(os.Stdout)
		}()

		status := 0

		opts := godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  []string{"features/feedback"},
		}

		status = godog.TestSuite{
			Name:                "feedback_tests",
			ScenarioInitializer: InitializeScenario,
			Options:             &opts,
		}.Run()

		fmt.Println("=================================")
		fmt.Printf("Component test coverage: %.2f%%\n", testing.Coverage()*100)
		fmt.Println("=================================")

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
