package routes

import (
	"context"
	"github.com/ONSdigital/dp-feedback-api/sdk"
	"net/smtp"

	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/handlers"

	render "github.com/ONSdigital/dp-renderer"

	cacheHelper "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/helper"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, rend *render.Render, hc health.HealthCheck, cacheService *cacheHelper.Helper) {
	auth := smtp.PlainAuth(
		"",
		cfg.MailUser,
		cfg.MailPassword,
		cfg.MailHost,
	)
	if cfg.MailHost == "localhost" {
		auth = unencryptedAuth{auth}
	}

	feedbackCfg := &config.FeedbackConfig{
		BindAddr:         "localhost:28600",
		ServiceAuthToken: "beehai7aeFoh4re8HaepaiFiwae9UXa6eeteimeil0ieyooyo5HohVoos2ahfeuw",
	}
	feedbackCfg.Client = sdk.New(feedbackCfg.BindAddr)

	log.Info(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	r.StrictSlash(true).Path("/feedback").Methods("GET").HandlerFunc(handlers.GetFeedback(rend, cacheService))
	r.StrictSlash(true).Path("/feedback").Methods("POST").HandlerFunc(handlers.AddFeedback(cfg.FeedbackTo, cfg.FeedbackFrom, false, rend, cacheService, feedbackCfg))
	r.StrictSlash(true).Path("/feedback/thanks").Methods("GET").HandlerFunc(handlers.FeedbackThanks(rend, cacheService))
	r.StrictSlash(true).Path("/feedback/thanks").Methods("POST").HandlerFunc(handlers.AddFeedback(cfg.FeedbackTo, cfg.FeedbackFrom, false, rend, cacheService, feedbackCfg))
}

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (proto string, toServer []byte, err error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}
