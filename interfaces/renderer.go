package interfaces

//TODO: Refactor this package to ensure less generic names are used

import (
	"io"

	coreModel "github.com/ONSdigital/dis-design-system-go/model"
)

//go:generate moq -out interfacestest/renderer.go -pkg interfacestest . Renderer

type Renderer interface {
	BuildPage(w io.Writer, pageModel interface{}, templateName string)
	NewBasePageModel() coreModel.Page
}
