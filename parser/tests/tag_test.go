package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestTagParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "*Hello world",
			node: nil,
		},
		{
			text: "# Hello World",
			node: nil,
		},
		{
			text: "#tag",
			node: &ast.Tag{
				Content: "tag",
			},
		},
		{
			text: "#tag/subtag 123",
			node: &ast.Tag{
				Content: "tag/subtag",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewTagParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
