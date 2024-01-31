package parser

import (
	"github.com/yourselfhosted/gomark/ast"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
)

type ParagraphParser struct {
	ContentTokens []*tokenizer.Token
}

func NewParagraphParser() *ParagraphParser {
	return &ParagraphParser{}
}

func (*ParagraphParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedTokens := tokenizer.GetFirstLine(tokens)
	if len(matchedTokens) == 0 {
		return nil, 0
	}

	children, err := ParseInline(matchedTokens)
	if err != nil {
		return nil, 0
	}
	return &ast.Paragraph{
		Children: children,
	}, len(matchedTokens)
}
