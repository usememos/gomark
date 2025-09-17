# gomark Architecture

This document explains the architectural decisions and design philosophy behind gomark.

## Design Philosophy

gomark is built on the principle of **pragmatic simplicity**:

> "Solve real problems efficiently without over-engineering"

### Core Principles

1. **Simplicity over Complexity**: Choose the simplest solution that works
2. **Performance over Features**: Fast, reliable parsing over theoretical completeness
3. **Maintainability over Flexibility**: Code that's easy to understand and modify
4. **Real Needs over Theoretical Needs**: Implement what's actually used
5. **Direct Solutions**: Avoid layers of abstraction when direct approaches work

## Architectural Decisions

### 1. Token-Based Parsing ✅

**Decision**: Use single-pass tokenization followed by token-based parsing

**Rationale**:
- **Performance**: Single-pass tokenization is very fast
- **Simplicity**: Tokens are easy to work with and debug
- **Reusability**: Tokens can be reused by multiple parsers
- **Memory Efficiency**: Tokens reference original string data

**Alternative Considered**: Text-based parsing (like goldmark)
**Why Rejected**: Added complexity without clear benefits for our use cases

### 2. Simple AST Interface ✅

**Decision**: Use minimal `Node` interface with direct field access

```go
type Node interface {
    Type() NodeType
    Restore() string
}
```

**Rationale**:
- **Performance**: Direct field access (`node.Children`) is faster than method calls
- **Simplicity**: Easy to understand and work with
- **Focused**: Only implements what's actually needed
- **Memory Efficient**: No overhead for unused tree navigation features

**Alternative Considered**: Complex tree interface (like goldmark)
**Why Rejected**: Analysis showed no actual usage of tree navigation in our codebase

### 3. Stateless Parsers ✅

**Decision**: Each parser is independent and stateless

**Rationale**:
- **Simplicity**: No complex context management
- **Debuggability**: Easy to test individual parsers
- **Performance**: No context overhead
- **Maintainability**: Clear separation of concerns

**Alternative Considered**: Context-heavy parsing
**Why Rejected**: Added complexity without clear benefits

### 4. String-Based Node Types ✅

**Decision**: Use `NodeType string` constants

```go
type NodeType string
const ParagraphNode NodeType = "PARAGRAPH"
```

**Rationale**:
- **Debuggability**: Easy to inspect and debug
- **Simplicity**: No complex type hierarchies
- **Extensibility**: Easy to add new types
- **JSON-Friendly**: Serializes naturally

**Alternative Considered**: Interface-based type system
**Why Rejected**: Unnecessary complexity for our needs

### 5. Configuration-Based Extensions ✅

**Decision**: Use configuration to enable/disable features

**Rationale**:
- **Performance**: Disabled features have zero overhead
- **Flexibility**: Easy to customize for different use cases
- **Maintainability**: Clear feature boundaries
- **User-Friendly**: Simple API for configuration

### 6. Buffer-Based Rendering ✅

**Decision**: Use `bytes.Buffer` for output accumulation

**Rationale**:
- **Performance**: Efficient string building
- **Memory**: Reusable buffers
- **Simplicity**: Standard Go pattern
- **Flexibility**: Easy to extend

## Package Organization

### Public vs Internal

**Public Packages** (goldmark-style):
```
├── ast/              # AST definitions - users need access
├── config/           # Configuration - users need to configure
├── parser/           # Parser interfaces - users may extend
├── renderer/         # Renderer interfaces - users may extend
```

**Internal Implementation**:
```
└── parser/internal/  # Parser implementations - users don't need access
```

**Rationale**:
- Public APIs allow extensibility where it matters
- Internal packages keep implementation details hidden
- Follows goldmark patterns for familiarity

## Performance Optimizations

### 1. Minimal Allocations
- Reuse token slices where possible
- Buffer pooling in renderers
- Direct field access instead of method calls

### 2. Single-Pass Processing
- Tokenization is single-pass
- No multiple traversals of input text
- Direct token-to-AST conversion

### 3. Focused Features
- Only implement actually-used functionality
- No complex tree operations unless needed
- Disable unused extensions for zero overhead

## Intentional Limitations

These are **conscious decisions**, not oversights:

### 1. HTML Attributes
**Current**: Basic HTML tags without attributes
**Rationale**: Complex attribute parsing adds significant complexity for minimal benefit

### 2. Multi-Character Tokens
**Current**: Single-character tokenization
**Rationale**: Works for all supported markdown features, simpler implementation

### 3. Complex Tree Navigation
**Current**: Direct field access only
**Rationale**: No actual usage found in codebase analysis

### 4. Parsing Context
**Current**: Stateless parsers
**Rationale**: Sufficient for current feature set, much simpler

## Recent Improvements

### Fixed Blockquote Blank Lines (GitHub Issue #19)
**Problem**: Blank lines in blockquotes weren't rendered correctly
**Solution**: Enhanced `Blockquote.Restore()` to handle `LineBreak` nodes properly
**Result**: Perfect preservation of blank lines in blockquotes

### Package Refactoring
**Problem**: Everything was in `internal/` packages
**Solution**: Moved key packages to public for extensibility
**Result**: goldmark-style architecture with better extensibility

## Comparison with goldmark

| Aspect | goldmark | gomark |
|--------|----------|--------|
| **Complexity** | High | Low |
| **Performance** | Good | Excellent |
| **Extensibility** | Very High | Moderate |
| **Maintainability** | Moderate | High |
| **Learning Curve** | Steep | Gentle |
| **Feature Set** | Comprehensive | Focused |

## When to Choose gomark

✅ **Choose gomark when**:
- You need fast, reliable markdown parsing
- You want simple, maintainable code
- You're building applications, not markdown libraries
- You need good performance with moderate extensibility

❌ **Choose goldmark when**:
- You need maximum extensibility
- You're building a markdown processing library
- You need complex AST transformations
- You need full CommonMark compliance edge cases

## Future Evolution

gomark is designed to evolve pragmatically:

1. **Add features only when needed**: No speculative features
2. **Maintain simplicity**: New features shouldn't complicate existing code
3. **Performance first**: New features shouldn't hurt performance
4. **Backward compatibility**: Changes should be additive

### Potential Future Additions

**Only if there's demonstrated need**:
- AST walking API (if users request it)
- More output formats (if users request them)
- Advanced HTML attributes (if simple approach proves insufficient)
- Text-based parsing (if token-based proves limiting)

## Conclusion

gomark represents a **pragmatic approach** to markdown parsing:

- **Goldmark-inspired architecture** for familiarity and extensibility
- **Performance-focused implementation** for real-world applications
- **Simple, maintainable code** that developers can understand and modify
- **Focused feature set** that solves real problems without over-engineering

This approach delivers excellent performance and maintainability while providing enough extensibility for most real-world use cases.