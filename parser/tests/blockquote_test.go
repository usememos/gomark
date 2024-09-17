package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestBlockquoteParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: ">Hello world",
			node: nil,
		},
		{
			text: "> Hello world",
			node: &ast.Blockquote{
				Children: []ast.Node{
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{
								Content: "Hello world",
							},
						},
					},
				},
			},
		},
		{
			text: "> 你好",
			node: &ast.Blockquote{
				Children: []ast.Node{
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{
								Content: "你好",
							},
						},
					},
				},
			},
		},
		{
			text: "> Hello\n> world",
			node: &ast.Blockquote{
				Children: []ast.Node{
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{
								Content: "Hello",
							},
						},
					},
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{
								Content: "world",
							},
						},
					},
				},
			},
		},
		{
			text: "> Hello\n> > world",
			node: &ast.Blockquote{
				Children: []ast.Node{
					&ast.Paragraph{
						Children: []ast.Node{
							&ast.Text{
								Content: "Hello",
							},
						},
					},
					&ast.Blockquote{
						Children: []ast.Node{
							&ast.Paragraph{
								Children: []ast.Node{
									&ast.Text{
										Content: "world",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewBlockquoteParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
