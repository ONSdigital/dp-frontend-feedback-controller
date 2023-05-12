package model

import "github.com/ONSdigital/dp-renderer/v2/model"

// Page contains data reused for feedback model
type Feedback struct {
	model.Page
	Contact          []model.TextField   `json:"contact"`
	TypeRadios       model.RadioFieldset `json:"type_radios"`
	DescriptionField model.TextareaField `json:"description_field"`
	PreviousURL      string              `json:"previous_url"`
	ReturnTo         string              `json:"return_to"`
}

// FeedbackForm represents the user feedback form
type FeedbackForm struct {
	FormLocation     string `schema:"feedback-form-type"`
	Type             string `schema:"type"`
	IsTypeErr        bool   `schema:"is_type_err"`
	URI              string `schema:":uri"`
	URL              string `schema:"url"`
	IsURLErr         bool   `schema:"is_url_err"`
	Description      string `schema:"description"`
	IsDescriptionErr bool   `schema:"is_description_err"`
	Name             string `schema:"name"`
	Email            string `schema:"email"`
	IsEmailErr       bool   `schema:"is_email_err"`
}
