package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

// We use aliases to avoid circular imports
// The actual interfaces are defined in the parent package

// BaseParser defines the basic parsing interface
type BaseParser interface {
	Match(tokens []*tokenizer.Token) (ast.Node, int)
}

// InlineParser represents a parser for inline elements
// This is an alias that matches the parent package interface
type InlineParser = BaseParser

// BlockParser represents a parser for block elements
// This is an alias that matches the parent package interface
type BlockParser = BaseParser