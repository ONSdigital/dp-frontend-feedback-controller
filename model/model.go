package model

import "github.com/ONSdigital/dp-renderer/model"

// Page contains data reused for feedback model
type Feedback struct {
	model.Page
	Radio              string `json:"radio"`
	IsRadioErr         bool   `json:"is_radio_err"`
	Feedback           string `json:"feedback"`
	IsFeedbackErr      bool   `json:"is_feedback_err"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	IsEmailErr         bool   `json:"is_email_err"`
	PreviousURL        string `json:"previous_url"`
	IsURLErr           bool   `json:"is_url_err"`
	ReturnTo           string `json:"return_to"`
	ServiceDescription string `json:"service_description"`
}
