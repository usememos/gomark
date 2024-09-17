package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestItalicParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "*Hello world!",
			node: nil,
		},
		{
			text: "*Hello*",
			node: &ast.Italic{
				Symbol:  "*",
				Content: "Hello",
			},
		},
		{
			text: "* Hello *",
			node: &ast.Italic{
				Symbol:  "*",
				Content: " Hello ",
			},
		},
		{
			text: "*1* Hello * *",
			node: &ast.Italic{
				Symbol:  "*",
				Content: "1",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewItalicParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
