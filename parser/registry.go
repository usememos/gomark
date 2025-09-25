package parser

import (
	"fmt"
	"sort"

	"github.com/usememos/gomark/parser/internal"
)

// ParserRegistry manages the collection of available parsers.
type ParserRegistry struct {
	blockParsers  map[string]BlockParserFactory
	inlineParsers map[string]InlineParserFactory
}

// BlockParserFactory creates block parsers.
type BlockParserFactory func() BlockParser

// InlineParserFactory creates inline parsers.
type InlineParserFactory func() InlineParser


// NewParserRegistry creates a new parser registry with all parsers registered.
func NewParserRegistry() *ParserRegistry {
	registry := &ParserRegistry{
		blockParsers:  make(map[string]BlockParserFactory),
		inlineParsers: make(map[string]InlineParserFactory),
	}

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

// GetBlockParsers returns all block parsers in priority order.
func (r *ParserRegistry) GetBlockParsers() []BlockParser {
	var parsers []BlockParser

	// All block parsers in priority order (paragraph must be last)
	parserNames := []string{
		"code_block",
		"heading",
		"horizontal_rule",
		"blockquote",
		"ordered_list",
		"unordered_list",
		"table",
		"task_list",
		"math_block",
		"embedded_content",
	}

	// Add all parsers except paragraph
	for _, name := range parserNames {
		if factory, exists := r.blockParsers[name]; exists {
			parsers = append(parsers, factory())
		}
	}

	// Paragraph parser must be last (lowest priority)
	var paragraph BlockParser = internal.NewParagraphParser()
	if factory, exists := r.blockParsers["paragraph"]; exists {
		paragraph = factory()
	}
	parsers = append(parsers, paragraph)

	// Set lookahead parsers for paragraph
	if p, ok := paragraph.(*internal.ParagraphParser); ok {
		baseParsers := make([]internal.BaseParser, len(parsers))
		for i, bp := range parsers {
			baseParsers[i] = bp
		}
		p.SetLookaheadParsers(baseParsers)
	}

	return parsers
}

// GetInlineParsers returns all inline parsers in priority order.
func (r *ParserRegistry) GetInlineParsers() []InlineParser {
	var parsers []InlineParser

	// All inline parsers in priority order (text must be last)
	parserNames := []string{
		"escaping_character",
		"line_break",
		"html_element",
		"image",
		"link",
		"autolink",
		"bold_italic",
		"bold",
		"italic",
		"strikethrough",
		"highlight",
		"subscript",
		"superscript",
		"math",
		"spoiler",
		"referenced_content",
		"tag",
		"code",
	}

	// Add all parsers except text
	for _, name := range parserNames {
		if factory, exists := r.inlineParsers[name]; exists {
			parsers = append(parsers, factory())
		}
	}

	// Text parser must be last (lowest priority)
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
