package handlers

import (
	"io"

	"github.com/ONSdigital/dis-design-system-go/model"
)

//go:generate moq -out clients_mock.go -pkg handlers . ClientError RenderClient

// RenderClient interface defines page rendering
type RenderClient interface {
	BuildPage(w io.Writer, pageModel interface{}, templateName string)
	NewBasePageModel() model.Page
}
