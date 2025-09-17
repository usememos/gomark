package markdown

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
)

func TestMarkdownRenderer(t *testing.T) {
	doc := &ast.Document{
		Children: []ast.Node{
			&ast.Paragraph{Children: []ast.Node{&ast.Text{Content: "Hello"}}},
			&ast.LineBreak{},
			&ast.Paragraph{Children: []ast.Node{&ast.Text{Content: "World"}}},
		},
	}

	r := NewMarkdownRenderer()
	r.RenderDocument(doc)
	require.Equal(t, doc.Restore(), r.String())
}
