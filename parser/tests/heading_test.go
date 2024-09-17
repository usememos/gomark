package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestHeadingParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "*Hello world",
			node: nil,
		},
		{
			text: "## Hello World\n123",
			node: &ast.Heading{
				Level: 2,
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello World",
					},
				},
			},
		},
		{
			text: "# # Hello World",
			node: &ast.Heading{
				Level: 1,
				Children: []ast.Node{
					&ast.Text{
						Content: "# Hello World",
					},
				},
			},
		},
		{
			text: " # 123123 Hello World",
			node: nil,
		},
		{
			text: `# 123 
Hello World`,
			node: &ast.Heading{
				Level: 1,
				Children: []ast.Node{
					&ast.Text{
						Content: "123 ",
					},
				},
			},
		},
		{
			text: "### **Hello** World",
			node: &ast.Heading{
				Level: 3,
				Children: []ast.Node{
					&ast.Bold{
						Symbol: "*",
						Children: []ast.Node{
							&ast.Text{
								Content: "Hello",
							},
						},
					},
					&ast.Text{
						Content: " World",
					},
				},
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewHeadingParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
