package routes

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/ONSdigital/dp-frontend-feedback-controller/email"

	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/handlers"

	render "github.com/ONSdigital/dp-renderer/v2"

	cacheHelper "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/helper"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, rend *render.Render, hc health.HealthCheck, cacheService *cacheHelper.Helper) {
	var auth smtp.Auth
	if cfg.MailEncrypted {
		auth = smtp.PlainAuth(
			"",
			cfg.MailUser,
			cfg.MailPassword,
			cfg.MailHost,
		)
	} else {
		auth = smtp.CRAMMD5Auth(cfg.MailUser, cfg.MailPassword)
	}
	mailAddr := fmt.Sprintf("%s:%s", cfg.MailHost, cfg.MailPort)

	emailSender := email.SMTPSender{
		Addr: mailAddr,
		Auth: auth,
	}

	f := handlers.NewFeedback(rend, cacheService, cfg, emailSender)

	log.Info(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	r.StrictSlash(true).Path("/feedback").Methods("GET").HandlerFunc(f.GetFeedback())
	r.StrictSlash(true).Path("/feedback").Methods("POST").HandlerFunc(f.AddFeedback())
	r.StrictSlash(true).Path("/feedback/thanks").Methods("GET").HandlerFunc(f.FeedbackThanks())
	r.StrictSlash(true).Path("/feedback/thanks").Methods("POST").HandlerFunc(f.AddFeedback())
}

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (proto string, toServer []byte, err error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}
