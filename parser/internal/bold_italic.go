package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type BoldItalicParser struct{}

func NewBoldItalicParser() InlineParser {
	return &BoldItalicParser{}
}

func (*BoldItalicParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedTokens := tokenizer.GetFirstLine(tokens)
	if len(matchedTokens) < 7 {
		return nil, 0
	}

	if matchedTokens[0].Type != matchedTokens[1].Type || matchedTokens[0].Type != matchedTokens[2].Type {
		return nil, 0
	}

	symbol := matchedTokens[0].Type
	if symbol != tokenizer.Asterisk && symbol != tokenizer.Underscore {
		return nil, 0
	}

	cursor := 3
	for cursor <= len(matchedTokens)-3 {
		t1, t2, t3 := matchedTokens[cursor], matchedTokens[cursor+1], matchedTokens[cursor+2]
		if t1.Type == tokenizer.NewLine || t2.Type == tokenizer.NewLine || t3.Type == tokenizer.NewLine {
			return nil, 0
		}
		if t1.Type == symbol && t2.Type == symbol && t3.Type == symbol {
			contentTokens := matchedTokens[3:cursor]
			if len(contentTokens) == 0 {
				return nil, 0
			}

			children, err := ParseInline(contentTokens)
			if err != nil || len(children) == 0 {
				return nil, 0
			}

			return &ast.BoldItalic{
				Symbol:   symbol,
				Children: children,
			}, cursor + 3
		}
		cursor++
	}

	return nil, 0
}
