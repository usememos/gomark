package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type HorizontalRuleParser struct{}

func NewHorizontalRuleParser() *HorizontalRuleParser {
	return &HorizontalRuleParser{}
}

func (*HorizontalRuleParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedTokens := tokenizer.GetFirstLine(tokens)
	if len(matchedTokens) < 1 {
		return nil, 0
	}

	// Check for new-style triple hyphen token
	if matchedTokens[0].Type == tokenizer.TripleHyphen {
		if len(matchedTokens) > 1 && matchedTokens[1].Type != tokenizer.NewLine {
			return nil, 0
		}
		return &ast.HorizontalRule{
			Symbol: "-",
		}, 1
	}

	// Legacy parsing for individual tokens
	if len(matchedTokens) < 3 {
		return nil, 0
	}
	if len(matchedTokens) > 3 && matchedTokens[3].Type != tokenizer.NewLine {
		return nil, 0
	}
	if matchedTokens[0].Type != matchedTokens[1].Type || matchedTokens[0].Type != matchedTokens[2].Type || matchedTokens[1].Type != matchedTokens[2].Type {
		return nil, 0
	}
	if matchedTokens[0].Type != tokenizer.Hyphen && matchedTokens[0].Type != tokenizer.Asterisk {
		return nil, 0
	}
	return &ast.HorizontalRule{
		Symbol: matchedTokens[0].Type,
	}, 3
}
