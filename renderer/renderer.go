package renderer

import (
	htmlrenderer "github.com/usememos/gomark/renderer/html"
	stringrenderer "github.com/usememos/gomark/renderer/string"
)

func NewHTMLRenderer() *htmlrenderer.HTMLRenderer {
	return htmlrenderer.NewHTMLRenderer()
}

func NewStringRenderer() *stringrenderer.StringRenderer {
	return stringrenderer.NewStringRenderer()
}
