package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestListParser(t *testing.T) {
	tests := []struct {
		text  string
		nodes []ast.Node
	}{
		{
			text: "1. hello\n\n",
			nodes: []ast.Node{
				&ast.List{
					Kind: ast.OrderedList,
					Children: []ast.Node{
						&ast.OrderedListItem{
							Number: "1",
							Children: []ast.Node{
								&ast.Text{
									Content: "hello",
								},
							},
						},
						&ast.LineBreak{},
						&ast.LineBreak{},
					},
				},
			},
		},
		{
			text: "1. hello\n2. world",
			nodes: []ast.Node{
				&ast.List{
					Kind: ast.OrderedList,
					Children: []ast.Node{
						&ast.OrderedListItem{
							Number: "1",
							Children: []ast.Node{
								&ast.Text{
									Content: "hello",
								},
							},
						},
						&ast.LineBreak{},
						&ast.OrderedListItem{
							Number: "2",
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
		{
			text: "1. hello\n  2. world",
			nodes: []ast.Node{
				&ast.List{
					Kind: ast.OrderedList,
					Children: []ast.Node{
						&ast.OrderedListItem{
							Number: "1",
							Children: []ast.Node{
								&ast.Text{
									Content: "hello",
								},
							},
						},
						&ast.LineBreak{},
						&ast.List{
							Kind:   ast.OrderedList,
							Indent: 2,
							Children: []ast.Node{
								&ast.OrderedListItem{
									Number: "2",
									Indent: 2,
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
		},
		{
			text: "* hello\n  * world\n  * gomark",
			nodes: []ast.Node{
				&ast.List{
					Kind: ast.UnorderedList,
					Children: []ast.Node{
						&ast.UnorderedListItem{
							Symbol: "*",
							Children: []ast.Node{
								&ast.Text{
									Content: "hello",
								},
							},
						},
						&ast.LineBreak{},
						&ast.List{
							Kind:   ast.UnorderedList,
							Indent: 2,
							Children: []ast.Node{
								&ast.UnorderedListItem{
									Symbol: "*",
									Indent: 2,
									Children: []ast.Node{
										&ast.Text{
											Content: "world",
										},
									},
								},
								&ast.LineBreak{},
								&ast.UnorderedListItem{
									Symbol: "*",
									Indent: 2,
									Children: []ast.Node{
										&ast.Text{
											Content: "gomark",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			text: "* hello\n  * world\n* gomark",
			nodes: []ast.Node{
				&ast.List{
					Kind: ast.UnorderedList,
					Children: []ast.Node{
						&ast.UnorderedListItem{
							Symbol: "*",
							Children: []ast.Node{
								&ast.Text{
									Content: "hello",
								},
							},
						},
						&ast.LineBreak{},
						&ast.List{
							Kind:   ast.UnorderedList,
							Indent: 2,
							Children: []ast.Node{
								&ast.UnorderedListItem{
									Symbol: "*",
									Indent: 2,
									Children: []ast.Node{
										&ast.Text{
											Content: "world",
										},
									},
								},
								&ast.LineBreak{},
							},
						},
						&ast.UnorderedListItem{
							Symbol: "*",
							Children: []ast.Node{
								&ast.Text{
									Content: "gomark",
								},
							},
						},
					},
				},
			},
		},
		{
			text: "* hello\nparagraph\n* world",
			nodes: []ast.Node{
				&ast.List{
					Kind: ast.UnorderedList,
					Children: []ast.Node{
						&ast.UnorderedListItem{
							Symbol: "*",
							Children: []ast.Node{
								&ast.Text{
									Content: "hello",
								},
							},
						},
						&ast.LineBreak{},
					},
				},
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "paragraph",
						},
					},
				},
				&ast.LineBreak{},
				&ast.List{
					Kind: ast.UnorderedList,
					Children: []ast.Node{
						&ast.UnorderedListItem{
							Symbol: "*",
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
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		nodes, _ := parser.Parse(tokens)
		require.ElementsMatch(t, test.nodes, nodes, fmt.Sprintf("Test case: %s", test.text))
	}
}
