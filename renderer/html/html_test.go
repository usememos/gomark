package html

import (
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
			expected: `<blockquote><p>Hello</p><p>world!</p></blockquote>`,
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
			text:     "#article #memo",
			expected: `<p><span>#article</span> <span>#memo</span></p>`,
		},
		{
			text:     "#article \\#memo",
			expected: `<p><span>#article</span> \#memo</p>`,
		},
		{
			text:     "* Hello\n* world!",
			expected: `<dl><li>Hello</li><br><li>world!</li></dl>`,
		},
		{
			text:     "1. Hello\n2. world\n* !",
			expected: `<dl><li>Hello</li><br><li>world</li><br><li>!</li></dl>`,
		},
		{
			text:     "- [ ] hello\n- [x] world",
			expected: `<dl><li><input type="checkbox" disabled />hello</li><br><li><input type="checkbox" checked disabled />world</li></dl>`,
		},
		{
			text:     "1. ordered\n* unorder\n- [ ] checkbox\n- [x] checked",
			expected: `<dl><li>ordered</li><br><li>unorder</li><br><li><input type="checkbox" disabled />checkbox</li><br><li><input type="checkbox" checked disabled />checked</li></dl>`,
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		nodes, err := parser.Parse(tokens)
		require.NoError(t, err)
		actual := NewHTMLRenderer().Render(nodes)
		require.Equal(t, test.expected, actual)
	}
}
