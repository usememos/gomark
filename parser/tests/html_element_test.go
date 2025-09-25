package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/internal"
	"github.com/usememos/gomark/parser/tokenizer"
)

func TestHTMLElementParser(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected ast.Node
	}{
		// Self-closing elements
		{
			name: "br element",
			text: "<br>",
			expected: &ast.HTMLElement{
				TagName:       "br",
				Attributes:    map[string]string{},
				IsSelfClosing: true,
			},
		},
		{
			name: "br element self-closing",
			text: "<br />",
			expected: &ast.HTMLElement{
				TagName:       "br",
				Attributes:    map[string]string{},
				IsSelfClosing: true,
			},
		},
		{
			name: "img element basic",
			text: `<img src="test.jpg" alt="Test">`,
			expected: &ast.HTMLElement{
				TagName: "img",
				Attributes: map[string]string{
					"src": "test.jpg",
					"alt": "Test",
				},
				IsSelfClosing: true,
			},
		},
		{
			name: "img element with attributes",
			text: `<img src="image.png" alt="Description" width="100" height="50">`,
			expected: &ast.HTMLElement{
				TagName: "img",
				Attributes: map[string]string{
					"src":    "image.png",
					"alt":    "Description",
					"width":  "100",
					"height": "50",
				},
				IsSelfClosing: true,
			},
		},

		// Simple text elements
		{
			name: "kbd element",
			text: "<kbd>Ctrl</kbd>",
			expected: &ast.HTMLElement{
				TagName:    "kbd",
				Attributes: map[string]string{},
				Children:   []ast.Node{&ast.Text{Content: "Ctrl"}},
			},
		},
		{
			name: "kbd with class attribute",
			text: `<kbd class="key">Enter</kbd>`,
			expected: &ast.HTMLElement{
				TagName:    "kbd",
				Attributes: map[string]string{"class": "key"},
				Children:   []ast.Node{&ast.Text{Content: "Enter"}},
			},
		},

		// Container elements
		{
			name: "small element",
			text: "<small>Copyright 2024</small>",
			expected: &ast.HTMLElement{
				TagName:    "small",
				Attributes: map[string]string{},
				Children:   []ast.Node{&ast.Text{Content: "Copyright 2024"}},
			},
		},
		{
			name: "mark element",
			text: "<mark>highlighted text</mark>",
			expected: &ast.HTMLElement{
				TagName:    "mark",
				Attributes: map[string]string{},
				Children:   []ast.Node{&ast.Text{Content: "highlighted text"}},
			},
		},
		{
			name: "mark with class",
			text: `<mark class="warning">important</mark>`,
			expected: &ast.HTMLElement{
				TagName:    "mark",
				Attributes: map[string]string{"class": "warning"},
				Children:   []ast.Node{&ast.Text{Content: "important"}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := tokenizer.Tokenize(test.text)
			node, _ := internal.NewHTMLElementParser().Match(tokens)
			require.Equal(t, test.expected, node)
		})
	}
}

func TestHTMLElementParserFailures(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "unsupported element",
			text: "<div>content</div>",
		},
		{
			name: "img without src",
			text: `<img alt="Test">`,
		},
		{
			name: "img without alt",
			text: `<img src="test.jpg">`,
		},
		{
			name: "unclosed element",
			text: "<kbd>text",
		},
		{
			name: "malformed tag",
			text: "<kbd>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens := tokenizer.Tokenize(test.text)
			node, _ := internal.NewHTMLElementParser().Match(tokens)
			require.Nil(t, node, "Expected parser to fail for: %s", test.text)
		})
	}
}
