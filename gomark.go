// Package gomark exposes high-level parsing and rendering helpers.
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
		registry: parser.NewParserRegistry(cfg),
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
		e.registry = parser.NewParserRegistry(cfg)
	}
}

// WithStrictMode enables or disables strict parsing mode.
func WithStrictMode(strict bool) EngineOption {
	return func(e *Engine) {
		e.config = e.config.WithStrictMode(strict)
		e.registry.UpdateConfig(e.config)
	}
}

// WithExtension enables or disables a specific extension.
func WithExtension(name string, enabled bool) EngineOption {
	return func(e *Engine) {
		e.config = e.config.WithExtension(name, enabled)
		e.registry.UpdateConfig(e.config)
	}
}

var defaultEngine = NewEngine()

// Parse parses markdown text into an AST document using the default engine.
func Parse(markdown string) (*ast.Document, error) {
	return defaultEngine.Parse(markdown)
}

// Restore renders the document back to markdown using the default engine.
func Restore(doc *ast.Document) string {
	return defaultEngine.Restore(doc)
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
