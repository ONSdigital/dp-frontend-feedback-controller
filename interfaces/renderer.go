package interfaces

import (
	"io"

	coreModel "github.com/ONSdigital/dp-renderer/model"
)

//go:generate moq -out interfacestest/renderer.go -pkg interfacestest . Renderer

type Renderer interface {
	BuildPage(w io.Writer, pageModel interface{}, templateName string)
	NewBasePageModel() coreModel.Page
}
