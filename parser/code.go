package parser

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type CodeParser struct{}

func NewCodeParser() *CodeParser {
	return &CodeParser{}
}

func (*CodeParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedTokens := tokenizer.GetFirstLine(tokens)
	if len(matchedTokens) < 3 {
		return nil, 0
	}
	if matchedTokens[0].Type != tokenizer.Backtick {
		return nil, 0
	}
	nextBacktickIndex := tokenizer.FindUnescaped(matchedTokens[1:], tokenizer.Backtick)
	if nextBacktickIndex < 0 {
		return nil, 0
	}
	matchedTokens = matchedTokens[:1+nextBacktickIndex+1]
	contentTokens := matchedTokens[1 : len(matchedTokens)-1]
	if len(contentTokens) == 0 {
		return nil, 0
	}
	return &ast.Code{
		Content: tokenizer.Stringify(contentTokens),
	}, len(matchedTokens)
}
