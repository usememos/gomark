package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type UnorderedListItemParser struct{}

func NewUnorderedListItemParser() *UnorderedListItemParser {
	return &UnorderedListItemParser{}
}

func (*UnorderedListItemParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
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

	if len(matchedTokens) < consumedTokens+2 {
		return nil, 0
	}

	symbolToken := matchedTokens[consumedTokens]
	if (symbolToken.Type != tokenizer.Hyphen && symbolToken.Type != tokenizer.Asterisk && symbolToken.Type != tokenizer.PlusSign) || matchedTokens[consumedTokens+1].Type != tokenizer.Space {
		return nil, 0
	}

	contentTokens := matchedTokens[consumedTokens+2:]
	if len(contentTokens) == 0 {
		return nil, 0
	}
	children, err := ParseInline(contentTokens)
	if err != nil {
		return nil, 0
	}
	return &ast.UnorderedListItem{
		Symbol:   symbolToken.Type,
		Indent:   indent,
		Children: children,
	}, consumedTokens + 2 + len(contentTokens)
}
