package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-frontend-feedback-controller/features/steps"
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
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
	})
}

func TestComponent(t *testing.T) {
	if *componentFlag {

		status := 0

		opts := godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  []string{"features/feedback"},
		}

		status = godog.TestSuite{
			Name:                 "feedback_tests",
			ScenarioInitializer:  InitializeScenario,
			TestSuiteInitializer: InitializeTestSuite,
			Options:              &opts,
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
