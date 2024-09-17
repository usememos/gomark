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
	return mergeListItemNodes(nodes), nil
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
	if len(nodes) == 0 {
		return nodes
	}
	result := []ast.Node{}
	for i := 0; i < len(nodes); i++ {
		var prevNode, prevResultNode ast.Node
		if i > 0 {
			prevNode = nodes[i-1]
		}
		if len(result) > 0 {
			prevResultNode = result[len(result)-1]
		}
		switch nodes[i].(type) {
		case *ast.OrderedListItem, *ast.UnorderedListItem, *ast.TaskListItem:
			if prevResultNode == nil || prevResultNode.Type() != ast.ListNode {
				prevResultNode = &ast.List{
					BaseBlock: ast.BaseBlock{},
				}
				result = append(result, prevResultNode)
			}
			prevResultNode.(*ast.List).Children = append(prevResultNode.(*ast.List).Children, nodes[i])
		case *ast.LineBreak:
			if prevResultNode != nil && prevResultNode.Type() == ast.ListNode && (prevNode == nil || prevNode.Type() != ast.LineBreakNode) {
				prevResultNode.(*ast.List).Children = append(prevResultNode.(*ast.List).Children, nodes[i])
			} else {
				result = append(result, nodes[i])
			}
		default:
			result = append(result, nodes[i])
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
