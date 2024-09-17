package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestCodeBlockParser(t *testing.T) {
	tests := []struct {
		text string
		node ast.Node
	}{
		{
			text: "```Hello world!```",
			node: nil,
		},
		{
			text: "```\nHello\n```",
			node: &ast.CodeBlock{
				Language: "",
				Content:  "Hello",
			},
		},
		{
			text: "```\nHello world!\n```",
			node: &ast.CodeBlock{
				Language: "",
				Content:  "Hello world!",
			},
		},
		{
			text: "```java\nHello \n world!\n```",
			node: &ast.CodeBlock{
				Language: "java",
				Content:  "Hello \n world!",
			},
		},
		{
			text: "```java\nHello \n world!\n```111",
			node: nil,
		},
		{
			text: "```java\nHello \n world!\n``` 111",
			node: nil,
		},
		{
			text: "```java\nHello \n world!\n```\n123123",
			node: &ast.CodeBlock{
				Language: "java",
				Content:  "Hello \n world!",
			},
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := parser.NewCodeBlockParser().Match(tokens)
		require.Equal(t, test.node, node)
	}
}
