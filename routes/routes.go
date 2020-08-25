package routes

import (
	"context"

	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/handlers"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, hc health.HealthCheck) {
	log.Event(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	r.StrictSlash(true).Path("/feedback").Methods("POST").HandlerFunc(handlers.AddFeedback(auth, mailAddr, cfg.FeedbackTo, cfg.FeedbackFrom, cfg.RendererURL, false))
	r.StrictSlash(true).Path("/feedback/positive").Methods("POST").HandlerFunc(handlers.AddFeedback(auth, mailAddr, cfg.FeedbackTo, cfg.FeedbackFrom, cfg.RendererURL, false))
	r.StrictSlash(true).Path("/feedback").Methods("GET").HandlerFunc(handlers.GetFeedback(cfg.RendererURL))
	r.StrictSlash(true).Path("/feedback/thanks").Methods("GET").HandlerFunc(handlers.FeedbackThanks(cfg.RendererURL))
}
