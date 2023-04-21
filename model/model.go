package model

import "github.com/ONSdigital/dp-renderer/model"

// Page contains data reused for feedback model
type Feedback struct {
	model.Page
	Contact          []model.TextField   `json:"contact"`
	TypeRadios       model.RadioFieldset `json:"type_radios"`
	DescriptionField model.TextareaField `json:"description_field"`
	PreviousURL      string              `json:"previous_url"`
	ReturnTo         string              `json:"return_to"`
}
