package model

import "github.com/ONSdigital/dp-renderer/model"

// Page contains data reused for feedback model
type Feedback struct {
	model.Page
	Radio              string `json:"radio"`
	Feedback           string `json:"feedback"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	ErrorType          string `json:"error_type"`
	PreviousURL        string `json:"previous_url"`
	ReturnTo           string `json:"return_to"`
	ServiceDescription string `json:"service_description"`
}
