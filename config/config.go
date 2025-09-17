package config

// ParserConfig holds configuration options for the parser.
type ParserConfig struct {
	// EnableExtensions controls which markdown extensions are enabled.
	EnableExtensions ExtensionConfig

	// StrictMode controls whether to use strict markdown parsing.
	StrictMode bool

	// MaxDepth limits the maximum nesting depth to prevent stack overflow.
	MaxDepth int

	// AllowHTML controls whether raw HTML is allowed in markdown.
	AllowHTML bool

	// SafeMode controls whether to sanitize potentially dangerous content.
	SafeMode bool

	// MaxFileSize limits the maximum file size that can be parsed.
	MaxFileSize int64
}

// ExtensionConfig controls which markdown extensions are enabled.
type ExtensionConfig struct {
	// Tables enables table parsing.
	Tables bool

	// Strikethrough enables ~~strikethrough~~ syntax.
	Strikethrough bool

	// Autolinks enables automatic URL linking.
	Autolinks bool

	// TaskLists enables task list syntax like - [ ] and - [x].
	TaskLists bool

	// Math enables math syntax like $inline$ and $$block$$.
	Math bool

	// Highlighting enables ==highlight== syntax.
	Highlighting bool

	// Subscript enables ~subscript~ syntax.
	Subscript bool

	// Superscript enables ^superscript^ syntax.
	Superscript bool

	// Spoilers enables ||spoiler|| syntax.
	Spoilers bool

	// EmbeddedContent enables ![[content]] syntax.
	EmbeddedContent bool

	// ReferencedContent enables [[references]] syntax.
	ReferencedContent bool

	// Tags enables #tag syntax.
	Tags bool
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *ParserConfig {
	return &ParserConfig{
		EnableExtensions: ExtensionConfig{
			Tables:            true,
			Strikethrough:     true,
			Autolinks:         true,
			TaskLists:         true,
			Math:              true,
			Highlighting:      true,
			Subscript:         true,
			Superscript:       true,
			Spoilers:          true,
			EmbeddedContent:   true,
			ReferencedContent: true,
			Tags:              true,
		},
		StrictMode:  false,
		MaxDepth:    100,
		AllowHTML:   true,
		SafeMode:    false,
		MaxFileSize: 10 * 1024 * 1024, // 10MB
	}
}

// StrictConfig returns a configuration for strict CommonMark parsing.
func StrictConfig() *ParserConfig {
	return &ParserConfig{
		EnableExtensions: ExtensionConfig{
			Tables:            false,
			Strikethrough:     false,
			Autolinks:         true, // Part of CommonMark
			TaskLists:         false,
			Math:              false,
			Highlighting:      false,
			Subscript:         false,
			Superscript:       false,
			Spoilers:          false,
			EmbeddedContent:   false,
			ReferencedContent: false,
			Tags:              false,
		},
		StrictMode:  true,
		MaxDepth:    50,
		AllowHTML:   false,
		SafeMode:    true,
		MaxFileSize: 1 * 1024 * 1024, // 1MB
	}
}

// GitHubFlavoredConfig returns a configuration similar to GitHub Flavored Markdown.
func GitHubFlavoredConfig() *ParserConfig {
	return &ParserConfig{
		EnableExtensions: ExtensionConfig{
			Tables:            true,
			Strikethrough:     true,
			Autolinks:         true,
			TaskLists:         true,
			Math:              false, // GitHub uses different math syntax
			Highlighting:      false,
			Subscript:         false,
			Superscript:       false,
			Spoilers:          false,
			EmbeddedContent:   false,
			ReferencedContent: false,
			Tags:              false,
		},
		StrictMode:  false,
		MaxDepth:    100,
		AllowHTML:   true,
		SafeMode:    true,
		MaxFileSize: 5 * 1024 * 1024, // 5MB
	}
}

// Clone creates a deep copy of the configuration.
func (c *ParserConfig) Clone() *ParserConfig {
	return &ParserConfig{
		EnableExtensions: ExtensionConfig{
			Tables:            c.EnableExtensions.Tables,
			Strikethrough:     c.EnableExtensions.Strikethrough,
			Autolinks:         c.EnableExtensions.Autolinks,
			TaskLists:         c.EnableExtensions.TaskLists,
			Math:              c.EnableExtensions.Math,
			Highlighting:      c.EnableExtensions.Highlighting,
			Subscript:         c.EnableExtensions.Subscript,
			Superscript:       c.EnableExtensions.Superscript,
			Spoilers:          c.EnableExtensions.Spoilers,
			EmbeddedContent:   c.EnableExtensions.EmbeddedContent,
			ReferencedContent: c.EnableExtensions.ReferencedContent,
			Tags:              c.EnableExtensions.Tags,
		},
		StrictMode:  c.StrictMode,
		MaxDepth:    c.MaxDepth,
		AllowHTML:   c.AllowHTML,
		SafeMode:    c.SafeMode,
		MaxFileSize: c.MaxFileSize,
	}
}

// WithExtension enables or disables a specific extension.
func (c *ParserConfig) WithExtension(name string, enabled bool) *ParserConfig {
	config := c.Clone()
	switch name {
	case "tables":
		config.EnableExtensions.Tables = enabled
	case "strikethrough":
		config.EnableExtensions.Strikethrough = enabled
	case "autolinks":
		config.EnableExtensions.Autolinks = enabled
	case "task_lists":
		config.EnableExtensions.TaskLists = enabled
	case "math":
		config.EnableExtensions.Math = enabled
	case "highlighting":
		config.EnableExtensions.Highlighting = enabled
	case "subscript":
		config.EnableExtensions.Subscript = enabled
	case "superscript":
		config.EnableExtensions.Superscript = enabled
	case "spoilers":
		config.EnableExtensions.Spoilers = enabled
	case "embedded_content":
		config.EnableExtensions.EmbeddedContent = enabled
	case "referenced_content":
		config.EnableExtensions.ReferencedContent = enabled
	case "tags":
		config.EnableExtensions.Tags = enabled
	}
	return config
}

// WithStrictMode sets the strict mode flag.
func (c *ParserConfig) WithStrictMode(strict bool) *ParserConfig {
	config := c.Clone()
	config.StrictMode = strict
	return config
}

// WithMaxDepth sets the maximum nesting depth.
func (c *ParserConfig) WithMaxDepth(depth int) *ParserConfig {
	config := c.Clone()
	config.MaxDepth = depth
	return config
}

// WithAllowHTML sets the HTML allowance flag.
func (c *ParserConfig) WithAllowHTML(allow bool) *ParserConfig {
	config := c.Clone()
	config.AllowHTML = allow
	return config
}

// WithSafeMode sets the safe mode flag.
func (c *ParserConfig) WithSafeMode(safe bool) *ParserConfig {
	config := c.Clone()
	config.SafeMode = safe
	return config
}

// WithMaxFileSize sets the maximum file size.
func (c *ParserConfig) WithMaxFileSize(size int64) *ParserConfig {
	config := c.Clone()
	config.MaxFileSize = size
	return config
}