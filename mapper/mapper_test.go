package mapper

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-frontend-feedback-controller/config"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	ctx := context.Background()

	Convey("test CreateFilterOverview correctly maps item to filterOverview page model", t, func() {
		cfg := config.Config{
			BindAddr:                   "1234",
			GracefulShutdownTimeout:    0,
			HealthCheckInterval:        0,
			HealthCheckCriticalTimeout: 0,
		}

		hm := HelloModel{
			Greeting: "Hello",
			Who:      "World",
		}

		hw := HelloWorld(ctx, hm, cfg)
		So(hw, ShouldResemble, HelloWorldModel{"Hello World!"})
	})
}
