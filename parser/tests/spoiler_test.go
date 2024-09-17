package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestSpoilerParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "*Hello world!",
			node: nil,
		},
		{
			text: "||Hello||",
			node: &ast.Spoiler{
				Content: "Hello",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewSpoilerParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
