package parser

import (
	"github.com/usememos/gomark/ast"
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

func Parse(tokens []*tokenizer.Token) ([]ast.Node, error) {
	return ParseBlock(tokens)
}

var defaultBlockParsers = []BlockParser{
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
	NewParagraphParser(),
	NewLineBreakParser(),
}

func ParseBlock(tokens []*tokenizer.Token) ([]ast.Node, error) {
	return ParseBlockWithParsers(tokens, defaultBlockParsers)
}

func ParseBlockWithParsers(tokens []*tokenizer.Token, blockParsers []BlockParser) ([]ast.Node, error) {
	nodes := []ast.Node{}
	for len(tokens) > 0 {
		for _, blockParser := range blockParsers {
			node, size := blockParser.Match(tokens)
			if node != nil && size != 0 {
				// Consume matched tokens.
				tokens = tokens[size:]
				nodes = append(nodes, node)
				break
			}
		}
	}

	nodes = mergeListItemNodes(nodes)
	return nodes, nil
}

var defaultInlineParsers = []InlineParser{
	NewEscapingCharacterParser(),
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
	NewLineBreakParser(),
	NewTextParser(),
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
