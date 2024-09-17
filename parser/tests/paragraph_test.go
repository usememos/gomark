package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestParagraphParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "",
			node: nil,
		},
		{
			text: "\n",
			node: nil,
		},
		{
			text: "Hello world!",
			node: &ast.Paragraph{
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello world!",
					},
				},
			},
		},
		{
			text: "Hello world!\n",
			node: &ast.Paragraph{
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello world!",
					},
				},
			},
		},
		{
			text: "Hello world!\n\nNew paragraph.",
			node: &ast.Paragraph{
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello world!",
					},
				},
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewParagraphParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
