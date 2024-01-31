package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yourselfhosted/gomark/ast"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
	"github.com/yourselfhosted/gomark/restore"
)

func TestHeadingParser(t *testing.T) {
	tests := []struct {
		text    string
		heading ast.Node
	}{
		{
			text:    "*Hello world",
			heading: nil,
		},
		{
			text: "## Hello World\n123",
			heading: &ast.Heading{
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
			heading: &ast.Heading{
				Level: 1,
				Children: []ast.Node{
					&ast.Text{
						Content: "# Hello World",
					},
				},
			},
		},
		{
			text:    " # 123123 Hello World",
			heading: nil,
		},
		{
			text: `# 123 
Hello World`,
			heading: &ast.Heading{
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
			heading: &ast.Heading{
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
		node, _ := NewHeadingParser().Match(tokens)
		require.Equal(t, restore.Restore([]ast.Node{test.heading}), restore.Restore([]ast.Node{node}))
	}
}
