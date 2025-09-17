# AST Package

The `ast` package provides the Abstract Syntax Tree (AST) node definitions for the gomark markdown parser. It defines all the node types and structures used to represent parsed markdown content.

## Overview

This package contains the core AST node definitions that represent the structure of a markdown document after parsing. All nodes implement the `Node` interface and fall into two main categories:

- **Block-level nodes**: Represent structural elements like paragraphs, headings, code blocks
- **Inline nodes**: Represent text formatting within block elements like bold, italic, links

## Design Philosophy

gomark's AST is designed for **simplicity and performance** rather than complex tree manipulation:

- **Direct field access**: `node.Children` instead of method calls
- **Simple interface**: Only `Type()` and `Restore()` methods required
- **Minimal overhead**: No parent/sibling tracking unless needed
- **Focused scope**: Implements what's actually used, not theoretical features

This approach provides excellent performance while maintaining clean, readable code.

## Core Types

### Node Interface

All AST nodes implement the `Node` interface:

```go
type Node interface {
    Type() NodeType
    Restore() string
}
```

### Node Categories

#### Block-level Nodes
- `DocumentNode` - Root document container
- `ParagraphNode` - Text paragraphs
- `HeadingNode` - Headings (H1-H6)
- `CodeBlockNode` - Code blocks
- `BlockquoteNode` - Quoted text blocks
- `ListNode` - List containers
- `OrderedListItemNode` / `UnorderedListItemNode` - List items
- `TaskListItemNode` - Checkbox list items
- `TableNode` - Tables
- `HorizontalRuleNode` - Horizontal rules
- `MathBlockNode` - Math expressions
- `EmbeddedContentNode` - Embedded content

#### Inline Nodes
- `TextNode` - Plain text content
- `BoldNode` - Bold text formatting
- `ItalicNode` - Italic text formatting
- `BoldItalicNode` - Bold italic combination
- `CodeNode` - Inline code
- `LinkNode` - Hyperlinks
- `AutoLinkNode` - Automatic links
- `ImageNode` - Images
- `TagNode` - Custom tags
- `StrikethroughNode` - Strikethrough text
- `HighlightNode` - Highlighted text
- `SubscriptNode` / `SuperscriptNode` - Sub/superscript
- `SpoilerNode` - Spoiler text
- `MathNode` - Inline math
- `EscapingCharacterNode` - Escaped characters
- `ReferencedContentNode` - Referenced content
- `HTMLElementNode` - HTML elements

## Key Structures

### Base Types

```go
type BaseNode struct{}
type BaseBlock struct{ BaseNode }
type BaseInline struct{ BaseNode }
```

All concrete node types embed one of these base structures.

### Document Structure

The `Document` node serves as the root container:

```go
type Document struct {
    BaseBlock
    Children []Node
}
```

### Example Node Implementation

```go
type Paragraph struct {
    BaseBlock
    Children []Node
}

func (*Paragraph) Type() NodeType {
    return ParagraphNode
}

func (n *Paragraph) Restore() string {
    var result string
    for _, child := range n.Children {
        result += child.Restore()
    }
    return result
}
```

## Usage

The AST package is typically used in conjunction with the parser and renderer packages:

1. **Parser** creates AST nodes from markdown tokens
2. **AST nodes** represent the document structure
3. **Renderer** converts AST nodes back to output format

### Creating Nodes

```go
// Create a text node
textNode := &ast.Text{Content: "Hello, world!"}

// Create a paragraph with text
paragraph := &ast.Paragraph{
    Children: []ast.Node{textNode},
}
```

### Traversing the AST

```go
func walkNode(node ast.Node) {
    fmt.Printf("Node type: %s\n", node.Type())

    // Handle nodes with children using direct field access
    switch n := node.(type) {
    case *ast.Paragraph:
        for _, child := range n.Children {
            walkNode(child)
        }
    case *ast.Document:
        for _, child := range n.Children {
            walkNode(child)
        }
    case *ast.Blockquote:
        for _, child := range n.Children {
            walkNode(child)
        }
    }
}
```

### Practical Example

```go
// Parse markdown and examine AST
doc, err := gomark.Parse("# Title\n\nThis is **bold** text.")
if err != nil {
    panic(err)
}

// Walk through the document
for _, child := range doc.Children {
    switch node := child.(type) {
    case *ast.Heading:
        fmt.Printf("Found heading: %s\n", node.Children[0].(*ast.Text).Content)
    case *ast.Paragraph:
        fmt.Printf("Found paragraph with %d children\n", len(node.Children))
    }
}
```

## File Structure

- `ast.go` - Core node type definitions and interface
- `block.go` - Block-level node implementations
- `inline.go` - Inline node implementations
- `document.go` - Document root node
- `util.go` - Utility functions for AST operations