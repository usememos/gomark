package parser

import (
	"github.com/yourselfhosted/gomark/ast"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
)

type LineBreakParser struct{}

func NewLineBreakParser() *LineBreakParser {
	return &LineBreakParser{}
}

func (*LineBreakParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	if len(tokens) == 0 {
		return nil, 0
	}
	if tokens[0].Type != tokenizer.NewLine {
		return nil, 0
	}
	return &ast.LineBreak{
		BaseBlock: ast.NewBaseBlock(ast.LineBreakNode),
	}, 1
}
