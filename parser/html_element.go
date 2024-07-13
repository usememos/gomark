package parser

import (
	"slices"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type HTMLElementParser struct{}

func NewHTMLElementParser() *HTMLElementParser {
	return &HTMLElementParser{}
}

var (
	availableHTMLElements = []string{
		"br",
	}
)

func (*HTMLElementParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	if len(tokens) < 5 {
		return nil, 0
	}
	if tokens[0].Type != tokenizer.LessThan {
		return nil, 0
	}
	tagName := tokenizer.Stringify([]*tokenizer.Token{tokens[1]})
	if !slices.Contains(availableHTMLElements, tagName) {
		return nil, 0
	}

	greaterThanIndex := tokenizer.FindUnescaped(tokens, tokenizer.GreaterThan)
	if greaterThanIndex+1 < 5 || tokens[greaterThanIndex-1].Type != tokenizer.Slash || tokens[greaterThanIndex-2].Type != tokenizer.Space {
		return nil, 0
	}

	matchedTokens := tokens[:greaterThanIndex]
	attributeTokens := matchedTokens[2 : greaterThanIndex-2]
	// TODO: Implement attribute parser.
	if len(attributeTokens) != 0 {
		return nil, 0
	}

	return &ast.HTMLElement{
		TagName:    tagName,
		Attributes: make(map[string]string),
	}, len(matchedTokens)
}
