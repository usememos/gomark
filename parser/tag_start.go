package parser

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

const TagMinLen = 6    // #[[x]]
const TagPrefixLen = 3 // #[[

type TagStartParser struct{}

func NewTagStartParser() InlineParser {
	return &TagStartParser{}
}

func (*TagStartParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedTokens := tokenizer.GetFirstLine(tokens)
	if len(matchedTokens) < TagMinLen {
		return nil, 0
	}
	if matchedTokens[0].Type != tokenizer.PoundSign || matchedTokens[1].Type != tokenizer.LeftSquareBracket || matchedTokens[2].Type != tokenizer.LeftSquareBracket {
		return nil, 0
	}

	contentTokens := []*tokenizer.Token{}
	tagEndMatched := false

	for cursor := TagPrefixLen; cursor < len(matchedTokens)-1; cursor++ {
		token, nextToken := matchedTokens[cursor], matchedTokens[cursor+1]

		if token.Type == tokenizer.Space || token.Type == tokenizer.PoundSign {
			break
		}

		if token.Type == tokenizer.RightSquareBracket && nextToken.Type == tokenizer.RightSquareBracket {
			tagEndMatched = true
			break
		}

		contentTokens = append(contentTokens, token)
	}

	if !tagEndMatched || len(contentTokens) == 0 {
		return nil, 0
	}

	return &ast.TagStart{
		Content: tokenizer.Stringify(contentTokens),
	}, len(contentTokens) + 1
}
