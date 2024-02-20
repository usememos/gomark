package parser

import (
	"github.com/yourselfhosted/gomark/ast"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
)

type SpoilerParser struct{}

func NewSpoilerParser() InlineParser {
	return &SpoilerParser{}
}

func (*SpoilerParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedTokens := tokenizer.GetFirstLine(tokens)
	if len(matchedTokens) < 5 {
		return nil, 0
	}

	prefixTokens := matchedTokens[:2]
	if prefixTokens[0].Type != prefixTokens[1].Type {
		return nil, 0
	}
	prefixTokenType := prefixTokens[0].Type
	if prefixTokenType != tokenizer.Pipe {
		return nil, 0
	}

	cursor, matched := 2, false
	for ; cursor < len(matchedTokens)-1; cursor++ {
		token, nextToken := matchedTokens[cursor], matchedTokens[cursor+1]
		if token.Type == tokenizer.NewLine || nextToken.Type == tokenizer.NewLine {
			return nil, 0
		}
		if token.Type == prefixTokenType && nextToken.Type == prefixTokenType {
			matchedTokens = matchedTokens[:cursor+2]
			matched = true
			break
		}
	}
	if !matched {
		return nil, 0
	}

	size := len(matchedTokens)
	content := tokenizer.Stringify(matchedTokens[2 : size-2])
	return &ast.Spoiler{
		Content: content,
	}, size
}
