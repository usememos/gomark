package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type OrderedListItemParser struct{}

func NewOrderedListItemParser() *OrderedListItemParser {
	return &OrderedListItemParser{}
}

func (*OrderedListItemParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedTokens := tokenizer.GetFirstLine(tokens)
	indent := 0
	consumedTokens := 0

	// Handle both new-style consolidated spaces and legacy individual spaces.
spaceScan:
	for i, token := range matchedTokens {
		switch token.Type {
		case tokenizer.Space:
			indent++
			consumedTokens = i + 1
		case tokenizer.MultipleSpaces:
			indent = len(token.Value) // Count the actual spaces in the consolidated token.
			consumedTokens = i + 1
		default:
			break spaceScan
		}
	}

	if len(matchedTokens) < consumedTokens+3 {
		return nil, 0
	}

	cursor := consumedTokens
	if matchedTokens[cursor].Type != tokenizer.Number || matchedTokens[cursor+1].Type != tokenizer.Dot || matchedTokens[cursor+2].Type != tokenizer.Space {
		return nil, 0
	}

	contentTokens := matchedTokens[cursor+3:]
	if len(contentTokens) == 0 {
		return nil, 0
	}
	children, err := ParseInline(contentTokens)
	if err != nil {
		return nil, 0
	}
	return &ast.OrderedListItem{
		Number:   matchedTokens[cursor].Value,
		Indent:   indent,
		Children: children,
	}, cursor + 3 + len(contentTokens)
}
