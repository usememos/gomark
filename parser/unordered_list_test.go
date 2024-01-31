package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yourselfhosted/gomark/ast"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
	"github.com/yourselfhosted/gomark/restore"
)

func TestUnorderedListParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "*asd",
			node: nil,
		},
		{
			text: "+ Hello World",
			node: &ast.UnorderedList{
				Symbol: tokenizer.PlusSign,
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello World",
					},
				},
			},
		},
		{
			text: "* **Hello**",
			node: &ast.UnorderedList{
				Symbol: tokenizer.Asterisk,
				Children: []ast.Node{
					&ast.Bold{
						Symbol: "*",
						Children: []ast.Node{
							&ast.Text{
								Content: "Hello",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := NewUnorderedListParser().Match(tokens)
		require.Equal(t, restore.Restore([]ast.Node{test.node}), restore.Restore([]ast.Node{node}))
	}
}
