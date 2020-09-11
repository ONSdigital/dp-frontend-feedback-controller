package routes

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/renderer"
	"github.com/ONSdigital/dp-frontend-feedback-controller/email"
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
	mailAddr := fmt.Sprintf("%s:%s", cfg.MailHost, cfg.MailPort)

	emailSender := email.SMTPSender{
		Addr: mailAddr,
		Auth: auth,
	}

	renderer := renderer.New(cfg.RendererURL)

	log.Event(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	r.StrictSlash(true).Path("/feedback").Methods("POST").HandlerFunc(handlers.AddFeedback(cfg.FeedbackTo, cfg.FeedbackFrom, false, renderer, emailSender))
	r.StrictSlash(true).Path("/feedback/positive").Methods("POST").HandlerFunc(handlers.AddFeedback(cfg.FeedbackTo, cfg.FeedbackFrom, false, renderer, emailSender))
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
