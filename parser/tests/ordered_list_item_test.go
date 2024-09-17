package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestOrderedListItemParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "1.asd",
			node: nil,
		},
		{
			text: "1. Hello World",
			node: &ast.OrderedListItem{
				Number: "1",
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello World",
					},
				},
			},
		},
		{
			text: "  1. Hello World",
			node: &ast.OrderedListItem{
				Number: "1",
				Indent: 2,
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello World",
					},
				},
			},
		},
		{
			text: "1aa. Hello World",
			node: nil,
		},
		{
			text: "22. Hello *World*",
			node: &ast.OrderedListItem{
				Number: "22",
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello ",
					},
					&ast.Italic{
						Symbol:  "*",
						Content: "World",
					},
				},
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewOrderedListItemParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
