package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type BlockquoteParser struct{}

func NewBlockquoteParser() *BlockquoteParser {
	return &BlockquoteParser{}
}

func (*BlockquoteParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	if len(tokens) == 0 {
		return nil, 0
	}

	inner := make([]*tokenizer.Token, 0, len(tokens))
	consumed := 0
	lines := 0

	for consumed < len(tokens) {
		line, lineSize, newlineCount := consumeLine(tokens[consumed:])
		if len(line) == 0 {
			break
		}

		trimmedIndex := 0
		for trimmedIndex < len(line) && trimmedIndex < 3 && line[trimmedIndex].Type == tokenizer.Space {
			trimmedIndex++
		}
		if trimmedIndex >= len(line) || line[trimmedIndex].Type != tokenizer.GreaterThan {
			break
		}
		trimmedIndex++
		if trimmedIndex < len(line) && line[trimmedIndex].Type == tokenizer.Space {
			trimmedIndex++
		}

		content := line[trimmedIndex:]
		inner = append(inner, content...)
		consumed += lineSize
		lines++

		if newlineCount == 0 {
			break
		}

		for i := 0; i < newlineCount; i++ {
			inner = append(inner, &tokenizer.Token{Type: tokenizer.NewLine, Value: "\n"})
		}
	}

	if lines == 0 {
		return nil, 0
	}

	for len(inner) > 0 && inner[len(inner)-1].Type == tokenizer.NewLine {
		inner = inner[:len(inner)-1]
	}

	children, err := ParseBlock(inner)
	if err != nil {
		return nil, 0
	}

	return &ast.Blockquote{Children: children}, consumed
}
