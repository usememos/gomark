package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestHorizontalRuleParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "---",
			node: &ast.HorizontalRule{
				Symbol: "-",
			},
		},
		{
			text: "---\naaa",
			node: &ast.HorizontalRule{
				Symbol: "-",
			},
		},
		{
			text: "****",
			node: nil,
		},
		{
			text: "***",
			node: &ast.HorizontalRule{
				Symbol: "*",
			},
		},
		{
			text: "-*-",
			node: nil,
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewHorizontalRuleParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
