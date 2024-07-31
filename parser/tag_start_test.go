package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
	"github.com/usememos/gomark/restore"
)

func TestTagStartParser(t *testing.T) {
	tests := []struct {
		text string
		tag  ast.Node
	}{
		{
			text: "# Hello World",
			tag:  nil,
		},
		{
			text: "#HelloWorld",
			tag:  nil,
		},
		{
			text: "#[[]]",
			tag:  nil,
		},
		{
			text: "#[[ ]]",
			tag:  nil,
		},
		{
			text: "#[[tag",
			tag:  nil,
		},
		{
			text: "#[[tag ]]",
			tag:  nil,
		},
		{
			text: " #[[inline]] tag",
			tag:  nil,
		},
		{
			text: "#[[x]]",
			tag: &ast.Tag{
				Content: "x",
			},
		},
		{
			text: "#[[foo]]",
			tag: &ast.Tag{
				Content: "foo",
			},
		},
		{
			text: "#[[foo]]bar",
			tag: &ast.Tag{
				Content: "foo",
			},
		},
		{
			text: "#[[foo]] bar",
			tag: &ast.Tag{
				Content: "foo",
			},
		},
		{
			text: "#tag/subtag",
			tag:  nil,
		},
		{
			text: "#[[tag/subtag]]",
			tag: &ast.Tag{
				Content: "tag/subtag",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := NewTagStartParser().Match(tokens)
		require.Equal(t, restore.Restore([]ast.Node{test.tag}), restore.Restore([]ast.Node{node}))
	}
}
