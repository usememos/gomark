package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yourselfhosted/gomark/ast"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
	"github.com/yourselfhosted/gomark/restore"
)

func TestItalicParser(t *testing.T) {
	tests := []struct {
		text   string
		italic ast.Node
	}{
		{
			text:   "*Hello world!",
			italic: nil,
		},
		{
			text: "*Hello*",
			italic: &ast.Italic{
				Symbol:  "*",
				Content: "Hello",
			},
		},
		{
			text: "* Hello *",
			italic: &ast.Italic{
				Symbol:  "*",
				Content: " Hello ",
			},
		},
		{
			text: "*1* Hello * *",
			italic: &ast.Italic{
				Symbol:  "*",
				Content: "1",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := NewItalicParser().Match(tokens)
		require.Equal(t, restore.Restore([]ast.Node{test.italic}), restore.Restore([]ast.Node{node}))
	}
}
