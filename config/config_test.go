package config

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Given an environment with no environment variables set", t, func() {
		Convey("When the config values are retrieved", func() {
			cfg, err := Get()
			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the values should be set to the expected defaults", func() {
				So(cfg.EnableNewNavBar, ShouldEqual, false)
				So(cfg.GracefulShutdownTimeout, ShouldEqual, 5*time.Second)
				So(cfg.HealthCheckInterval, ShouldEqual, 30*time.Second)
				So(cfg.HealthCheckCriticalTimeout, ShouldEqual, 90*time.Second)
				So(cfg.MailHost, ShouldEqual, "localhost")
				So(cfg.MailPort, ShouldEqual, "1025")
				So(cfg.MailUser, ShouldEqual, "")
				So(cfg.MailPassword, ShouldEqual, "")
				So(cfg.FeedbackTo, ShouldEqual, "to@gmail.com")
				So(cfg.FeedbackFrom, ShouldEqual, "from@gmail.com")
				So(cfg.FeedbackAPIEnabled, ShouldEqual, false)
				So(cfg.SiteDomain, ShouldEqual, "localhost")
				So(cfg.Debug, ShouldEqual, false)
				So(cfg.SupportedLanguages, ShouldResemble, []string{"en", "cy"})
				So(cfg.IsPublishing, ShouldEqual, false)
				So(cfg.EnableCensusTopicSubsection, ShouldEqual, false)
				So(cfg.OTExporterOTLPEndpoint, ShouldEqual, "localhost:4317")
				So(cfg.OTServiceName, ShouldEqual, "dp-frontend-feedback-controller")
				So(cfg.OTBatchTimeout, ShouldEqual, 5*time.Second)
			})

			Convey("Then a second call to config should return the same config", func() {
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})
}
