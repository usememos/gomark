package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
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

	rightSquareBracketIndex := tokenizer.FindUnescaped(matchedTokens[1:], tokenizer.RightSquareBracket)
	if rightSquareBracketIndex == -1 {
		return nil, 0
	}
	contentTokens := matchedTokens[1 : rightSquareBracketIndex+1]
	if tokenizer.FindUnescaped(contentTokens, tokenizer.LeftSquareBracket) != -1 {
		return nil, 0
	}
	if len(contentTokens)+4 >= len(matchedTokens) {
		return nil, 0
	}
	if matchedTokens[2+len(contentTokens)].Type != tokenizer.LeftParenthesis {
		return nil, 0
	}
	urlTokens, matched := []*tokenizer.Token{}, false
	for _, token := range matchedTokens[3+len(contentTokens):] {
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

	contentNodes, err := ParseInlineWithParsers(contentTokens, []InlineParser{NewEscapingCharacterParser(), NewTextParser()})
	if err != nil {
		return nil, 0
	}
	return &ast.Link{
		Content: contentNodes,
		URL:     tokenizer.Stringify(urlTokens),
	}, 4 + len(contentTokens) + len(urlTokens)
}
