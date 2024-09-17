package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestReferencedContentParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "[[Hello world]",
			node: nil,
		},
		{
			text: "[[Hello world]]",
			node: &ast.ReferencedContent{
				ResourceName: "Hello world",
			},
		},
		{
			text: "[[memos/1]]",
			node: &ast.ReferencedContent{
				ResourceName: "memos/1",
			},
		},
		{
			text: "[[resources/101]]111\n123",
			node: &ast.ReferencedContent{
				ResourceName: "resources/101",
			},
		},
		{
			text: "[[resources/101?align=center]]",
			node: &ast.ReferencedContent{
				ResourceName: "resources/101",
				Params:       "align=center",
			},
		},
		{
			text: "[[resources/6uxnhT98q8vN8anBbUbRGu?align=center]]",
			node: &ast.ReferencedContent{
				ResourceName: "resources/6uxnhT98q8vN8anBbUbRGu",
				Params:       "align=center",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewReferencedContentParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
