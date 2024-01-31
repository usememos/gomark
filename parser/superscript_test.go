package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yourselfhosted/gomark/ast"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
	"github.com/yourselfhosted/gomark/restore"
)

func TestSuperscriptParser(t *testing.T) {
	tests := []struct {
		text        string
		superscript ast.Node
	}{
		{
			text:        "^Hello world!",
			superscript: nil,
		},
		{
			text: "^Hello^",
			superscript: &ast.Superscript{
				Content: "Hello",
			},
		},
		{
			text: "^ Hello ^",
			superscript: &ast.Superscript{
				Content: " Hello ",
			},
		},
		{
			text: "^1^ Hello ^ ^",
			superscript: &ast.Superscript{
				Content: "1",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := NewSuperscriptParser().Match(tokens)
		require.Equal(t, restore.Restore([]ast.Node{test.superscript}), restore.Restore([]ast.Node{node}))
	}
}
