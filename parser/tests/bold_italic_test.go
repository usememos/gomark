package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/internal"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestBoldItalicParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "*Hello world!",
			node: nil,
		},
		{
			text: "*** Hello * *",
			node: nil,
		},
		{
			text: "*** Hello **",
			node: nil,
		},
		{
			text: "***Hello***",
			node: &ast.BoldItalic{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Text{Content: "Hello"},
				},
			},
		},
		{
			text: "___Hello___",
			node: &ast.BoldItalic{
				Symbol: "_",
				Children: []ast.Node{
					&ast.Text{Content: "Hello"},
				},
			},
		},
		{
			text: "*** Hello ***",
			node: &ast.BoldItalic{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Text{Content: " Hello "},
				},
			},
		},
		{
			text: "***Hello _world_***",
			node: &ast.BoldItalic{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Text{Content: "Hello "},
					&ast.Italic{
						Symbol: "_",
						Children: []ast.Node{
							&ast.Text{Content: "world"},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := internal.NewBoldItalicParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
