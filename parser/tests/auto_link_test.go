package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestAutoLinkParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "<https://example.com)",
			node: nil,
		},
		{
			text: "<https://example.com>",
			node: &ast.AutoLink{
				URL: "https://example.com",
			},
		},
		{
			text: "https://example.com",
			node: &ast.AutoLink{
				URL:       "https://example.com",
				IsRawText: true,
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewAutoLinkParser().Match(tokens)
		require.Equal(t, node, test.node)
	}
}
