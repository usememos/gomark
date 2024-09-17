package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestSubscriptParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "~Hello world!",
			node: nil,
		},
		{
			text: "~Hello~",
			node: &ast.Subscript{
				Content: "Hello",
			},
		},
		{
			text: "~ Hello ~",
			node: &ast.Subscript{
				Content: " Hello ",
			},
		},
		{
			text: "~1~ Hello ~ ~",
			node: &ast.Subscript{
				Content: "1",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewSubscriptParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
