package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestLinkParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "[](https://example.com)",
			node: &ast.Link{
				Text: "",
				URL:  "https://example.com",
			},
		},
		{
			text: "! [](https://example.com)",
			node: nil,
		},
		{
			text: "[alte]( htt ps :/ /example.com)",
			node: nil,
		},
		{
			text: "[your/slash](https://example.com)",
			node: &ast.Link{
				Text: "your/slash",
				URL:  "https://example.com",
			},
		},
		{
			text: "[hello world](https://example.com)",
			node: &ast.Link{
				Text: "hello world",
				URL:  "https://example.com",
			},
		},
	}
	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewLinkParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
