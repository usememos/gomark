package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestTaskListItemParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "*asd",
			node: nil,
		},
		{
			text: "+ [ ] Hello World",
			node: &ast.TaskListItem{
				Symbol: tokenizer.PlusSign,
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello World",
					},
				},
			},
		},
		{
			text: "  + [ ] Hello World",
			node: &ast.TaskListItem{
				Symbol:   tokenizer.PlusSign,
				Indent:   2,
				Complete: false,
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello World",
					},
				},
			},
		},
		{
			text: "* [x] **Hello**",
			node: &ast.TaskListItem{
				Symbol:   tokenizer.Asterisk,
				Complete: true,
				Children: []ast.Node{
					&ast.Bold{
						Symbol: "*",
						Children: []ast.Node{
							&ast.Text{
								Content: "Hello",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewTaskListItemParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
