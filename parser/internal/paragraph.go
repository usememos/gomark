package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type ParagraphParser struct {
	lookahead []BaseParser
}

func NewParagraphParser() *ParagraphParser {
	return &ParagraphParser{}
}

func (p *ParagraphParser) SetLookaheadParsers(parsers []BaseParser) {
	p.lookahead = p.lookahead[:0]
	for _, parser := range parsers {
		if parser == nil || parser == p {
			continue
		}
		if _, ok := parser.(*LineBreakParser); ok {
			continue
		}
		p.lookahead = append(p.lookahead, parser)
	}
}

func (p *ParagraphParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	if len(tokens) == 0 {
		return nil, 0
	}
	if tokens[0].Type == tokenizer.NewLine {
		return nil, 0
	}

	var matchedTokens []*tokenizer.Token
	consumed := 0

	for consumed < len(tokens) {
		line, _, newlineCount := consumeLine(tokens[consumed:])
		lineLen := len(line)

		if isBlankLine(line) {
			consumed += lineLen
			break
		}

		matchedTokens = append(matchedTokens, line...)
		consumed += lineLen

		if newlineCount == 0 {
			break
		}

		if newlineCount > 1 {
			break
		}

		// Single newline: treat as soft break if the next segment is not a new block.
		consumed++
		remainder := tokens[consumed:]
		if len(remainder) == 0 {
			break
		}
		if p.startsNewBlock(remainder) {
			consumed--
			break
		}

		matchedTokens = append(matchedTokens, &tokenizer.Token{Type: tokenizer.NewLine, Value: "\n"})
	}

	if len(matchedTokens) == 0 {
		return nil, 0
	}

	children, err := ParseInline(matchedTokens)
	if err != nil {
		return nil, 0
	}

	return &ast.Paragraph{Children: children}, consumed
}

func (p *ParagraphParser) startsNewBlock(tokens []*tokenizer.Token) bool {
	if len(tokens) == 0 {
		return false
	}
	for _, parser := range p.lookahead {
		node, size := parser.Match(tokens)
		if node != nil && size > 0 {
			return true
		}
	}
	return false
}

func consumeLine(tokens []*tokenizer.Token) ([]*tokenizer.Token, int, int) {
	length := 0
	for length < len(tokens) && tokens[length].Type != tokenizer.NewLine {
		length++
	}

	line := tokens[:length]
	newlineCount := 0
	for cursor := length; cursor < len(tokens) && tokens[cursor].Type == tokenizer.NewLine; cursor++ {
		newlineCount++
	}

	return line, length + newlineCount, newlineCount
}

func isBlankLine(tokens []*tokenizer.Token) bool {
	if len(tokens) == 0 {
		return true
	}
	for _, token := range tokens {
		if token.Type != tokenizer.Space {
			return false
		}
	}
	return true
}
