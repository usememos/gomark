package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type TableParser struct{}

func NewTableParser() *TableParser {
	return &TableParser{}
}

func (*TableParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	rawRows := tokenizer.Split(tokens, tokenizer.NewLine)
	if len(rawRows) < 3 {
		return nil, 0
	}

	headerTokens := rawRows[0]
	if len(headerTokens) < 3 {
		return nil, 0
	}

	delimiterTokens := rawRows[1]
	if len(delimiterTokens) < 3 {
		return nil, 0
	}

	// Check header.
	headerCells, ok := matchTableCellTokens(headerTokens)
	if !ok || len(headerCells) == 0 {
		return nil, 0
	}

	// Check delimiter.
	delimiterCells, ok := matchTableCellTokens(delimiterTokens)
	if !ok || len(delimiterCells) != len(headerCells) {
		return nil, 0
	}
	delimiter := make([]string, 0, len(delimiterCells))
	for _, cellTokens := range delimiterCells {
		if !isValidDelimiterCell(cellTokens) {
			return nil, 0
		}
		delimiter = append(delimiter, tokenizer.Stringify(cellTokens))
	}

	// Check rows.
	rows := rawRows[2:]
	matchedRows := make([][][]*tokenizer.Token, 0)
	matchedRawRows := make([][]*tokenizer.Token, 0)
	for _, rowTokens := range rows {
		cells, ok := matchTableCellTokens(rowTokens)
		if !ok || len(cells) != len(headerCells) {
			break
		}
		matchedRows = append(matchedRows, cells)
		matchedRawRows = append(matchedRawRows, rowTokens)
	}
	if len(matchedRows) == 0 {
		return nil, 0
	}

	headerNodes := make([]ast.Node, 0, len(headerCells))
	for _, cellTokens := range headerCells {
		if len(cellTokens) == 0 {
			headerNodes = append(headerNodes, &ast.Text{})
			continue
		}
		nodes, err := ParseBlockWithParsers(cellTokens, []BlockParser{NewHeadingParser(), NewParagraphParser()})
		if err != nil {
			return nil, 0
		}
		if len(nodes) != 1 {
			return nil, 0
		}
		headerNodes = append(headerNodes, nodes[0])
	}

	rowsNodes := make([][]ast.Node, 0, len(matchedRows))
	for _, rowCells := range matchedRows {
		rowNodes := make([]ast.Node, 0, len(rowCells))
		for _, cellTokens := range rowCells {
			if len(cellTokens) == 0 {
				rowNodes = append(rowNodes, &ast.Text{})
				continue
			}
			nodes, err := ParseBlockWithParsers(cellTokens, []BlockParser{NewHeadingParser(), NewParagraphParser()})
			if err != nil {
				return nil, 0
			}
			if len(nodes) != 1 {
				return nil, 0
			}
			rowNodes = append(rowNodes, nodes[0])
		}
		rowsNodes = append(rowsNodes, rowNodes)
	}

	size := len(headerTokens) + len(delimiterTokens) + 2
	for _, row := range matchedRawRows {
		size += len(row)
	}
	if len(matchedRawRows) > 0 {
		size += len(matchedRawRows) - 1
	}

	return &ast.Table{
		Header:    headerNodes,
		Delimiter: delimiter,
		Rows:      rowsNodes,
	}, size
}

func matchTableCellTokens(tokens []*tokenizer.Token) ([][]*tokenizer.Token, bool) {
	if len(tokens) == 0 {
		return nil, false
	}

	pipeCount := 0
	for _, token := range tokens {
		if token.Type == tokenizer.Pipe {
			pipeCount++
		}
	}
	if pipeCount == 0 {
		return nil, false
	}

	cells := tokenizer.Split(tokens, tokenizer.Pipe)
	start := 0
	end := len(cells)

	if len(cells) > 0 && len(cells[0]) == 0 {
		start++
	}
	if end > start && len(cells[len(cells)-1]) == 0 {
		end--
	}
	if start >= end {
		return nil, false
	}

	trimmed := make([][]*tokenizer.Token, 0, end-start)
	for _, cell := range cells[start:end] {
		trimmed = append(trimmed, trimTableCellSpaces(cell))
	}

	return trimmed, true
}

func trimTableCellSpaces(tokens []*tokenizer.Token) []*tokenizer.Token {
	start := 0
	end := len(tokens)
	for start < end && isSpaceToken(tokens[start]) {
		start++
	}
	for end > start && isSpaceToken(tokens[end-1]) {
		end--
	}
	return tokens[start:end]
}

func isSpaceToken(token *tokenizer.Token) bool {
	return token.Type == tokenizer.Space || token.Type == tokenizer.MultipleSpaces
}

func isValidDelimiterCell(tokens []*tokenizer.Token) bool {
	if len(tokens) == 0 {
		return false
	}

	hyphenCount := 0
	for index, token := range tokens {
		switch token.Type {
		case tokenizer.Hyphen:
			hyphenCount++
		case tokenizer.Colon:
			if index != 0 && index != len(tokens)-1 {
				return false
			}
		default:
			return false
		}
	}

	return hyphenCount >= 3
}
