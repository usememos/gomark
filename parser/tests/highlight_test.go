package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestHighlightParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "==Hello world!",
			node: nil,
		},
		{
			text: "==Hello==",
			node: &ast.Highlight{
				Content: "Hello",
			},
		},
		{
			text: "==Hello world==",
			node: &ast.Highlight{
				Content: "Hello world",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewHighlightParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}