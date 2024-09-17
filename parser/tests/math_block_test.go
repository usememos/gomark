package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestMathBlockParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "$$\n(1+x)^2\n$$",
			node: &ast.MathBlock{
				Content: "(1+x)^2",
			},
		},
		{
			text: "$$\na=3\n$$",
			node: &ast.MathBlock{
				Content: "a=3",
			},
		},
	}
	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewMathBlockParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
