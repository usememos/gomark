package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yourselfhosted/gomark/ast"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
	"github.com/yourselfhosted/gomark/restore"
)

func TestHighlightParser(t *testing.T) {
	tests := []struct {
		text string
		bold ast.Node
	}{
		{
			text: "==Hello world!",
			bold: nil,
		},
		{
			text: "==Hello==",
			bold: &ast.Highlight{
				Content: "Hello",
			},
		},
		{
			text: "==Hello world==",
			bold: &ast.Highlight{
				Content: "Hello world",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := NewHighlightParser().Match(tokens)
		require.Equal(t, restore.Restore([]ast.Node{test.bold}), restore.Restore([]ast.Node{node}))
	}
}
