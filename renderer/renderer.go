// Package renderer provides adapters for converting AST documents to various formats.
package renderer

import (
	"github.com/usememos/gomark/ast"
	htmlrenderer "github.com/usememos/gomark/renderer/html"
	markdownrenderer "github.com/usememos/gomark/renderer/markdown"
	stringrenderer "github.com/usememos/gomark/renderer/string"
)

// Renderer defines the common interface for all renderers.
type Renderer interface {
	// RenderDocument renders the entire document and returns the result as a string.
	RenderDocument(doc *ast.Document) string

	// RenderNode renders a single node (for advanced use cases).
	RenderNode(node ast.Node)

	// String returns the accumulated output as a string.
	String() string

	// Reset clears the internal buffer to allow reuse.
	Reset()
}

// RendererType defines the available renderer types.
type RendererType string

const (
	HTMLRenderer     RendererType = "html"
	MarkdownRenderer RendererType = "markdown"
	StringRenderer   RendererType = "string"
)

// RendererFactory creates renderers by type.
type RendererFactory struct{}

// NewRendererFactory creates a new renderer factory.
func NewRendererFactory() *RendererFactory {
	return &RendererFactory{}
}

// CreateRenderer creates a renderer of the specified type.
func (*RendererFactory) CreateRenderer(rendererType RendererType) Renderer {
	switch rendererType {
	case HTMLRenderer:
		return NewHTMLRenderer()
	case MarkdownRenderer:
		return NewMarkdownRenderer()
	case StringRenderer:
		return NewStringRenderer()
	default:
		return NewMarkdownRenderer() // Default fallback
	}
}

// NewHTMLRenderer creates an HTML renderer that implements the Renderer interface.
func NewHTMLRenderer() Renderer {
	return &HTMLRendererAdapter{htmlrenderer.NewHTMLRenderer()}
}

// NewStringRenderer creates a string renderer that implements the Renderer interface.
func NewStringRenderer() Renderer {
	return &StringRendererAdapter{stringrenderer.NewStringRenderer()}
}

// NewMarkdownRenderer creates a markdown renderer that implements the Renderer interface.
func NewMarkdownRenderer() Renderer {
	return &MarkdownRendererAdapter{markdownrenderer.NewMarkdownRenderer()}
}

// HTMLRendererAdapter adapts the concrete HTML renderer to the Renderer interface.
type HTMLRendererAdapter struct {
	renderer *htmlrenderer.HTMLRenderer
}

func (a *HTMLRendererAdapter) RenderDocument(doc *ast.Document) string {
	return a.renderer.RenderDocument(doc)
}

func (a *HTMLRendererAdapter) RenderNode(node ast.Node) {
	a.renderer.RenderNode(node)
}

func (a *HTMLRendererAdapter) String() string {
	// HTML renderer doesn't expose its output buffer directly
	// so we need to render a temporary empty document to get current state
	return a.renderer.RenderDocument(&ast.Document{})
}

func (a *HTMLRendererAdapter) Reset() {
	a.renderer = htmlrenderer.NewHTMLRenderer()
}

// StringRendererAdapter adapts the concrete String renderer to the Renderer interface.
type StringRendererAdapter struct {
	renderer *stringrenderer.StringRenderer
}

func (a *StringRendererAdapter) RenderDocument(doc *ast.Document) string {
	return a.renderer.RenderDocument(doc)
}

func (a *StringRendererAdapter) RenderNode(node ast.Node) {
	a.renderer.RenderNode(node)
}

func (a *StringRendererAdapter) String() string {
	// String renderer doesn't expose its output buffer directly
	// so we need to render a temporary empty document to get current state
	return a.renderer.RenderDocument(&ast.Document{})
}

func (a *StringRendererAdapter) Reset() {
	a.renderer = stringrenderer.NewStringRenderer()
}

// MarkdownRendererAdapter adapts the concrete Markdown renderer to the Renderer interface.
type MarkdownRendererAdapter struct {
	renderer *markdownrenderer.MarkdownRenderer
}

func (a *MarkdownRendererAdapter) RenderDocument(doc *ast.Document) string {
	a.renderer.RenderDocument(doc)
	return a.renderer.String()
}

func (a *MarkdownRendererAdapter) RenderNode(node ast.Node) {
	// MarkdownRenderer doesn't have RenderNode, so we render it directly
	// This is a limitation of the current markdown renderer design
	if doc, ok := node.(*ast.Document); ok {
		a.renderer.RenderDocument(doc)
	} else {
		// For individual nodes, we create a temporary document
		tempDoc := &ast.Document{Children: []ast.Node{node}}
		a.renderer.RenderDocument(tempDoc)
	}
}

func (a *MarkdownRendererAdapter) String() string {
	return a.renderer.String()
}

func (a *MarkdownRendererAdapter) Reset() {
	a.renderer = markdownrenderer.NewMarkdownRenderer()
}
