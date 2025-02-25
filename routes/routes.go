package routes

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/handlers"

	feedbackAPI "github.com/ONSdigital/dp-feedback-api/sdk"

	render "github.com/ONSdigital/dp-renderer/v2"

	cacheHelper "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/helper"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	HealthCheckHandler func(w http.ResponseWriter, req *http.Request)
	Renderer           *render.Render
	FeedbackAPI        *feedbackAPI.Client
}

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, c Clients, cacheService *cacheHelper.Helper) {
	f := handlers.NewFeedback(c.Renderer, cacheService, cfg)

	log.Info(ctx, "adding routes")
	r.StrictSlash(true).Path("/health").HandlerFunc(c.HealthCheckHandler)
	r.StrictSlash(true).Path("/feedback").Methods("GET").HandlerFunc(f.GetFeedback())
	r.StrictSlash(true).Path("/feedback").Methods("POST").HandlerFunc(f.AddFeedback())
	r.StrictSlash(true).Path("/feedback/thanks").Methods("GET").HandlerFunc(f.FeedbackThanks())
	r.StrictSlash(true).Path("/feedback/thanks").Methods("POST").HandlerFunc(f.AddFeedback())
}
