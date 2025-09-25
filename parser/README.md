# Parser Package

The `parser` package contains the core parsing logic for converting markdown tokens into AST nodes. It provides a modular, extensible parser architecture with support for various markdown elements and extensions.

## Overview

This package implements a simple, efficient two-stage parsing approach:
1. **Tokenization**: Raw markdown text is converted into tokens (single-pass)
2. **Parsing**: Tokens are processed into AST nodes using specialized parsers

The parser architecture is highly modular, with separate parsers for different markdown elements that can be enabled/disabled based on configuration.

## Design Philosophy

gomark's parser prioritizes **simplicity and performance**:

- **Token-based parsing**: Simple, fast single-character tokenization
- **Stateless parsers**: Each parser is independent and easy to understand
- **Direct approach**: No complex state machines or context tracking
- **Focused functionality**: Implements what's needed without over-engineering

This approach provides excellent performance while maintaining code that's easy to understand and maintain.

## Core Architecture

### Parser Interfaces

```go
type BaseParser interface {
    Match(tokens []*tokenizer.Token) (ast.Node, int)
}

type InlineParser interface {
    BaseParser
}

type BlockParser interface {
    BaseParser
}
```

All parsers implement one of these interfaces, returning:
- `ast.Node`: The parsed AST node (nil if no match)
- `int`: Number of tokens consumed during parsing

### Parser Registry

The `ParserRegistry` manages all available parsers and their configuration:

```go
type ParserRegistry struct {
    blockParsers  map[string]BlockParserFactory
    inlineParsers map[string]InlineParserFactory
    config        *config.ParserConfig
}
```

#### Factory Functions

```go
type BlockParserFactory func() BlockParser
type InlineParserFactory func() InlineParser
```

## Available Parsers

### Block-level Parsers

#### Core Markdown
- **Paragraph**: Text paragraphs
- **Heading**: Headers (H1-H6)
- **CodeBlock**: Fenced and indented code blocks
- **Blockquote**: Quoted text blocks
- **HorizontalRule**: Horizontal dividers
- **LineBreak**: Line breaks

#### List Parsers
- **OrderedListItem**: Numbered list items
- **UnorderedListItem**: Bulleted list items
- **TaskListItem**: Checkbox list items (extension)

#### Extension Parsers
- **Table**: Table structures (extension)
- **MathBlock**: Block math expressions (extension)
- **EmbeddedContent**: Embedded content blocks (extension)

### Inline Parsers

#### Text Formatting
- **Text**: Plain text content
- **Bold**: Bold text (`**text**` or `__text__`)
- **Italic**: Italic text (`*text*` or `_text_`)
- **BoldItalic**: Combined bold-italic formatting
- **Code**: Inline code (`` `code` ``)
- **Strikethrough**: Strikethrough text (`~~text~~`, extension)

#### Links and Media
- **Link**: Manual links (`[text](url)`)
- **AutoLink**: Automatic URL detection (extension)
- **Image**: Images (`![alt](src)`)

#### Extensions
- **Math**: Inline math expressions (`$math$`)
- **Highlight**: Highlighted text (`==text==`)
- **Subscript**: Subscript text (`~text~`)
- **Superscript**: Superscript text (`^text^`)
- **Spoiler**: Spoiler text (`||text||`)
- **Tag**: Tag syntax (`#tag`)
- **ReferencedContent**: Referenced content (`[[ref]]`)
- **EscapingCharacter**: Escaped characters (`\char`)
- **HTMLElement**: HTML elements (`<kbd>`, `<br>`, `<img>`, `<small>`, `<mark>`)

## Parsing Process

### Main Parse Functions

```go
// Parse with default configuration
func Parse(tokens []*tokenizer.Token) (*ast.Document, error)

// Parse with custom configuration
func ParseWithConfig(tokens []*tokenizer.Token, cfg *config.ParserConfig) (*ast.Document, error)

// Parse with specific registry
func ParseWithRegistry(tokens []*tokenizer.Token, registry *ParserRegistry) (*ast.Document, error)
```

### Parse Flow

1. **Tokenization**: Text → Tokens (via tokenizer package)
2. **Block Parsing**: Tokens → Block-level AST nodes
3. **Inline Parsing**: Process inline content within blocks
4. **Document Assembly**: Combine all nodes into final AST

### Error Handling

The parser includes comprehensive error tracking:

```go
type ParseResult struct {
    Errors   []ParseError
    Warnings []ParseWarning
}
```

Errors include:
- **Syntax errors**: Malformed markdown syntax
- **Depth errors**: Excessive nesting beyond configured limits
- **Extension errors**: Extension-specific parsing issues

## Parser Registration

### Automatic Registration

Default parsers are automatically registered based on configuration:

```go
registry := NewParserRegistry(config)
// All appropriate parsers registered based on config.EnableExtensions
```

### Manual Registration

```go
// Register custom block parser
registry.RegisterBlockParser("custom", func() BlockParser {
    return &CustomBlockParser{}
})

// Register custom inline parser
registry.RegisterInlineParser("custom", func() InlineParser {
    return &CustomInlineParser{}
})
```

### Parser Priority

Parsers are executed in priority order. Higher priority parsers are tried first:

```go
type ParserInfo struct {
    Name        string
    Description string
    Priority    int      // Higher = earlier execution
    Extension   string   // Which config extension enables this
}
```

## Tokenizer Subpackage

### Token Types

The tokenizer recognizes various markdown-significant characters:

```go
const (
    Underscore         = "_"
    Asterisk          = "*"
    PoundSign         = "#"
    Backtick          = "`"
    LeftSquareBracket = "["
    // ... and many more
)
```

### Tokenization

```go
tokens := tokenizer.Tokenize(markdown)
```

Converts raw markdown text into a sequence of tokens for parsing.

## Usage Examples

### Basic Parsing (Legacy API)

```go
import (
    "github.com/usememos/gomark/parser"
    "github.com/usememos/gomark/parser/tokenizer"
)

// Tokenize markdown
tokens := tokenizer.Tokenize("# Hello\n\nThis is **bold** text with <kbd>Ctrl</kbd> keys.")

// Parse into AST (all features enabled by default)
doc, err := parser.Parse(tokens)
if err != nil {
    // Handle parsing errors
}
```

### Recommended High-Level API

For most use cases, prefer the high-level API which includes tokenization:

```go
import "github.com/usememos/gomark"

// Simple usage - all features including HTML elements enabled
doc, err := gomark.Parse("Press <kbd>Ctrl</kbd>+<kbd>C</kbd> to copy")
```

### Custom Configuration

```go
import (
    "github.com/usememos/gomark/config"
    "github.com/usememos/gomark/parser"
    "github.com/usememos/gomark/parser/tokenizer"
)

// Create custom config if needed
cfg := config.DefaultConfig().                  // Start with default configuration
    WithMaxDepth(100).                          // Limit nesting depth
    WithMaxFileSize(1024 * 1024)                // 1MB file size limit

// Parse with config
tokens := tokenizer.Tokenize(markdown)
doc, err := parser.ParseWithConfig(tokens, cfg)

// Alternative: Use high-level API with config
import "github.com/usememos/gomark"
engine := gomark.NewEngine(gomark.WithConfig(cfg))
doc, err := engine.Parse(markdown)
```

### Custom Parser Implementation

```go
type CustomInlineParser struct{}

func (p *CustomInlineParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
    // Check if tokens match custom pattern
    if len(tokens) >= 2 && tokens[0].Type == "!!" && tokens[1].Type == "text" {
        // Create custom AST node
        node := &ast.Text{Content: "Custom: " + tokens[1].Content}
        return node, 2 // Consumed 2 tokens
    }
    return nil, 0 // No match
}

// Register the parser
registry.RegisterInlineParser("custom", func() InlineParser {
    return &CustomInlineParser{}
})
```

## File Structure

- `parser.go` - Main parsing functions and interfaces
- `registry.go` - Parser registry and management
- `errors.go` - Error handling and reporting
- `tokenizer/` - Tokenization subpackage
- `internal/` - Internal parser implementations:
  - `paragraph.go`, `heading.go`, `code_block.go` - Block parsers
  - `bold.go`, `italic.go`, `link.go` - Inline parsers
  - Extension parsers for specialized features
- `tests/` - Comprehensive test suite for all parsers

## Extension System

The parser supports a rich extension system where individual features can be enabled/disabled:

- Extensions are controlled via `config.ExtensionConfig`
- Each parser checks if its required extension is enabled
- Disabled parsers are not registered in the registry
- This allows for lightweight, customized parsing profiles

## Performance Benefits

Our token-based approach provides several performance advantages:

- **Single-pass tokenization**: Text is processed once into reusable tokens
- **Minimal state**: Parsers are stateless and don't require complex context
- **Direct field access**: AST nodes use simple struct fields instead of method calls
- **Focused scope**: Only implements features that are actually used

## Intentional Design Choices

Several features common in other markdown parsers are intentionally simplified:

- **No complex tree navigation**: AST nodes don't track parents/siblings (not needed)
- **Single-character tokens**: Multi-character tokenization disabled for simplicity
- **Stateless parsing**: No complex parsing context (simpler and faster)
- **Direct HTML support**: Basic HTML elements without complex attribute parsing

These choices prioritize maintainability and performance over theoretical completeness.