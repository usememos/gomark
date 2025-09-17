# Config Package

The `config` package provides configuration management for the gomark markdown parser. It defines parser settings, extension controls, and safety options.

## Overview

This package centralizes all configuration options for the markdown parser, allowing fine-grained control over parsing behavior, enabled extensions, and safety features.

## Design Philosophy

gomark's configuration system is designed for **ease of use and flexibility**:

- **Fluent API**: Method chaining for readable configuration
- **Sensible defaults**: DefaultConfig() works well for most use cases
- **Extension control**: Fine-grained enable/disable of features
- **Safety options**: Built-in protection against malicious content
- **Immutable configs**: Configuration methods return new instances for safety

## Core Types

### ParserConfig

The main configuration structure that controls parser behavior:

```go
type ParserConfig struct {
    EnableExtensions ExtensionConfig
    StrictMode       bool
    MaxDepth         int
    AllowHTML        bool
    SafeMode         bool
    MaxFileSize      int64
}
```

#### Fields

- **EnableExtensions**: Controls which markdown extensions are enabled
- **StrictMode**: When true, uses strict CommonMark parsing rules
- **MaxDepth**: Limits maximum nesting depth to prevent stack overflow (default: 100)
- **AllowHTML**: Whether raw HTML is allowed in markdown (default: true)
- **SafeMode**: Whether to sanitize potentially dangerous content (default: false)
- **MaxFileSize**: Maximum file size that can be parsed in bytes (default: 10MB)

### ExtensionConfig

Controls which markdown extensions are enabled:

```go
type ExtensionConfig struct {
    Tables            bool  // Table parsing
    Strikethrough     bool  // ~~strikethrough~~ syntax
    Autolinks         bool  // Automatic URL linking
    TaskLists         bool  // - [ ] and - [x] syntax
    Math              bool  // $inline$ and $$block$$ math
    Highlighting      bool  // ==highlight== syntax
    Subscript         bool  // ~subscript~ syntax
    Superscript       bool  // ^superscript^ syntax
    Spoilers          bool  // ||spoiler|| syntax
    EmbeddedContent   bool  // ![[content]] syntax
    ReferencedContent bool  // [[references]] syntax
    Tags              bool  // #tag syntax
}
```

## Predefined Configurations

### DefaultConfig()

Returns a configuration with all extensions enabled and relaxed parsing:

```go
config := config.DefaultConfig()
// All extensions enabled
// StrictMode: false
// AllowHTML: true
// SafeMode: false
```

### StrictConfig()

Returns a configuration for strict CommonMark compliance:

```go
config := config.StrictConfig()
// Only CommonMark-compatible features enabled
// StrictMode: true
// Most extensions disabled
```

## Configuration Methods

### WithStrictMode(strict bool)

Creates a new config with strict mode enabled/disabled:

```go
config := cfg.WithStrictMode(true)
```

### WithExtension(name string, enabled bool)

Enables or disables a specific extension by name:

```go
config := cfg.WithExtension("tables", false)
config = config.WithExtension("math", true)
```

#### Supported Extension Names

- `"tables"` - Table parsing
- `"strikethrough"` - Strikethrough text
- `"autolinks"` - Automatic URL linking
- `"tasklists"` - Task list items
- `"math"` - Math expressions
- `"highlighting"` - Text highlighting
- `"subscript"` - Subscript text
- `"superscript"` - Superscript text
- `"spoilers"` - Spoiler text
- `"embedded"` - Embedded content
- `"references"` - Referenced content
- `"tags"` - Tag syntax

### WithSafeMode(safe bool)

Enables or disables safe mode:

```go
config := cfg.WithSafeMode(true)
```

### WithMaxDepth(depth int)

Sets the maximum parsing depth:

```go
config := cfg.WithMaxDepth(50)
```

## Usage Examples

### Basic Usage

```go
import "github.com/usememos/gomark/config"

// Use default configuration
cfg := config.DefaultConfig()

// Create strict CommonMark parser
strictCfg := config.StrictConfig()
```

### Custom Configuration

```go
// Start with defaults and customize
cfg := config.DefaultConfig()
cfg = cfg.WithStrictMode(true)
cfg = cfg.WithExtension("tables", false)
cfg = cfg.WithSafeMode(true)
```

### Configuration for Different Use Cases

#### Blog Content (Safe)
```go
cfg := config.DefaultConfig().
    WithSafeMode(true).
    WithExtension("math", false).
    WithMaxDepth(50)
```

#### Documentation (Full Featured)
```go
cfg := config.DefaultConfig().
    WithExtension("tables", true).
    WithExtension("math", true).
    WithStrictMode(false)
```

#### Strict CommonMark Only
```go
cfg := config.StrictConfig().
    WithSafeMode(true)
```

## Security Considerations

- **SafeMode**: When enabled, sanitizes potentially dangerous content
- **AllowHTML**: Controls whether raw HTML is processed
- **MaxDepth**: Prevents stack overflow from deeply nested structures
- **MaxFileSize**: Prevents processing of excessively large files

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

## Configuration Best Practices

### For Different Use Cases

**Blog/Content Sites** (balanced features + safety):
```go
cfg := config.DefaultConfig().
    WithSafeMode(true).
    WithExtension("math", false).       // Disable if not needed
    WithExtension("embedded", false)    // Disable advanced features
```

**Documentation** (full features):
```go
cfg := config.DefaultConfig().
    WithExtension("tables", true).
    WithExtension("math", true).
    WithStrictMode(false)               // Allow extensions
```

**Strict CommonMark** (standards compliance):
```go
cfg := config.StrictConfig().
    WithSafeMode(true)                  // Extra safety
```

**High Performance** (minimal features):
```go
cfg := config.DefaultConfig().
    WithExtension("tables", false).
    WithExtension("math", false).
    WithExtension("highlighting", false).
    WithExtension("spoilers", false)
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