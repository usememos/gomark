package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type ItalicParser struct {
	ContentTokens []*tokenizer.Token
}

func NewItalicParser() *ItalicParser {
	return &ItalicParser{}
}

func (*ItalicParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedTokens := tokenizer.GetFirstLine(tokens)
	if len(matchedTokens) < 3 {
		return nil, 0
	}

	prefixTokens := matchedTokens[:1]
	if prefixTokens[0].Type != tokenizer.Asterisk && prefixTokens[0].Type != tokenizer.Underscore {
		return nil, 0
	}
	prefixTokenType := prefixTokens[0].Type
	contentTokens := []*tokenizer.Token{}
	closingIndex := -1

	for i := 1; i < len(matchedTokens); {
		token := matchedTokens[i]
		if token.Type == tokenizer.NewLine {
			return nil, 0
		}

		if token.Type != prefixTokenType {
			contentTokens = append(contentTokens, token)
			i++
			continue
		}

		runLen := 1
		for i+runLen < len(matchedTokens) && matchedTokens[i+runLen].Type == prefixTokenType {
			runLen++
		}

		if runLen%2 == 0 {
			for j := 0; j < runLen; j++ {
				contentTokens = append(contentTokens, matchedTokens[i+j])
			}
			i += runLen
			continue
		}

		if runLen == 1 {
			closingIndex = i
			break
		}

		for j := 0; j < runLen-1; j++ {
			contentTokens = append(contentTokens, matchedTokens[i+j])
		}
		closingIndex = i + runLen - 1
		break
	}

	if closingIndex == -1 || len(contentTokens) == 0 {
		return nil, 0
	}

	children, err := ParseInline(contentTokens)
	if err != nil || len(children) == 0 {
		return nil, 0
	}

	return &ast.Italic{
		Symbol:   prefixTokenType,
		Children: children,
	}, closingIndex + 1
}
