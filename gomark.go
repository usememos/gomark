// Package gomark provides a fast, extensible markdown parser with zero-configuration usage.
//
// All markdown features are enabled by default, including HTML elements, tables, math,
// highlighting, and more. Simply call Parse() to convert markdown to AST nodes.
package gomark

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/config"
	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
	"github.com/usememos/gomark/renderer/markdown"
)

// Engine coordinates parsing and rendering.
type Engine struct {
	config   *config.ParserConfig
	registry *parser.ParserRegistry
}

// EngineOption configures Engine behaviour.
type EngineOption func(*Engine)

// NewEngine creates an Engine with the given options applied.
func NewEngine(opts ...EngineOption) *Engine {
	cfg := config.DefaultConfig()
	engine := &Engine{
		config:   cfg,
		registry: parser.NewParserRegistry(),
	}
	for _, opt := range opts {
		opt(engine)
	}
	return engine
}

// WithConfig sets the parser configuration.
func WithConfig(cfg *config.ParserConfig) EngineOption {
	return func(e *Engine) {
		e.config = cfg
	}
}

var defaultEngine = NewEngine()

// Parse parses markdown text into an AST document using the default engine.
// This uses DefaultConfig with all features enabled including HTML elements.
func Parse(markdown string) (*ast.Document, error) {
	return defaultEngine.Parse(markdown)
}

// Restore renders the document back to markdown using the default engine.
func Restore(doc *ast.Document) string {
	return defaultEngine.Restore(doc)
}

// NewMemosEngine creates an engine with default configuration and all features enabled.
func NewMemosEngine() *Engine {
	return NewEngine(WithConfig(config.DefaultConfig()))
}

// Parse parses markdown into an AST document.
func (e *Engine) Parse(markdown string) (*ast.Document, error) {
	tokens := tokenizer.Tokenize(markdown)
	return parser.ParseWithRegistry(tokens, e.registry)
}

// Restore renders the AST document to markdown.
func (*Engine) Restore(doc *ast.Document) string {
	if doc == nil {
		return ""
	}
	r := markdown.NewMarkdownRenderer()
	r.RenderDocument(doc)
	return r.String()
}
