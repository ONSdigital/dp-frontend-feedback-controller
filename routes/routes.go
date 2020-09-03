package routes

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/renderer"
	"net/smtp"

	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/handlers"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, hc health.HealthCheck) {

	auth := smtp.PlainAuth(
		"",
		cfg.MailUser,
		cfg.MailPassword,
		cfg.MailHost,
	)

	if cfg.MailHost == "localhost" {
		auth = unencryptedAuth{auth}
	}

	renderer := renderer.New(cfg.RendererURL)

	mailAddr := fmt.Sprintf("%s:%s", cfg.MailHost, cfg.MailPort)

	log.Event(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	r.StrictSlash(true).Path("/feedback").Methods("POST").HandlerFunc(handlers.AddFeedback(auth, mailAddr, cfg.FeedbackTo, cfg.FeedbackFrom, false, renderer))
	r.StrictSlash(true).Path("/feedback/positive").Methods("POST").HandlerFunc(handlers.AddFeedback(auth, mailAddr, cfg.FeedbackTo, cfg.FeedbackFrom, false, renderer))
	r.StrictSlash(true).Path("/feedback").Methods("GET").HandlerFunc(handlers.GetFeedback(renderer))
	r.StrictSlash(true).Path("/feedback/thanks").Methods("GET").HandlerFunc(handlers.FeedbackThanks(renderer))

}

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}
