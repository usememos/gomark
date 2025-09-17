package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/internal"
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
				Content: []ast.Node{},
				URL:     "https://example.com",
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
				Content: []ast.Node{&ast.Text{Content: "your/slash"}},
				URL:     "https://example.com",
			},
		},
		{
			text: "[hello world](https://example.com)",
			node: &ast.Link{
				Content: []ast.Node{&ast.Text{Content: "hello world"}},
				URL:     "https://example.com",
			},
		},
		{
			text: "[hello world](https://example.com)",
			node: &ast.Link{
				Content: []ast.Node{&ast.Text{Content: "hello world"}},
				URL:     "https://example.com",
			},
		},
		{
			text: `[\[link\]](https://example.com)`,
			node: &ast.Link{
				Content: []ast.Node{&ast.EscapingCharacter{Symbol: "["}, &ast.Text{Content: `link`}, &ast.EscapingCharacter{Symbol: "]"}},
				URL:     "https://example.com",
			},
		},
	}
	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := internal.NewLinkParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
