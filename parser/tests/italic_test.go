package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/internal"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestItalicParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "*Hello world!",
			node: nil,
		},
		{
			text: "*Hello*",
			node: &ast.Italic{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello",
					},
				},
			},
		},
		{
			text: "*`code`*",
			node: &ast.Italic{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Code{Content: "code"},
				},
			},
		},
		{
			text: "* Hello *",
			node: &ast.Italic{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Text{
						Content: " Hello ",
					},
				},
			},
		},
		{
			text: "*1* Hello * *",
			node: &ast.Italic{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Text{
						Content: "1",
					},
				},
			},
		},
		{
			text: "_Hello_",
			node: &ast.Italic{
				Symbol: "_",
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello",
					},
				},
			},
		},
		{
			text: "_ Hello _",
			node: &ast.Italic{
				Symbol: "_",
				Children: []ast.Node{
					&ast.Text{
						Content: " Hello ",
					},
				},
			},
		},
		{
			text: "_1_ Hello _ _",
			node: &ast.Italic{
				Symbol: "_",
				Children: []ast.Node{
					&ast.Text{
						Content: "1",
					},
				},
			},
		},
		{
			text: "*[Hello](https://example.com)*",
			node: &ast.Italic{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Link{
						Content: []ast.Node{
							&ast.Text{
								Content: "Hello",
							},
						},
						URL: "https://example.com",
					},
				},
			},
		},
		{
			text: "*Hello **world***",
			node: &ast.Italic{
				Symbol: "*",
				Children: []ast.Node{
					&ast.Text{Content: "Hello "},
					&ast.Bold{
						Symbol: "*",
						Children: []ast.Node{
							&ast.Text{Content: "world"},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := internal.NewItalicParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
