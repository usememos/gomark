package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestHTMLElementParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "<br />",
			node: &ast.HTMLElement{
				TagName:    "br",
				Attributes: map[string]string{},
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewHTMLElementParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
