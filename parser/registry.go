package parser

import (
	"fmt"
	"sort"

	"github.com/usememos/gomark/config"
	"github.com/usememos/gomark/parser/internal"
)

// ParserRegistry manages the collection of available parsers.
type ParserRegistry struct {
	blockParsers  map[string]BlockParserFactory
	inlineParsers map[string]InlineParserFactory
	config        *config.ParserConfig
}

// BlockParserFactory creates block parsers.
type BlockParserFactory func() BlockParser

// InlineParserFactory creates inline parsers.
type InlineParserFactory func() InlineParser

// ParserInfo holds metadata about a parser.
type ParserInfo struct {
	Name        string
	Description string
	Priority    int
	Extension   string // Optional: which extension enables this parser
}

// NewParserRegistry creates a new parser registry.
func NewParserRegistry(cfg *config.ParserConfig) *ParserRegistry {
	registry := &ParserRegistry{
		blockParsers:  make(map[string]BlockParserFactory),
		inlineParsers: make(map[string]InlineParserFactory),
		config:        cfg,
	}

	// Register default parsers
	registry.registerDefaultParsers()
	return registry
}

// RegisterBlockParser registers a block parser factory.
func (r *ParserRegistry) RegisterBlockParser(name string, factory BlockParserFactory) {
	r.blockParsers[name] = factory
}

// RegisterInlineParser registers an inline parser factory.
func (r *ParserRegistry) RegisterInlineParser(name string, factory InlineParserFactory) {
	r.inlineParsers[name] = factory
}

// UnregisterBlockParser removes a block parser.
func (r *ParserRegistry) UnregisterBlockParser(name string) {
	delete(r.blockParsers, name)
}

// UnregisterInlineParser removes an inline parser.
func (r *ParserRegistry) UnregisterInlineParser(name string) {
	delete(r.inlineParsers, name)
}

// GetBlockParsers returns enabled block parsers in priority order.
func (r *ParserRegistry) GetBlockParsers() []BlockParser {
	var parsers []BlockParser

	// Core parsers (always enabled)
	coreParsers := []struct {
		name    string
		factory BlockParserFactory
	}{
		{"code_block", func() BlockParser { return internal.NewCodeBlockParser() }},
		{"heading", func() BlockParser { return internal.NewHeadingParser() }},
		{"horizontal_rule", func() BlockParser { return internal.NewHorizontalRuleParser() }},
		{"blockquote", func() BlockParser { return internal.NewBlockquoteParser() }},
		{"ordered_list", func() BlockParser { return internal.NewOrderedListItemParser() }},
		{"unordered_list", func() BlockParser { return internal.NewUnorderedListItemParser() }},
	}

	for _, p := range coreParsers {
		if factory, exists := r.blockParsers[p.name]; exists {
			parsers = append(parsers, factory())
		}
	}

	// Extension-based parsers
	extensionParsers := []struct {
		name      string
		factory   BlockParserFactory
		condition func() bool
	}{
		{"table", func() BlockParser { return internal.NewTableParser() }, func() bool { return r.config.EnableExtensions.Tables }},
		{"task_list", func() BlockParser { return internal.NewTaskListItemParser() }, func() bool { return r.config.EnableExtensions.TaskLists }},
		{"math_block", func() BlockParser { return internal.NewMathBlockParser() }, func() bool { return r.config.EnableExtensions.Math }},
		{"embedded_content", func() BlockParser { return internal.NewEmbeddedContentParser() }, func() bool { return r.config.EnableExtensions.EmbeddedContent }},
	}

	for _, p := range extensionParsers {
		if p.condition() {
			if factory, exists := r.blockParsers[p.name]; exists {
				parsers = append(parsers, factory())
			}
		}
	}

	// Paragraph parser should be last (lowest priority)
	var paragraph BlockParser = internal.NewParagraphParser()
	if factory, exists := r.blockParsers["paragraph"]; exists {
		paragraph = factory()
	}
	parsers = append(parsers, paragraph)

	// Set lookahead parsers for paragraph
	if p, ok := paragraph.(*internal.ParagraphParser); ok {
		// Convert BlockParser slice to BaseParser slice
		baseParsers := make([]internal.BaseParser, len(parsers))
		for i, bp := range parsers {
			baseParsers[i] = bp
		}
		p.SetLookaheadParsers(baseParsers)
	}

	return parsers
}

// GetInlineParsers returns enabled inline parsers in priority order.
func (r *ParserRegistry) GetInlineParsers() []InlineParser {
	var parsers []InlineParser

	// Core parsers (always enabled, in priority order)
	coreParsers := []struct {
		name    string
		factory InlineParserFactory
	}{
		{"escaping_character", func() InlineParser { return internal.NewEscapingCharacterParser() }},
		{"line_break", func() InlineParser { return internal.NewLineBreakParser() }},
		{"image", func() InlineParser { return internal.NewImageParser() }},
		{"link", func() InlineParser { return internal.NewLinkParser() }},
		{"bold_italic", func() InlineParser { return internal.NewBoldItalicParser() }},
		{"bold", func() InlineParser { return internal.NewBoldParser() }},
		{"italic", func() InlineParser { return internal.NewItalicParser() }},
		{"code", func() InlineParser { return internal.NewCodeParser() }},
		{"text", func() InlineParser { return internal.NewTextParser() }}, // Should be last
	}

	// Add core parsers first
	for _, p := range coreParsers {
		if p.name == "text" {
			continue // Handle text parser separately at the end
		}
		if factory, exists := r.inlineParsers[p.name]; exists {
			parsers = append(parsers, factory())
		}
	}

	// Extension-based parsers
	extensionParsers := []struct {
		name      string
		factory   InlineParserFactory
		condition func() bool
	}{
		{"html_element", func() InlineParser { return internal.NewHTMLElementParser() }, func() bool { return r.config.AllowHTML }},
		{"autolink", func() InlineParser { return internal.NewAutoLinkParser() }, func() bool { return r.config.EnableExtensions.Autolinks }},
		{"strikethrough", func() InlineParser { return internal.NewStrikethroughParser() }, func() bool { return r.config.EnableExtensions.Strikethrough }},
		{"highlight", func() InlineParser { return internal.NewHighlightParser() }, func() bool { return r.config.EnableExtensions.Highlighting }},
		{"subscript", func() InlineParser { return internal.NewSubscriptParser() }, func() bool { return r.config.EnableExtensions.Subscript }},
		{"superscript", func() InlineParser { return internal.NewSuperscriptParser() }, func() bool { return r.config.EnableExtensions.Superscript }},
		{"math", func() InlineParser { return internal.NewMathParser() }, func() bool { return r.config.EnableExtensions.Math }},
		{"spoiler", func() InlineParser { return internal.NewSpoilerParser() }, func() bool { return r.config.EnableExtensions.Spoilers }},
		{"referenced_content", func() InlineParser { return internal.NewReferencedContentParser() }, func() bool { return r.config.EnableExtensions.ReferencedContent }},
		{"tag", func() InlineParser { return internal.NewTagParser() }, func() bool { return r.config.EnableExtensions.Tags }},
	}

	for _, p := range extensionParsers {
		if p.condition() {
			if factory, exists := r.inlineParsers[p.name]; exists {
				parsers = append(parsers, factory())
			}
		}
	}

	// Text parser should always be last (lowest priority)
	if factory, exists := r.inlineParsers["text"]; exists {
		parsers = append(parsers, factory())
	}

	return parsers
}

// ListBlockParsers returns the names of all registered block parsers.
func (r *ParserRegistry) ListBlockParsers() []string {
	var names []string
	for name := range r.blockParsers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ListInlineParsers returns the names of all registered inline parsers.
func (r *ParserRegistry) ListInlineParsers() []string {
	var names []string
	for name := range r.inlineParsers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// HasBlockParser checks if a block parser is registered.
func (r *ParserRegistry) HasBlockParser(name string) bool {
	_, exists := r.blockParsers[name]
	return exists
}

// HasInlineParser checks if an inline parser is registered.
func (r *ParserRegistry) HasInlineParser(name string) bool {
	_, exists := r.inlineParsers[name]
	return exists
}

// UpdateConfig updates the registry's configuration.
func (r *ParserRegistry) UpdateConfig(cfg *config.ParserConfig) {
	r.config = cfg
}

// GetConfig returns the current configuration.
func (r *ParserRegistry) GetConfig() *config.ParserConfig {
	return r.config
}

// Clone creates a copy of the registry with the same parsers but a new config.
func (r *ParserRegistry) Clone(cfg *config.ParserConfig) *ParserRegistry {
	newRegistry := &ParserRegistry{
		blockParsers:  make(map[string]BlockParserFactory),
		inlineParsers: make(map[string]InlineParserFactory),
		config:        cfg,
	}

	// Copy all registered parsers
	for name, factory := range r.blockParsers {
		newRegistry.blockParsers[name] = factory
	}
	for name, factory := range r.inlineParsers {
		newRegistry.inlineParsers[name] = factory
	}

	return newRegistry
}

// registerDefaultParsers registers the built-in parsers.
func (r *ParserRegistry) registerDefaultParsers() {
	// Register block parsers
	r.RegisterBlockParser("code_block", func() BlockParser { return internal.NewCodeBlockParser() })
	r.RegisterBlockParser("table", func() BlockParser { return internal.NewTableParser() })
	r.RegisterBlockParser("horizontal_rule", func() BlockParser { return internal.NewHorizontalRuleParser() })
	r.RegisterBlockParser("heading", func() BlockParser { return internal.NewHeadingParser() })
	r.RegisterBlockParser("blockquote", func() BlockParser { return internal.NewBlockquoteParser() })
	r.RegisterBlockParser("ordered_list", func() BlockParser { return internal.NewOrderedListItemParser() })
	r.RegisterBlockParser("task_list", func() BlockParser { return internal.NewTaskListItemParser() })
	r.RegisterBlockParser("unordered_list", func() BlockParser { return internal.NewUnorderedListItemParser() })
	r.RegisterBlockParser("math_block", func() BlockParser { return internal.NewMathBlockParser() })
	r.RegisterBlockParser("embedded_content", func() BlockParser { return internal.NewEmbeddedContentParser() })
	r.RegisterBlockParser("paragraph", func() BlockParser { return internal.NewParagraphParser() })

	// Register inline parsers
	r.RegisterInlineParser("escaping_character", func() InlineParser { return internal.NewEscapingCharacterParser() })
	r.RegisterInlineParser("line_break", func() InlineParser { return internal.NewLineBreakParser() })
	r.RegisterInlineParser("html_element", func() InlineParser { return internal.NewHTMLElementParser() })
	r.RegisterInlineParser("bold_italic", func() InlineParser { return internal.NewBoldItalicParser() })
	r.RegisterInlineParser("image", func() InlineParser { return internal.NewImageParser() })
	r.RegisterInlineParser("link", func() InlineParser { return internal.NewLinkParser() })
	r.RegisterInlineParser("autolink", func() InlineParser { return internal.NewAutoLinkParser() })
	r.RegisterInlineParser("bold", func() InlineParser { return internal.NewBoldParser() })
	r.RegisterInlineParser("italic", func() InlineParser { return internal.NewItalicParser() })
	r.RegisterInlineParser("spoiler", func() InlineParser { return internal.NewSpoilerParser() })
	r.RegisterInlineParser("highlight", func() InlineParser { return internal.NewHighlightParser() })
	r.RegisterInlineParser("code", func() InlineParser { return internal.NewCodeParser() })
	r.RegisterInlineParser("subscript", func() InlineParser { return internal.NewSubscriptParser() })
	r.RegisterInlineParser("superscript", func() InlineParser { return internal.NewSuperscriptParser() })
	r.RegisterInlineParser("math", func() InlineParser { return internal.NewMathParser() })
	r.RegisterInlineParser("referenced_content", func() InlineParser { return internal.NewReferencedContentParser() })
	r.RegisterInlineParser("tag", func() InlineParser { return internal.NewTagParser() })
	r.RegisterInlineParser("strikethrough", func() InlineParser { return internal.NewStrikethroughParser() })
	r.RegisterInlineParser("text", func() InlineParser { return internal.NewTextParser() })
}

// ValidateParsers checks that all required parsers are registered.
func (r *ParserRegistry) ValidateParsers() error {
	required := []string{"paragraph", "text"}

	for _, name := range required {
		if !r.HasBlockParser(name) && !r.HasInlineParser(name) {
			return fmt.Errorf("required parser '%s' is not registered", name)
		}
	}

	return nil
}
