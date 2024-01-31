package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yourselfhosted/gomark/ast"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
	"github.com/yourselfhosted/gomark/restore"
)

func TestTaskListParser(t *testing.T) {
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
			node: &ast.TaskList{
				Symbol:   tokenizer.PlusSign,
				Complete: false,
				Children: []ast.Node{
					&ast.Text{
						Content: "Hello World",
					},
				},
			},
		},
		{
			text: "  + [ ] Hello World",
			node: &ast.TaskList{
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
			node: &ast.TaskList{
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
		node, _ := NewTaskListParser().Match(tokens)
		require.Equal(t, restore.Restore([]ast.Node{test.node}), restore.Restore([]ast.Node{node}))
	}
}
