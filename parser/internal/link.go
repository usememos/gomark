package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type LinkParser struct{}

func NewLinkParser() *LinkParser {
	return &LinkParser{}
}

// findMatchingRightSquareBracket finds the index of the closing ] that matches the opening [
// It handles nested brackets by tracking bracket depth
func findMatchingRightSquareBracket(tokens []*tokenizer.Token) int {
	depth := 0
	for i, token := range tokens {
		// Skip escaped brackets
		if i > 0 && tokens[i-1].Type == tokenizer.Backslash {
			continue
		}

		if token.Type == tokenizer.LeftSquareBracket {
			depth++
		} else if token.Type == tokenizer.RightSquareBracket {
			if depth == 0 {
				return i
			}
			depth--
		}
	}
	return -1
}

func (*LinkParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	matchedTokens := tokenizer.GetFirstLine(tokens)
	if len(matchedTokens) < 5 {
		return nil, 0
	}
	if matchedTokens[0].Type != tokenizer.LeftSquareBracket {
		return nil, 0
	}

	rightSquareBracketIndex := findMatchingRightSquareBracket(matchedTokens[1:])
	if rightSquareBracketIndex == -1 {
		return nil, 0
	}
	contentTokens := matchedTokens[1 : rightSquareBracketIndex+1]

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

	contentNodes, err := ParseInlineWithParsers(contentTokens, []InlineParser{NewImageParser(), NewEscapingCharacterParser(), NewTextParser()})
	if err != nil {
		return nil, 0
	}
	return &ast.Link{
		Content: contentNodes,
		URL:     tokenizer.Stringify(urlTokens),
	}, 4 + len(contentTokens) + len(urlTokens)
}
