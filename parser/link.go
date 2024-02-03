package parser

import (
	"github.com/yourselfhosted/gomark/ast"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
)

type LinkParser struct{}

func NewLinkParser() *LinkParser {
	return &LinkParser{}
}

func (*LinkParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedTokens := tokenizer.GetFirstLine(tokens)
	if len(matchedTokens) < 5 {
		return nil, 0
	}
	if matchedTokens[0].Type != tokenizer.LeftSquareBracket {
		return nil, 0
	}

	textTokens := []*tokenizer.Token{}
	for _, token := range matchedTokens[1:] {
		if token.Type == tokenizer.RightSquareBracket {
			break
		}
		textTokens = append(textTokens, token)
	}
	if len(textTokens)+4 >= len(matchedTokens) {
		return nil, 0
	}
	if matchedTokens[2+len(textTokens)].Type != tokenizer.LeftParenthesis {
		return nil, 0
	}
	urlTokens, matched := []*tokenizer.Token{}, false
	for _, token := range matchedTokens[3+len(textTokens):] {
		if token.Type == tokenizer.Space {
			return nil, 0
		}
		if token.Type == tokenizer.RightParenthesis {
			matched = true
			break
		}
		urlTokens = append(urlTokens, token)
	}
	if !matched || len(urlTokens) == 0 {
		return nil, 0
	}

	return &ast.Link{
		Text: tokenizer.Stringify(textTokens),
		URL:  tokenizer.Stringify(urlTokens),
	}, 4 + len(textTokens) + len(urlTokens)
}
