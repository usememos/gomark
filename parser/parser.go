package parser

import (
	"slices"

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
		var prevNode, nextNode, prevResultNode ast.Node
		if i > 0 {
			prevNode = nodes[i-1]
		}
		if i < len(nodes)-1 {
			nextNode = nodes[i+1]
		}
		if len(result) > 0 {
			prevResultNode = result[len(result)-1]
		}
		switch nodes[i].(type) {
		case *ast.OrderedListItem, *ast.UnorderedListItem, *ast.TaskListItem:
			var listKind ast.ListKind
			var indent int
			switch item := nodes[i].(type) {
			case *ast.OrderedListItem:
				listKind = ast.OrderedList
				indent = item.Indent
			case *ast.UnorderedListItem:
				listKind = ast.UnorderedList
				indent = item.Indent
			case *ast.TaskListItem:
				listKind = ast.DescrpitionList
				indent = item.Indent
			}

			indent /= 2
			if prevResultNode == nil || prevResultNode.Type() != ast.ListNode || prevResultNode.(*ast.List).Kind != listKind || prevResultNode.(*ast.List).Indent > indent {
				prevResultNode = &ast.List{
					BaseBlock: ast.BaseBlock{},
					Kind:      listKind,
					Indent:    indent,
					Children:  []ast.Node{nodes[i]},
				}
				result = append(result, prevResultNode)
				continue
			}

			listNode, ok := prevResultNode.(*ast.List)
			if !ok {
				continue
			}
			if listNode.Indent != indent {
				parent := findListPossibleParent(listNode, indent)
				if parent == nil {
					parent = &ast.List{
						BaseBlock: ast.BaseBlock{},
						Kind:      listKind,
						Indent:    indent,
					}
					listNode.Children = append(listNode.Children, parent)
				}
				parent.Children = append(parent.Children, nodes[i])
			} else {
				listNode.Children = append(listNode.Children, nodes[i])
			}
		case *ast.LineBreak:
			if prevResultNode != nil && prevResultNode.Type() == ast.ListNode &&
				// Check if the prev node is not a line break node.
				(prevNode == nil || prevNode.Type() != ast.LineBreakNode) &&
				// Check if the next node is a list item node.
				(nextNode == nil || slices.Contains([]ast.NodeType{ast.OrderedListItemNode, ast.UnorderedListItemNode, ast.TaskListItemNode}, nextNode.Type())) {
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

func findListPossibleParent(listNode *ast.List, indent int) *ast.List {
	if listNode.Indent == indent {
		return listNode
	}
	if listNode.Indent < indent {
		return nil
	}
	if len(listNode.Children) == 0 {
		return nil
	}
	lastChild := listNode.Children[len(listNode.Children)-1]
	if lastChild.Type() != ast.ListNode {
		return nil
	}
	return findListPossibleParent(lastChild.(*ast.List), indent)
}
