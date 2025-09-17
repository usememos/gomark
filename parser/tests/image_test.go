package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/internal"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestImageParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "![](https://example.com)",
			node: &ast.Image{
				AltText: "",
				URL:     "https://example.com",
			},
		},
		{
			text: "! [](https://example.com)",
			node: nil,
		},
		{
			text: "![alte]( htt ps :/ /example.com)",
			node: nil,
		},
		{
			text: "![al te](https://example.com)",
			node: &ast.Image{
				AltText: "al te",
				URL:     "https://example.com",
			},
		},
	}
	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := internal.NewImageParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
