// Package config provides configuration structures for the gomark parser.
package config

// ParserConfig holds configuration options for the parser.
// For simplicity, most features are enabled by default.
type ParserConfig struct {
	// MaxDepth limits the maximum nesting depth to prevent stack overflow.
	MaxDepth int

	// MaxFileSize limits the maximum file size that can be parsed.
	MaxFileSize int64
}

// DefaultConfig returns a configuration with sensible defaults for most use cases.
// All markdown extensions are enabled by default.
func DefaultConfig() *ParserConfig {
	return &ParserConfig{
		MaxDepth:    200,              // Allow deeper nesting for complex content
		MaxFileSize: 50 * 1024 * 1024, // 50MB for large documents with embedded content
	}
}

// Clone creates a deep copy of the configuration.
func (c *ParserConfig) Clone() *ParserConfig {
	return &ParserConfig{
		MaxDepth:    c.MaxDepth,
		MaxFileSize: c.MaxFileSize,
	}
}

// WithMaxDepth sets the maximum nesting depth.
func (c *ParserConfig) WithMaxDepth(depth int) *ParserConfig {
	config := c.Clone()
	config.MaxDepth = depth
	return config
}

// WithMaxFileSize sets the maximum file size.
func (c *ParserConfig) WithMaxFileSize(size int64) *ParserConfig {
	config := c.Clone()
	config.MaxFileSize = size
	return config
}
