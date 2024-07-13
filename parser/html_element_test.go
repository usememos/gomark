package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
	"github.com/usememos/gomark/restore"
)

func TestHTMLElementParser(t *testing.T) {
	tests := []struct {
		text        string
		htmlElement ast.Node
	}{
		{
			text: "<br />",
			htmlElement: &ast.HTMLElement{
				TagName: "br",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := NewHTMLElementParser().Match(tokens)
		require.Equal(t, restore.Restore([]ast.Node{test.htmlElement}), restore.Restore([]ast.Node{node}))
	}
}
