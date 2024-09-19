package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestParser(t *testing.T) {
	tests := []struct {
		text  string
		nodes []ast.Node
	}{
		{
			text: "Hello world!",
			nodes: []ast.Node{
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "Hello world!",
						},
					},
				},
			},
		},
		{
			text: "# Hello world!",
			nodes: []ast.Node{
				&ast.Heading{
					Level: 1,
					Children: []ast.Node{
						&ast.Text{
							Content: "Hello world!",
						},
					},
				},
			},
		},
		{
			text: "\\# Hello world!",
			nodes: []ast.Node{
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.EscapingCharacter{
							Symbol: "#",
						},
						&ast.Text{
							Content: " Hello world!",
						},
					},
				},
			},
		},
		{
			text: "**Hello** world!",
			nodes: []ast.Node{
				&ast.Paragraph{
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
							Content: " world!",
						},
					},
				},
			},
		},
		{
			text: "Hello **world**!\nHere is a new line.",
			nodes: []ast.Node{
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "Hello ",
						},
						&ast.Bold{
							Symbol: "*",
							Children: []ast.Node{
								&ast.Text{
									Content: "world",
								},
							},
						},
						&ast.Text{
							Content: "!",
						},
					},
				},
				&ast.LineBreak{},
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "Here is a new line.",
						},
					},
				},
			},
		},
		{
			text: "Hello **world**!\n```javascript\nconsole.log(\"Hello world!\");\n```",
			nodes: []ast.Node{
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "Hello ",
						},
						&ast.Bold{
							Symbol: "*",
							Children: []ast.Node{
								&ast.Text{
									Content: "world",
								},
							},
						},
						&ast.Text{
							Content: "!",
						},
					},
				},
				&ast.LineBreak{},
				&ast.CodeBlock{
					Language: "javascript",
					Content:  "console.log(\"Hello world!\");",
				},
			},
		},
		{
			text: "Hello world!\n\nNew paragraph.",
			nodes: []ast.Node{
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "Hello world!",
						},
					},
				},
				&ast.LineBreak{},
				&ast.LineBreak{},
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "New paragraph.",
						},
					},
				},
			},
		},
		{
			text: "1. hello\n- [ ] world",
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
					},
				},
				&ast.List{
					Kind: ast.DescrpitionList,
					Children: []ast.Node{
						&ast.TaskListItem{
							Symbol:   tokenizer.Hyphen,
							Complete: false,
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
			text: "- [ ] hello\n- [x] world",
			nodes: []ast.Node{
				&ast.List{
					Kind: ast.DescrpitionList,
					Children: []ast.Node{
						&ast.TaskListItem{
							Symbol:   tokenizer.Hyphen,
							Complete: false,
							Children: []ast.Node{
								&ast.Text{
									Content: "hello",
								},
							},
						},
						&ast.LineBreak{},
						&ast.TaskListItem{
							Symbol:   tokenizer.Hyphen,
							Complete: true,
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
			text: "\n\n",
			nodes: []ast.Node{
				&ast.LineBreak{},
				&ast.LineBreak{},
			},
		},
		{
			text: "\n$$\na=3\n$$",
			nodes: []ast.Node{
				&ast.LineBreak{},
				&ast.MathBlock{
					Content: "a=3",
				},
			},
		},
		{
			text: "Hello\n![[memos/101]]",
			nodes: []ast.Node{
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "Hello",
						},
					},
				},
				&ast.LineBreak{},
				&ast.EmbeddedContent{
					ResourceName: "memos/101",
				},
			},
		},
		{
			text: "Hello\nworld<br />",
			nodes: []ast.Node{
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "Hello",
						},
					},
				},
				&ast.LineBreak{},
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "world",
						},
						&ast.HTMLElement{
							TagName:    "br",
							Attributes: map[string]string{},
						},
					},
				},
			},
		},
		{
			text: "Hello <br /> world",
			nodes: []ast.Node{
				&ast.Paragraph{
					Children: []ast.Node{
						&ast.Text{
							Content: "Hello ",
						},
						&ast.HTMLElement{
							TagName:    "br",
							Attributes: map[string]string{},
						},
						&ast.Text{
							Content: " world",
						},
					},
				},
			},
		},
		{
			text: "* unordered list item 1\n* unordered list item 2",
			nodes: []ast.Node{
				&ast.List{
					Kind: ast.UnorderedList,
					Children: []ast.Node{
						&ast.UnorderedListItem{
							Symbol: tokenizer.Asterisk,
							Children: []ast.Node{
								&ast.Text{
									Content: "unordered list item 1",
								},
							},
						},
						&ast.LineBreak{},
						&ast.UnorderedListItem{
							Symbol: tokenizer.Asterisk,
							Children: []ast.Node{
								&ast.Text{
									Content: "unordered list item 2",
								},
							},
						},
					},
				},
			},
		},
		{
			text: "* unordered list item\n\n1. ordered list item",
			nodes: []ast.Node{
				&ast.List{
					Kind: ast.UnorderedList,
					Children: []ast.Node{
						&ast.UnorderedListItem{
							Symbol: tokenizer.Asterisk,
							Children: []ast.Node{
								&ast.Text{
									Content: "unordered list item",
								},
							},
						},
					},
				},
				&ast.LineBreak{},
				&ast.LineBreak{},
				&ast.List{
					Kind: ast.OrderedList,
					Children: []ast.Node{
						&ast.OrderedListItem{
							Number: "1",
							Children: []ast.Node{
								&ast.Text{
									Content: "ordered list item",
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
		require.Equal(t, test.nodes, nodes, fmt.Sprintf("Test case: %s", test.text))
	}
}
