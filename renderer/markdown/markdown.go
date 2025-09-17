package markdown

import (
	"bytes"

	"github.com/usememos/gomark/ast"
)

// MarkdownRenderer renders an AST back into markdown text.
type MarkdownRenderer struct {
	output *bytes.Buffer
}

// NewMarkdownRenderer constructs a markdown renderer.
func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{output: new(bytes.Buffer)}
}

// RenderDocument renders the full document tree.
func (r *MarkdownRenderer) RenderDocument(doc *ast.Document) {
	if doc == nil {
		return
	}
	if len(doc.Children) == 0 {
		r.renderNode(doc)
		return
	}
	r.RenderNodes(doc.Children)
}

// RenderNodes renders an ordered list of nodes.
func (r *MarkdownRenderer) RenderNodes(nodes []ast.Node) {
	for _, node := range nodes {
		r.renderNode(node)
	}
}

func (r *MarkdownRenderer) renderNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.Document:
		r.RenderNodes(n.Children)
	default:
		r.output.WriteString(n.Restore())
	}
}

// String returns the accumulated markdown string.
func (r *MarkdownRenderer) String() string {
	return r.output.String()
}
