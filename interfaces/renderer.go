package interfaces

//go:generate moq -out interfacestest/renderer.go -pkg interfacestest . Renderer

type Renderer interface {
	Do(path string, b []byte) ([]byte, error)
}
