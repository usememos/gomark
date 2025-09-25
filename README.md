# gomark

A fast, extensible, and well-structured markdown parser for Go, optimized for simplicity and performance.

## Features

- 🚀 **Fast & Lightweight**: Simple token-based parsing with excellent performance
- 🧩 **Extensible**: Modular architecture with configurable extensions
- 📝 **Standards Compliant**: Supports CommonMark with useful extensions
- 🔄 **Roundtrip Support**: Parse → modify → restore back to markdown
- 🛡️ **Type Safe**: Clean Go interfaces with comprehensive error handling
- 🧪 **Well Tested**: Comprehensive test suite with high coverage

## Supported Markdown Features

### Core Markdown
- **Text formatting**: Bold, italic, bold-italic, strikethrough
- **Code**: Inline code and fenced code blocks
- **Links & Images**: Manual links, auto-links, and images
- **Lists**: Ordered, unordered, and task lists
- **Structure**: Headings, paragraphs, blockquotes, horizontal rules

### Extensions
- **Tables**: GitHub-flavored table syntax
- **Math**: Inline (`$math$`) and block (`$$math$$`) expressions
- **Highlighting**: `==highlighted text==`
- **Subscript/Superscript**: `~sub~` and `^super^`
- **Spoilers**: `||spoiler text||`
- **Tags**: `#hashtag` syntax
- **Referenced Content**: `[[wiki-style]]` links
- **Embedded Content**: `![[embeds]]`
- **HTML Elements**: `<kbd>`, `<br>`, `<img>`, `<small>`, `<mark>`, and more

## Quick Start

gomark is designed for **zero-configuration usage** - all features are enabled by default:

```go
package main

import (
    "fmt"
    "github.com/usememos/gomark"
    "github.com/usememos/gomark/renderer/html"
)

func main() {
    // All features enabled by default - no configuration needed!
    markdown := `# Hello Memos!

**Bold** and *italic* text with ==highlighting==.

Press <kbd>Ctrl</kbd>+<kbd>C</kbd> to copy.

Math: $E = mc^2$ and task lists:
- [x] HTML elements supported
- [ ] Even more features coming

#gomark ||works great||!`

    // Parse with all features enabled
    doc, err := gomark.Parse(markdown)
    if err != nil {
        panic(err)
    }

    // Render to HTML
    renderer := html.NewHTMLRenderer()
    html := renderer.RenderDocument(doc)
    fmt.Println(html)

    // Or restore back to markdown
    restored := gomark.Restore(doc)
    fmt.Println(restored)
}
```

## Advanced Usage

### Custom Configuration

While gomark works great with zero configuration, you can customize it if needed:

```go
import (
    "github.com/usememos/gomark"
    "github.com/usememos/gomark/config"
)

// Customize limits if needed
engine := gomark.NewEngine(gomark.WithConfig(
    config.DefaultConfig().WithMaxDepth(100).WithMaxFileSize(1024 * 1024), // 1MB limit
))

doc, err := engine.Parse(markdown)
```

### Multiple Output Formats

```go
import (
    "github.com/usememos/gomark/renderer"
    "github.com/usememos/gomark/renderer/html"
    "github.com/usememos/gomark/renderer/string"
)

// Parse once, render to multiple formats
doc, _ := gomark.Parse(markdown)

// Render to HTML
htmlRenderer := html.NewHTMLRenderer()
htmlOutput := htmlRenderer.RenderDocument(doc)

// Extract plain text
stringRenderer := string.NewStringRenderer()
textOutput := stringRenderer.RenderDocument(doc)

// Render back to markdown (roundtrip)
markdownOutput := gomark.Restore(doc)
```

### Available Configurations

```go
import "github.com/usememos/gomark/config"

// Default configuration: all features enabled with generous limits
cfg := config.DefaultConfig()

// Custom limits if needed
cfg := config.DefaultConfig().
    WithMaxDepth(50).                   // Limit nesting depth
    WithMaxFileSize(1024 * 1024)        // 1MB file size limit
```

## Architecture

gomark follows a clean, modular architecture:

```
gomark/
├── ast/              # Abstract Syntax Tree definitions
├── config/           # Configuration management
├── parser/           # Parsing logic and tokenization
│   ├── internal/     # Parser implementations
│   ├── tokenizer/    # Tokenization engine
│   └── tests/        # Comprehensive parser tests
├── renderer/         # Output rendering
│   ├── html/         # HTML renderer
│   ├── markdown/     # Markdown renderer (roundtrip)
│   └── string/       # Plain text renderer
└── gomark.go         # Main engine and public API
```

### Design Principles

1. **Simplicity over Complexity**: Clean, understandable code
2. **Performance over Features**: Fast parsing with minimal overhead
3. **Extensibility**: Easy to add new features without breaking existing code
4. **Type Safety**: Proper Go interfaces and error handling
5. **Testing**: Comprehensive test coverage for reliability

## Performance

gomark is designed for performance:

- **Token-based parsing**: Efficient single-pass tokenization
- **Minimal allocations**: Reuses buffers and objects where possible
- **Simple AST**: Direct field access instead of complex tree operations
- **Focused scope**: Implements what's needed, not theoretical features

## Extensions

All extensions can be configured via the config package:

```go
cfg := config.DefaultConfig()

// Disable specific extensions
cfg = cfg.WithExtension("tables", false)
cfg = cfg.WithExtension("math", false)

// Enable strict parsing
cfg = cfg.WithStrictMode(true)

// Enable safe mode (sanitizes content)
cfg = cfg.WithSafeMode(true)
```

### Available Extensions

| Extension | Default | Description |
|-----------|---------|-------------|
| `tables` | ✅ | GitHub-flavored table syntax |
| `strikethrough` | ✅ | `~~strikethrough~~` text |
| `autolinks` | ✅ | Automatic URL detection |
| `tasklists` | ✅ | `- [ ]` and `- [x]` checkboxes |
| `math` | ✅ | `$inline$` and `$$block$$` math |
| `highlighting` | ✅ | `==highlighted==` text |
| `subscript` | ✅ | `~subscript~` text |
| `superscript` | ✅ | `^superscript^` text |
| `spoilers` | ✅ | `\|\|spoiler\|\|` text |
| `embedded` | ✅ | `![[embedded]]` content |
| `references` | ✅ | `[[referenced]]` content |
| `tags` | ✅ | `#hashtag` syntax |

## Recent Improvements

- ✅ **Phase 1 HTML Elements**: Added support for `<kbd>`, `<br>`, `<img>`, `<small>`, `<mark>`
- ✅ **Simplified Configuration**: Zero-config usage with sensible defaults
- ✅ **Enhanced HTML Parsing**: Proper attribute handling and self-closing tag support
- ✅ **Fixed blockquote blank lines** (GitHub issue #19)
- ✅ **Refactored to modular architecture**
- ✅ **Improved package organization** with public APIs
- ✅ **Enhanced test coverage** with 26+ HTML element test cases

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Submit a pull request

## License

This project is part of the [Memos](https://github.com/usememos/memos) ecosystem.

## Inspiration

Designed for simplicity, performance, and ease of use. gomark focuses on practical markdown parsing with clean, maintainable code.