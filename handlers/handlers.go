package handlers

import (
	"net/http"

	cacheHelper "github.com/ONSdigital/dp-frontend-cache-helper/pkg/navigation/helper"
	"github.com/ONSdigital/dp-frontend-feedback-controller/config"
	"github.com/ONSdigital/dp-frontend-feedback-controller/email"
	"github.com/ONSdigital/dp-frontend-feedback-controller/interfaces"
	"github.com/ONSdigital/log.go/v2/log"
)

// Feedback represents the handlers required to provide feedback
type Feedback struct {
	Render       interfaces.Renderer
	CacheService *cacheHelper.Helper
	Config       *config.Config
	EmailSender  email.Sender
}

// NewFeedback creates a new instance of Feedback
func NewFeedback(rc interfaces.Renderer, c *cacheHelper.Helper, cfg *config.Config, es email.Sender) *Feedback {
	return &Feedback{
		Render:       rc,
		CacheService: c,
		Config:       cfg,
		EmailSender:  es,
	}
}

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Error(req.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}
