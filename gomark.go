package gomark

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
	"github.com/usememos/gomark/restore"
)

func Parse(markdown string) (nodes []ast.Node, err error) {
	tokens := tokenizer.Tokenize(markdown)
	return parser.Parse(tokens)
}

func Restore(nodes []ast.Node) string {
	return restore.Restore(nodes)
}
