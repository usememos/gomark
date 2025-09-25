// Package parser provides the core markdown parsing logic for gomark.
package parser

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/config"
	"github.com/usememos/gomark/parser/internal"
	"github.com/usememos/gomark/parser/tokenizer"
)

type BaseParser interface {
	Match(tokens []*tokenizer.Token) (ast.Node, int)
}

type InlineParser interface {
	BaseParser
}

type BlockParser interface {
	BaseParser
}

func Parse(tokens []*tokenizer.Token) (*ast.Document, error) {
	// Use the original legacy logic for maximum compatibility
	nodes, err := ParseBlock(tokens)
	if err != nil {
		return nil, err
	}
	return &ast.Document{Children: nodes}, nil
}

// ParseWithConfig parses tokens using the provided configuration.
func ParseWithConfig(tokens []*tokenizer.Token, cfg *config.ParserConfig) (*ast.Document, error) {
	registry := NewParserRegistry()
	return ParseWithRegistry(tokens, registry)
}

// ParseWithRegistry parses tokens using the provided parser registry.
func ParseWithRegistry(tokens []*tokenizer.Token, registry *ParserRegistry) (*ast.Document, error) {
	result := NewParseResult()

	nodes, err := parseBlockWithRegistry(tokens, registry, result)
	if err != nil {
		return nil, err
	}

	// Return first error if any occurred during parsing
	if result.HasErrors() {
		return &ast.Document{Children: nodes}, result.FirstError()
	}

	return &ast.Document{Children: nodes}, nil
}

// parseBlockWithRegistry parses tokens using a registry with enhanced error handling.
func parseBlockWithRegistry(tokens []*tokenizer.Token, registry *ParserRegistry, result *ParseResult) ([]ast.Node, error) {
	blockParsers := registry.GetBlockParsers()
	return parseBlockWithParsersAndResult(tokens, blockParsers, result)
}

// parseBlockWithParsersAndResult includes enhanced error reporting.
func parseBlockWithParsersAndResult(tokens []*tokenizer.Token, blockParsers []BlockParser, result *ParseResult) ([]ast.Node, error) {
	// Set lookahead parsers for any paragraph parsers
	for _, parser := range blockParsers {
		if paragraph, ok := parser.(*internal.ParagraphParser); ok {
			// Convert BlockParser slice to BaseParser slice
			baseParsers := make([]internal.BaseParser, len(blockParsers))
			for i, bp := range blockParsers {
				baseParsers[i] = bp
			}
			paragraph.SetLookaheadParsers(baseParsers)
		}
	}

	nodes := []ast.Node{}
	position := 0

	for len(tokens) > 0 {
		lineBreakCount := 0
		for len(tokens) > 0 {
			if tokens[0].Type == tokenizer.NewLine {
				tokens = tokens[1:]
				lineBreakCount++
				position++
				continue
			}
			if tokens[0].Type == tokenizer.Space {
				i := 0
				for i < len(tokens) && tokens[i].Type == tokenizer.Space {
					i++
				}
				if i < len(tokens) && tokens[i].Type == tokenizer.NewLine {
					tokens = tokens[i+1:]
					lineBreakCount++
					position += i + 1
					continue
				}
			}
			break
		}

		if lineBreakCount > 0 {
			for i := 0; i < lineBreakCount; i++ {
				nodes = append(nodes, &ast.LineBreak{})
			}
		}

		if len(tokens) == 0 {
			break
		}

		matched := false
		for _, blockParser := range blockParsers {
			node, size := blockParser.Match(tokens)
			if node != nil && size != 0 && size <= len(tokens) {
				// Consume matched tokens
				tokens = tokens[size:]
				position += size
				nodes = append(nodes, node)
				matched = true
				break
			}
		}

		if !matched {
			// Enhanced error handling with better reporting
			token := tokens[0]
			text := token.Value
			if text == "" {
				text = token.String()
			}

			// Add a warning for unmatched content
			result.AddWarning(NewUnknownElementError("unmatched content", position, token))

			// Create a fallback paragraph with the unmatched token
			nodes = append(nodes, &ast.Paragraph{
				Children: []ast.Node{&ast.Text{Content: text}},
			})
			tokens = tokens[1:]
			position++
		}
	}

	nodes = mergeListItemNodes(nodes)
	return nodes, nil
}

func newDefaultBlockParsers() []BlockParser {
	paragraph := internal.NewParagraphParser()
	parsers := []BlockParser{
		internal.NewCodeBlockParser(),
		internal.NewTableParser(),
		internal.NewHorizontalRuleParser(),
		internal.NewHeadingParser(),
		internal.NewBlockquoteParser(),
		internal.NewOrderedListItemParser(),
		internal.NewTaskListItemParser(),
		internal.NewUnorderedListItemParser(),
		internal.NewMathBlockParser(),
		internal.NewEmbeddedContentParser(),
		paragraph,
	}
	// Convert BlockParser slice to BaseParser slice
	baseParsers := make([]internal.BaseParser, len(parsers))
	for i, bp := range parsers {
		baseParsers[i] = bp
	}
	paragraph.SetLookaheadParsers(baseParsers)
	return parsers
}

// ParseBlock is kept for backward compatibility - use ParseWithConfig for new code.
func ParseBlock(tokens []*tokenizer.Token) ([]ast.Node, error) {
	return ParseBlockWithParsers(tokens, newDefaultBlockParsers())
}

func ParseBlockWithParsers(tokens []*tokenizer.Token, blockParsers []BlockParser) ([]ast.Node, error) {
	for _, parser := range blockParsers {
		if paragraph, ok := parser.(*internal.ParagraphParser); ok {
			// Convert BlockParser slice to BaseParser slice
			baseParsers := make([]internal.BaseParser, len(blockParsers))
			for i, bp := range blockParsers {
				baseParsers[i] = bp
			}
			paragraph.SetLookaheadParsers(baseParsers)
		}
	}

	nodes := []ast.Node{}
	for len(tokens) > 0 {
		lineBreakCount := 0
		for len(tokens) > 0 {
			if tokens[0].Type == tokenizer.NewLine {
				tokens = tokens[1:]
				lineBreakCount++
				continue
			}
			if tokens[0].Type == tokenizer.Space {
				i := 0
				for i < len(tokens) && tokens[i].Type == tokenizer.Space {
					i++
				}
				if i < len(tokens) && tokens[i].Type == tokenizer.NewLine {
					tokens = tokens[i+1:]
					lineBreakCount++
					continue
				}
			}
			break
		}
		if lineBreakCount > 0 {
			for i := 0; i < lineBreakCount; i++ {
				nodes = append(nodes, &ast.LineBreak{})
			}
		}
		if len(tokens) == 0 {
			break
		}

		matched := false
		for _, blockParser := range blockParsers {
			node, size := blockParser.Match(tokens)
			if node != nil && size != 0 && size <= len(tokens) {
				// Consume matched tokens.
				tokens = tokens[size:]
				nodes = append(nodes, node)
				matched = true
				break
			}
		}

		if !matched {
			// Safety valve: prevent infinite loops by consuming a token
			// as plain text when no parser makes progress.
			text := tokens[0].Value
			if text == "" {
				text = tokens[0].String()
			}
			nodes = append(nodes, &ast.Paragraph{
				Children: []ast.Node{&ast.Text{Content: text}},
			})
			tokens = tokens[1:]
		}
	}

	nodes = mergeListItemNodes(nodes)
	return nodes, nil
}

var defaultInlineParsers = []InlineParser{
	internal.NewEscapingCharacterParser(),
	internal.NewLineBreakParser(),
	internal.NewHTMLElementParser(),
	internal.NewBoldItalicParser(),
	internal.NewImageParser(),
	internal.NewLinkParser(),
	internal.NewAutoLinkParser(),
	internal.NewBoldParser(),
	internal.NewItalicParser(),
	internal.NewSpoilerParser(),
	internal.NewHighlightParser(),
	internal.NewCodeParser(),
	internal.NewSubscriptParser(),
	internal.NewSuperscriptParser(),
	internal.NewMathParser(),
	internal.NewReferencedContentParser(),
	internal.NewTagParser(),
	internal.NewStrikethroughParser(),
	internal.NewTextParser(),
}

func ParseInline(tokens []*tokenizer.Token) ([]ast.Node, error) {
	return ParseInlineWithParsers(tokens, defaultInlineParsers)
}

func ParseInlineWithParsers(tokens []*tokenizer.Token, inlineParsers []InlineParser) ([]ast.Node, error) {
	nodes := []ast.Node{}
	for len(tokens) > 0 {
		for _, inlineParser := range inlineParsers {
			node, size := inlineParser.Match(tokens)
			if node != nil && size != 0 {
				// Consume matched tokens.
				tokens = tokens[size:]
				nodes = append(nodes, node)
				break
			}
		}
	}
	return mergeTextNodes(nodes), nil
}

func mergeListItemNodes(nodes []ast.Node) []ast.Node {
	var result []ast.Node
	var stack []*ast.List

	for _, node := range nodes {
		nodeType := node.Type()

		// Handle line breaks.
		if nodeType == ast.LineBreakNode {
			// If the stack is not empty and the last node is a list node, add the line break to the list.
			if len(stack) > 0 && len(result) > 0 && result[len(result)-1].Type() == ast.ListNode {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
			} else {
				result = append(result, node)
			}
			continue
		}

		if ast.IsListItemNode(node) {
			itemKind, itemIndent := ast.GetListItemKindAndIndent(node)

			// Create a new List node if the stack is empty or the current item should be a child of the last item.
			if len(stack) == 0 || (itemKind != stack[len(stack)-1].Kind || itemIndent > stack[len(stack)-1].Indent) {
				newList := &ast.List{
					Kind:     itemKind,
					Indent:   itemIndent,
					Children: []ast.Node{node},
				}

				// Add the new List node to the stack or the result.
				if len(stack) > 0 && itemIndent > stack[len(stack)-1].Indent {
					stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, newList)
				} else {
					result = append(result, newList)
				}
				stack = append(stack, newList)
			} else {
				// Pop the stack until the current item should be a sibling of the last item.
				for len(stack) > 0 && (itemKind != stack[len(stack)-1].Kind || itemIndent < stack[len(stack)-1].Indent) {
					stack = stack[:len(stack)-1]
				}

				// Add the current item to the last List node in the stack or the result.
				if len(stack) > 0 {
					stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
				} else {
					result = append(result, node)
				}
			}
		} else {
			result = append(result, node)
			// Reset the stack if the current node is not a list item node.
			stack = nil
		}
	}

	return result
}

func mergeTextNodes(nodes []ast.Node) []ast.Node {
	if len(nodes) == 0 {
		return nodes
	}
	result := []ast.Node{nodes[0]}
	for i := 1; i < len(nodes); i++ {
		if nodes[i].Type() == ast.TextNode && result[len(result)-1].Type() == ast.TextNode {
			result[len(result)-1].(*ast.Text).Content += nodes[i].(*ast.Text).Content
		} else {
			result = append(result, nodes[i])
		}
	}
	return result
}
