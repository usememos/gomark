package internal

import (
	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

// ParseInline parses inline tokens into AST nodes
func ParseInline(tokens []*tokenizer.Token) ([]ast.Node, error) {
	return ParseInlineWithParsers(tokens, getDefaultInlineParsers())
}

// ParseInlineWithParsers parses inline tokens using the provided parsers
func ParseInlineWithParsers(tokens []*tokenizer.Token, inlineParsers []InlineParser) ([]ast.Node, error) {
	nodes := []ast.Node{}
	for len(tokens) > 0 {
		matched := false
		for _, inlineParser := range inlineParsers {
			node, size := inlineParser.Match(tokens)
			if node != nil && size != 0 {
				tokens = tokens[size:]
				nodes = append(nodes, node)
				matched = true
				break
			}
		}
		if !matched {
			// Fallback: consume one token as text
			nodes = append(nodes, &ast.Text{Content: tokens[0].Value})
			tokens = tokens[1:]
		}
	}
	return mergeTextNodes(nodes), nil
}

// ParseBlock parses block tokens into AST nodes using default parsers
func ParseBlock(tokens []*tokenizer.Token) ([]ast.Node, error) {
	return ParseBlockWithParsers(tokens, getDefaultBlockParsers())
}

// ParseBlockWithParsers parses block tokens using the provided parsers
func ParseBlockWithParsers(tokens []*tokenizer.Token, blockParsers []BlockParser) ([]ast.Node, error) {
	// Set lookahead parsers for any paragraph parsers
	for _, parser := range blockParsers {
		if paragraph, ok := parser.(*ParagraphParser); ok {
			paragraph.SetLookaheadParsers(blockParsers)
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
				tokens = tokens[size:]
				nodes = append(nodes, node)
				matched = true
				break
			}
		}

		if !matched {
			// Safety valve: prevent infinite loops by consuming a token
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

// getDefaultInlineParsers returns the default set of inline parsers
func getDefaultInlineParsers() []InlineParser {
	return []InlineParser{
		NewEscapingCharacterParser(),
		NewLineBreakParser(),
		NewHTMLElementParser(),
		NewBoldItalicParser(),
		NewImageParser(),
		NewLinkParser(),
		NewAutoLinkParser(),
		NewBoldParser(),
		NewItalicParser(),
		NewSpoilerParser(),
		NewHighlightParser(),
		NewCodeParser(),
		NewSubscriptParser(),
		NewSuperscriptParser(),
		NewMathParser(),
		NewReferencedContentParser(),
		NewTagParser(),
		NewStrikethroughParser(),
		NewTextParser(),
	}
}

// getDefaultBlockParsers returns the default set of block parsers
func getDefaultBlockParsers() []BlockParser {
	paragraph := NewParagraphParser()
	parsers := []BlockParser{
		NewCodeBlockParser(),
		NewTableParser(),
		NewHorizontalRuleParser(),
		NewHeadingParser(),
		NewBlockquoteParser(),
		NewOrderedListItemParser(),
		NewTaskListItemParser(),
		NewUnorderedListItemParser(),
		NewMathBlockParser(),
		NewEmbeddedContentParser(),
		paragraph,
	}
	paragraph.SetLookaheadParsers(parsers)
	return parsers
}

// mergeTextNodes merges consecutive text nodes
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

// mergeListItemNodes merges list item nodes into list structures
func mergeListItemNodes(nodes []ast.Node) []ast.Node {
	var result []ast.Node
	var stack []*ast.List

	for _, node := range nodes {
		nodeType := node.Type()

		if nodeType == ast.LineBreakNode {
			if len(stack) > 0 && len(result) > 0 && result[len(result)-1].Type() == ast.ListNode {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
			} else {
				result = append(result, node)
			}
			continue
		}

		if ast.IsListItemNode(node) {
			itemKind, itemIndent := ast.GetListItemKindAndIndent(node)

			if len(stack) == 0 || (itemKind != stack[len(stack)-1].Kind || itemIndent > stack[len(stack)-1].Indent) {
				newList := &ast.List{
					Kind:     itemKind,
					Indent:   itemIndent,
					Children: []ast.Node{node},
				}

				if len(stack) > 0 && itemIndent > stack[len(stack)-1].Indent {
					stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, newList)
				} else {
					result = append(result, newList)
				}
				stack = append(stack, newList)
			} else {
				for len(stack) > 0 && (itemKind != stack[len(stack)-1].Kind || itemIndent < stack[len(stack)-1].Indent) {
					stack = stack[:len(stack)-1]
				}

				if len(stack) > 0 {
					stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
				} else {
					result = append(result, node)
				}
			}
		} else {
			result = append(result, node)
			stack = nil
		}
	}

	return result
}