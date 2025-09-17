package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type HighlightParser struct{}

func NewHighlightParser() InlineParser {
	return &HighlightParser{}
}

func (*HighlightParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedToken := tokenizer.GetFirstLine(tokens)
	if len(matchedToken) < 3 {
		return nil, 0
	}

	// Check for new-style double equal token
	if matchedToken[0].Type == tokenizer.DoubleEqual {
		// Find the closing token
		cursor, matched := 1, false
		for ; cursor < len(matchedToken); cursor++ {
			token := matchedToken[cursor]
			if token.Type == tokenizer.NewLine {
				return nil, 0
			}
			if token.Type == tokenizer.DoubleEqual {
				matchedToken = matchedToken[:cursor+1]
				matched = true
				break
			}
		}
		if !matched {
			return nil, 0
		}

		content := ""
		for _, token := range matchedToken[1 : len(matchedToken)-1] {
			content += token.Value
		}

		return &ast.Highlight{Content: content}, len(matchedToken)
	}

	// Legacy parsing for individual tokens
	if len(matchedToken) < 5 {
		return nil, 0
	}

	prefixTokens := matchedToken[:2]
	if prefixTokens[0].Type != prefixTokens[1].Type {
		return nil, 0
	}
	prefixTokenType := prefixTokens[0].Type
	if prefixTokenType != tokenizer.EqualSign {
		return nil, 0
	}

	cursor, matched := 2, false
	for ; cursor < len(matchedToken)-1; cursor++ {
		token, nextToken := matchedToken[cursor], matchedToken[cursor+1]
		if token.Type == prefixTokenType && nextToken.Type == prefixTokenType {
			matched = true
			break
		}
	}
	if !matched {
		return nil, 0
	}

	return &ast.Highlight{
		Content: tokenizer.Stringify(matchedToken[2:cursor]),
	}, cursor + 2
}
