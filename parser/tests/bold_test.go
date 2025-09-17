package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/internal"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestBoldParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "*Hello world!",
			node: nil,
		},
		{
			text: "****",
			node: nil,
		},
		{
			text: "**Hello**",
			node: &ast.Bold{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello",
					},
				},
			},
		},
		{
			text: "**`code`**",
			node: &ast.Bold{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Code{Content: "code"},
				},
			},
		},
		{
			text: "** Hello **",
			node: &ast.Bold{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Text{
						Content: " Hello ",
					},
				},
			},
		},
		{
			text: "** Hello * *",
			node: nil,
		},
		{
			text: "* * Hello **",
			node: nil,
		},
		{
			text: "__Hello__",
			node: &ast.Bold{
				Symbol: "_",
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello",
					},
				},
			},
		},
		{
			text: "**Nested _italic_**",
			node: &ast.Bold{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Text{Content: "Nested "},
					&ast.Italic{
						Symbol: "_",
						Children: []ast.Node{
							&ast.Text{Content: "italic"},
						},
					},
				},
			},
		},
		{
			text: "__ Hello __",
			node: &ast.Bold{
				Symbol: "_",
				Children: []ast.Node{
					&ast.Text{
						Content: " Hello ",
					},
				},
			},
		},
		{
			text: "__ Hello _ _",
			node: nil,
		},
		{
			text: "_ _ Hello __",
			node: nil,
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := internal.NewBoldParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
