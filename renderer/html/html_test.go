package html

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/gomark/parser"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestHTMLRenderer(t *testing.T) {
	tests := []struct {
		text     string
		expected string
	}{
		{
			text:     "Hello world!",
			expected: `<p>Hello world!</p>`,
		},
		{
			text:     "# Hello world!",
			expected: `<h1>Hello world!</h1>`,
		},
		{
			text:     "> Hello\n> world!",
			expected: `<blockquote><p>Hello<br>world!</p></blockquote>`,
		},
		{
			text:     "*Hello* world!",
			expected: `<p><em>Hello</em> world!</p>`,
		},
		{
			text:     "Hello world!\n\nNew paragraph.",
			expected: "<p>Hello world!</p><br><p>New paragraph.</p>",
		},
		{
			text:     "**Hello** world!",
			expected: `<p><strong>Hello</strong> world!</p>`,
		},
		{
			text:     "***Hello***",
			expected: `<p><strong><em>Hello</em></strong></p>`,
		},
		{
			text:     "***Hello _world_***",
			expected: `<p><strong><em>Hello <em>world</em></em></strong></p>`,
		},
		{
			text:     "#article #memo",
			expected: `<p><span>#article</span> <span>#memo</span></p>`,
		},
		{
			text:     "#article \\#memo",
			expected: `<p><span>#article</span> \#memo</p>`,
		},
		{
			text:     "* Hello\n* world!",
			expected: `<ul><li>Hello</li><br><li>world!</li></ul>`,
		},
		{
			text:     "- [ ] hello\n- [x] world",
			expected: `<dl><li><input type="checkbox" disabled />hello</li><br><li><input type="checkbox" checked disabled />world</li></dl>`,
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		doc, err := parser.Parse(tokens)
		require.NoError(t, err)
		actual := NewHTMLRenderer().RenderDocument(doc)
		require.Equal(t, test.expected, actual, fmt.Sprintf("Test case: %s", test.text))
	}
}
