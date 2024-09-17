package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestMathParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "$\\sqrt{3x-1}+(1+x)^2$",
			node: &ast.Math{
				Content: "\\sqrt{3x-1}+(1+x)^2",
			},
		},
	}
	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewMathParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
