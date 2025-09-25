# Config Package

The `config` package provides configuration management for the gomark markdown parser. It defines parser settings, extension controls, and safety options.

## Overview

This package centralizes all configuration options for the markdown parser, allowing fine-grained control over parsing behavior, enabled extensions, and safety features.

**Note**: For most use cases, **no configuration is needed** - gomark works great with defaults!

## Design Philosophy

gomark's configuration system is designed for **extreme simplicity**:

- **Zero-config by default**: All features enabled out of the box for memos
- **Minimal configuration**: Only MaxDepth and MaxFileSize need customization
- **Single configuration**: DefaultConfig with sensible defaults for all use cases
- **All extensions always enabled**: No complex feature toggling needed
- **Immutable configs**: Configuration methods return new instances for safety

## Core Types

### ParserConfig

The main configuration structure that controls parser limits:

```go
type ParserConfig struct {
    MaxDepth    int   // Maximum nesting depth to prevent stack overflow
    MaxFileSize int64 // Maximum file size that can be parsed in bytes
}
```

#### Fields

- **MaxDepth**: Limits maximum nesting depth to prevent stack overflow (default: 200 for memos)
- **MaxFileSize**: Maximum file size that can be parsed in bytes (default: 50MB for memos)

All markdown extensions (tables, math, HTML elements, etc.) are **always enabled** - no configuration needed!

## Configuration

### DefaultConfig()

Returns the standard configuration with sensible defaults for all use cases:

```go
config := config.DefaultConfig()
// MaxDepth: 200 (allow complex nested content)
// MaxFileSize: 50MB (for large documents with embedded content)
// All extensions enabled by default
```

## Configuration Methods

Only two methods are available for customizing limits:

### WithMaxDepth(depth int)

Sets the maximum parsing depth:

```go
config := cfg.WithMaxDepth(50)
```

### WithMaxFileSize(size int64)

Sets the maximum file size:

```go
config := cfg.WithMaxFileSize(1024 * 1024) // 1MB
```

## Usage Examples

### Zero Configuration (Recommended)

```go
import "github.com/usememos/gomark"

// No configuration needed - all features work out of the box!
doc, err := gomark.Parse("Press <kbd>Ctrl</kbd> with **bold** and ==highlighting==")
```

### Basic Configuration Usage

```go
import "github.com/usememos/gomark/config"

// Default configuration - all features enabled
cfg := config.DefaultConfig()
```

### Custom Configuration

```go
// Customize limits if needed
cfg := config.DefaultConfig().
    WithMaxDepth(100).                  // Limit nesting depth
    WithMaxFileSize(1024 * 1024)        // 1MB file size limit
```

### Configuration for Different Use Cases

#### Most Use Cases (Recommended)
```go
cfg := config.DefaultConfig()
// All features enabled with generous limits - works great for most needs!
```

#### Memory-Constrained Environments
```go
cfg := config.DefaultConfig().
    WithMaxDepth(50).                 // Limit nesting depth
    WithMaxFileSize(1024 * 1024)      // 1MB file size limit
```

#### High-Performance Applications
```go
cfg := config.DefaultConfig().
    WithMaxDepth(20).                 // Very shallow nesting
    WithMaxFileSize(100 * 1024)       // 100KB file limit
```

## Security Considerations

- **MaxDepth**: Prevents stack overflow from deeply nested structures
- **MaxFileSize**: Prevents processing of excessively large files
- **HTML Elements**: All supported HTML elements are safe and commonly used

## Integration

This configuration is used throughout the parser:

1. **Engine**: Accepts configuration via `WithConfig()` option
2. **Parser Registry**: Uses config to determine which parsers to enable
3. **Individual Parsers**: Check config for extension-specific behavior

```go
// Engine configuration
engine := gomark.NewEngine(
    gomark.WithConfig(cfg),
)

// Direct parser configuration
registry := parser.NewParserRegistry(cfg)
```


## Extension Impact

Understanding the performance and complexity impact of each extension:

| Extension | Performance Impact | Complexity | Recommended For |
|-----------|-------------------|------------|-----------------|
| `tables` | Low | Medium | Documentation, GitHub-style |
| `math` | Medium | High | Technical content, academia |
| `autolinks` | Low | Low | All use cases |
| `tasklists` | Low | Low | Project management, todos |
| `highlighting` | Low | Low | Note-taking, emphasis |
| `spoilers` | Low | Low | Forums, entertainment |
| `tags` | Low | Low | Note-taking, organization |