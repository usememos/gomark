# Renderer Package

The `renderer` package provides output rendering capabilities for gomark AST nodes. It includes multiple renderer implementations that can convert parsed AST back into various output formats.

## Overview

The renderer package implements a simple, efficient rendering system that can convert AST nodes into different output formats:

- **Markdown Renderer**: Converts AST back to markdown text (roundtrip support)
- **HTML Renderer**: Generates HTML output from AST
- **String Renderer**: Extracts plain text content from AST

All renderers implement a common interface, making it easy to add new output formats.

## Design Philosophy

gomark's renderers are designed for **simplicity and reliability**:

- **Direct AST traversal**: Simple recursive rendering without complex visitor patterns
- **Buffer-based output**: Efficient string building with `bytes.Buffer`
- **Format-specific optimization**: Each renderer optimized for its target format
- **Roundtrip accuracy**: Markdown renderer preserves original formatting where possible

Recent improvements include:
- ✅ **HTML Elements Support**: Added rendering for `<kbd>`, `<br>`, `<img>`, `<small>`, `<mark>`
- ✅ **Enhanced HTML Renderer**: Proper attribute handling and self-closing tag support
- ✅ **Security**: HTML attribute escaping and content sanitization
- ✅ **Fixed blockquote blank lines** (GitHub issue #19)

## Core Architecture

### Renderer Interface

All renderers implement the `Renderer` interface:

```go
type Renderer interface {
    // RenderDocument renders the entire document and returns the result as a string.
    RenderDocument(doc *ast.Document) string

    // RenderNode renders a single node (for advanced use cases).
    RenderNode(node ast.Node)

    // String returns the accumulated output as a string.
    String() string

    // Reset clears the internal buffer to allow reuse.
    Reset()
}
```

### Renderer Types

```go
type RendererType string

const (
    HTMLRenderer     RendererType = "html"
    MarkdownRenderer RendererType = "markdown"
    StringRenderer   RendererType = "string"
)
```

### Factory Pattern

The package includes a factory for creating renderers:

```go
type RendererFactory struct{}

func NewRendererFactory() *RendererFactory
func (f *RendererFactory) CreateRenderer(rendererType RendererType) Renderer
```

## Available Renderers

### Markdown Renderer

**Package**: `github.com/usememos/gomark/renderer/markdown`

Converts AST nodes back to markdown format, enabling roundtrip parsing and rendering.

```go
import markdownrenderer "github.com/usememos/gomark/renderer/markdown"

renderer := markdownrenderer.NewMarkdownRenderer()
markdown := renderer.RenderDocument(doc)
```

**Features**:
- **Roundtrip Support**: Parse markdown → AST → markdown maintains original structure
- **Preserve Formatting**: Maintains original markdown syntax choices where possible
- **Extension Support**: Handles all gomark extensions correctly
- **Efficient**: Leverages AST node `Restore()` methods for optimal performance

**Use Cases**:
- Document editing and transformation
- Markdown processing pipelines
- Content validation and normalization
- AST manipulation with markdown output

### HTML Renderer

**Package**: `github.com/usememos/gomark/renderer/html`

Generates semantic HTML from AST nodes.

```go
import htmlrenderer "github.com/usememos/gomark/renderer/html"

renderer := htmlrenderer.NewHTMLRenderer()
html := renderer.RenderDocument(doc)
```

**Features**:
- **Semantic HTML**: Generates proper HTML5 semantic elements
- **Extension Support**: Handles tables, math, highlighting, and other extensions
- **Context Aware**: Maintains rendering context for nested elements
- **Customizable**: Extensible for custom HTML generation needs

**HTML Element Mapping**:
- `# Heading` → `<h1>Heading</h1>`
- `**Bold**` → `<strong>Bold</strong>`
- `*Italic*` → `<em>Italic</em>`
- `[Link](url)` → `<a href="url">Link</a>`
- `` `code` `` → `<code>code</code>`
- Code blocks → `<pre><code>...</code></pre>`
- Tables → `<table><tr><td>...</td></tr></table>`

**Use Cases**:
- Web page generation
- Email content creation
- Documentation sites
- Rich text display

### String Renderer

**Package**: `github.com/usememos/gomark/renderer/string`

Extracts plain text content from AST, removing all markdown formatting.

```go
import stringrenderer "github.com/usememos/gomark/renderer/string"

renderer := stringrenderer.NewStringRenderer()
text := renderer.RenderDocument(doc)
```

**Features**:
- **Plain Text Output**: Strips all formatting, keeping only text content
- **Structure Preservation**: Maintains paragraph breaks and spacing
- **Link Handling**: Extracts link text and URLs appropriately
- **Table Support**: Formats table content as readable text

**Use Cases**:
- Search indexing
- Content analysis
- Text summarization
- Preview generation
- Accessibility applications

## Usage Examples

### Basic Rendering

```go
import (
    "github.com/usememos/gomark/renderer"
    markdownrenderer "github.com/usememos/gomark/renderer/markdown"
    htmlrenderer "github.com/usememos/gomark/renderer/html"
)

// Create document AST (from parser)
doc, _ := parser.Parse(tokens)

// Render to markdown
mdRenderer := markdownrenderer.NewMarkdownRenderer()
markdown := mdRenderer.RenderDocument(doc)

// Render to HTML
htmlRenderer := htmlrenderer.NewHTMLRenderer()
html := htmlRenderer.RenderDocument(doc)
```

### Factory-based Creation

```go
factory := renderer.NewRendererFactory()

// Create different renderers
htmlRenderer := factory.CreateRenderer(renderer.HTMLRenderer)
markdownRenderer := factory.CreateRenderer(renderer.MarkdownRenderer)
stringRenderer := factory.CreateRenderer(renderer.StringRenderer)

// Use them
html := htmlRenderer.RenderDocument(doc)
markdown := markdownRenderer.RenderDocument(doc)
text := stringRenderer.RenderDocument(doc)
```

### Node-by-Node Rendering

```go
renderer := htmlrenderer.NewHTMLRenderer()

// Render individual nodes
for _, node := range doc.Children {
    renderer.RenderNode(node)
}

// Get accumulated output
output := renderer.String()
```

### Renderer Reuse

```go
renderer := markdownrenderer.NewMarkdownRenderer()

// Render multiple documents
for _, doc := range documents {
    output := renderer.RenderDocument(doc)
    processOutput(output)
    renderer.Reset() // Clear buffer for next document
}
```

## Custom Renderer Implementation

To create a custom renderer, implement the `Renderer` interface:

```go
type CustomRenderer struct {
    output *bytes.Buffer
}

func NewCustomRenderer() *CustomRenderer {
    return &CustomRenderer{
        output: new(bytes.Buffer),
    }
}

func (r *CustomRenderer) RenderDocument(doc *ast.Document) string {
    r.Reset()
    if doc != nil {
        for _, child := range doc.Children {
            r.RenderNode(child)
        }
    }
    return r.String()
}

func (r *CustomRenderer) RenderNode(node ast.Node) {
    switch n := node.(type) {
    case *ast.Paragraph:
        r.output.WriteString("<p>")
        for _, child := range n.Children {
            r.RenderNode(child)
        }
        r.output.WriteString("</p>\n")
    case *ast.Text:
        r.output.WriteString(n.Content)
    // Handle other node types...
    default:
        // Fallback to node's Restore() method
        r.output.WriteString(n.Restore())
    }
}

func (r *CustomRenderer) String() string {
    return r.output.String()
}

func (r *CustomRenderer) Reset() {
    r.output.Reset()
}
```

## Performance Considerations

### Buffer Management
- All renderers use `bytes.Buffer` for efficient string building
- Call `Reset()` between documents to reuse renderer instances
- Avoid frequent string concatenation in custom renderers

### Node Processing
- Renderers process nodes recursively
- For large documents, consider iterative approaches for deep nesting
- Use type switches for efficient node type handling

### Memory Usage
- Renderers accumulate output in memory
- For very large documents, consider streaming approaches
- Reset buffers promptly to free memory

## File Structure

- `renderer.go` - Core renderer interface and factory
- `markdown/` - Markdown renderer implementation
  - `markdown.go` - Markdown rendering logic
  - `markdown_test.go` - Test suite
- `html/` - HTML renderer implementation
  - `html.go` - HTML rendering logic
  - `html_test.go` - Test suite
- `string/` - String renderer implementation
  - `string.go` - String rendering logic
  - `string_test.go` - Test suite

## Integration

The renderer package integrates with the main gomark engine:

```go
// Engine uses markdown renderer by default
engine := gomark.NewEngine()
doc, _ := engine.Parse(markdown)
output := engine.Restore(doc) // Uses markdown renderer

// Custom rendering
renderer := htmlrenderer.NewHTMLRenderer()
html := renderer.RenderDocument(doc)
```

This design allows for flexible output generation while maintaining clean separation between parsing and rendering concerns.

## Recent Improvements

- **✅ Fixed blockquote blank lines**: Resolved GitHub issue #19 where blank lines in blockquotes weren't properly rendered
- **✅ Improved error handling**: Better handling of edge cases in all renderers
- **✅ Performance optimization**: Buffer reuse and efficient string building
- **✅ Better roundtrip accuracy**: Markdown renderer more faithfully preserves original formatting

## Extending Renderers

Adding a new renderer is straightforward:

```go
type CustomRenderer struct {
    output *bytes.Buffer
}

func NewCustomRenderer() *CustomRenderer {
    return &CustomRenderer{output: new(bytes.Buffer)}
}

func (r *CustomRenderer) RenderDocument(doc *ast.Document) string {
    r.Reset()
    for _, child := range doc.Children {
        r.RenderNode(child)
    }
    return r.String()
}

func (r *CustomRenderer) RenderNode(node ast.Node) {
    switch n := node.(type) {
    case *ast.Paragraph:
        // Custom paragraph rendering
        for _, child := range n.Children {
            r.RenderNode(child)
        }
    case *ast.Text:
        r.output.WriteString(n.Content)
    default:
        // Fallback to node's built-in restoration
        r.output.WriteString(n.Restore())
    }
}

func (r *CustomRenderer) String() string { return r.output.String() }
func (r *CustomRenderer) Reset() { r.output.Reset() }
```