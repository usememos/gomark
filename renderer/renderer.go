package renderer

import (
	htmlrenderer "github.com/yourselfhosted/gomark/renderer/html"
	stringrenderer "github.com/yourselfhosted/gomark/renderer/string"
)

func NewHTMLRenderer() *htmlrenderer.HTMLRenderer {
	return htmlrenderer.NewHTMLRenderer()
}

func NewStringRenderer() *stringrenderer.StringRenderer {
	return stringrenderer.NewStringRenderer()
}
