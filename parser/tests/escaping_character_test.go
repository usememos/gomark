package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestEscapingCharacterParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: `\#`,
			node: &ast.EscapingCharacter{
				Symbol: "#",
			},
		},
		{
			text: `\' test`,
			node: &ast.EscapingCharacter{
				Symbol: "'",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewEscapingCharacterParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
