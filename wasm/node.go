package main

import (
	"github.com/usememos/gomark/ast"
	nodepb "github.com/usememos/gomark/proto/gen/node/v1"
)

func ConvertFromASTNodes(rawNodes []ast.Node) []*nodepb.Node {
	nodes := []*nodepb.Node{}
	for _, rawNode := range rawNodes {
		node := ConvertFromASTNode(rawNode)
		nodes = append(nodes, node)
	}
	return nodes
}

func ConvertFromASTNode(rawNode ast.Node) *nodepb.Node {
	node := &nodepb.Node{
		Type: nodepb.NodeType(rawNode.Type()),
	}

	switch n := rawNode.(type) {
	case *ast.LineBreak:
		node.Node = &nodepb.Node_LineBreakNode{}
	case *ast.Paragraph:
		children := ConvertFromASTNodes(n.Children)
		node.Node = &nodepb.Node_ParagraphNode{ParagraphNode: &nodepb.ParagraphNode{Children: children}}
	case *ast.CodeBlock:
		node.Node = &nodepb.Node_CodeBlockNode{CodeBlockNode: &nodepb.CodeBlockNode{Language: n.Language, Content: n.Content}}
	case *ast.Heading:
		children := ConvertFromASTNodes(n.Children)
		node.Node = &nodepb.Node_HeadingNode{HeadingNode: &nodepb.HeadingNode{Level: int32(n.Level), Children: children}}
	case *ast.HorizontalRule:
		node.Node = &nodepb.Node_HorizontalRuleNode{HorizontalRuleNode: &nodepb.HorizontalRuleNode{Symbol: n.Symbol}}
	case *ast.Blockquote:
		children := ConvertFromASTNodes(n.Children)
		node.Node = &nodepb.Node_BlockquoteNode{BlockquoteNode: &nodepb.BlockquoteNode{Children: children}}
	case *ast.OrderedList:
		children := ConvertFromASTNodes(n.Children)
		node.Node = &nodepb.Node_OrderedListNode{OrderedListNode: &nodepb.OrderedListNode{Number: n.Number, Indent: int32(n.Indent), Children: children}}
	case *ast.UnorderedList:
		children := ConvertFromASTNodes(n.Children)
		node.Node = &nodepb.Node_UnorderedListNode{UnorderedListNode: &nodepb.UnorderedListNode{Symbol: n.Symbol, Indent: int32(n.Indent), Children: children}}
	case *ast.TaskList:
		children := ConvertFromASTNodes(n.Children)
		node.Node = &nodepb.Node_TaskListNode{TaskListNode: &nodepb.TaskListNode{Symbol: n.Symbol, Indent: int32(n.Indent), Complete: n.Complete, Children: children}}
	case *ast.MathBlock:
		node.Node = &nodepb.Node_MathBlockNode{MathBlockNode: &nodepb.MathBlockNode{Content: n.Content}}
	case *ast.Table:
		node.Node = &nodepb.Node_TableNode{TableNode: convertTableFromASTNode(n)}
	case *ast.EmbeddedContent:
		node.Node = &nodepb.Node_EmbeddedContentNode{EmbeddedContentNode: &nodepb.EmbeddedContentNode{ResourceName: n.ResourceName, Params: n.Params}}
	case *ast.Text:
		node.Node = &nodepb.Node_TextNode{TextNode: &nodepb.TextNode{Content: n.Content}}
	case *ast.Bold:
		children := ConvertFromASTNodes(n.Children)
		node.Node = &nodepb.Node_BoldNode{BoldNode: &nodepb.BoldNode{Symbol: n.Symbol, Children: children}}
	case *ast.Italic:
		node.Node = &nodepb.Node_ItalicNode{ItalicNode: &nodepb.ItalicNode{Symbol: n.Symbol, Content: n.Content}}
	case *ast.BoldItalic:
		node.Node = &nodepb.Node_BoldItalicNode{BoldItalicNode: &nodepb.BoldItalicNode{Symbol: n.Symbol, Content: n.Content}}
	case *ast.Code:
		node.Node = &nodepb.Node_CodeNode{CodeNode: &nodepb.CodeNode{Content: n.Content}}
	case *ast.Image:
		node.Node = &nodepb.Node_ImageNode{ImageNode: &nodepb.ImageNode{AltText: n.AltText, Url: n.URL}}
	case *ast.Link:
		node.Node = &nodepb.Node_LinkNode{LinkNode: &nodepb.LinkNode{Text: n.Text, Url: n.URL}}
	case *ast.AutoLink:
		node.Node = &nodepb.Node_AutoLinkNode{AutoLinkNode: &nodepb.AutoLinkNode{Url: n.URL, IsRawText: n.IsRawText}}
	case *ast.Tag:
		node.Node = &nodepb.Node_TagNode{TagNode: &nodepb.TagNode{Content: n.Content}}
	case *ast.Strikethrough:
		node.Node = &nodepb.Node_StrikethroughNode{StrikethroughNode: &nodepb.StrikethroughNode{Content: n.Content}}
	case *ast.EscapingCharacter:
		node.Node = &nodepb.Node_EscapingCharacterNode{EscapingCharacterNode: &nodepb.EscapingCharacterNode{Symbol: n.Symbol}}
	case *ast.Math:
		node.Node = &nodepb.Node_MathNode{MathNode: &nodepb.MathNode{Content: n.Content}}
	case *ast.Highlight:
		node.Node = &nodepb.Node_HighlightNode{HighlightNode: &nodepb.HighlightNode{Content: n.Content}}
	case *ast.Subscript:
		node.Node = &nodepb.Node_SubscriptNode{SubscriptNode: &nodepb.SubscriptNode{Content: n.Content}}
	case *ast.Superscript:
		node.Node = &nodepb.Node_SuperscriptNode{SuperscriptNode: &nodepb.SuperscriptNode{Content: n.Content}}
	case *ast.ReferencedContent:
		node.Node = &nodepb.Node_ReferencedContentNode{ReferencedContentNode: &nodepb.ReferencedContentNode{ResourceName: n.ResourceName, Params: n.Params}}
	default:
		node.Node = &nodepb.Node_TextNode{TextNode: &nodepb.TextNode{}}
	}
	return node
}

func convertTableFromASTNode(node *ast.Table) *nodepb.TableNode {
	table := &nodepb.TableNode{
		Header:    node.Header,
		Delimiter: node.Delimiter,
	}
	for _, row := range node.Rows {
		table.Rows = append(table.Rows, &nodepb.TableNode_Row{Cells: row})
	}
	return table
}

func convertToASTNodes(nodes []*nodepb.Node) []ast.Node {
	rawNodes := []ast.Node{}
	for _, node := range nodes {
		rawNode := convertToASTNode(node)
		rawNodes = append(rawNodes, rawNode)
	}
	return rawNodes
}

func convertToASTNode(node *nodepb.Node) ast.Node {
	switch n := node.Node.(type) {
	case *nodepb.Node_LineBreakNode:
		return &ast.LineBreak{}
	case *nodepb.Node_ParagraphNode:
		children := convertToASTNodes(n.ParagraphNode.Children)
		return &ast.Paragraph{Children: children}
	case *nodepb.Node_CodeBlockNode:
		return &ast.CodeBlock{Language: n.CodeBlockNode.Language, Content: n.CodeBlockNode.Content}
	case *nodepb.Node_HeadingNode:
		children := convertToASTNodes(n.HeadingNode.Children)
		return &ast.Heading{Level: int(n.HeadingNode.Level), Children: children}
	case *nodepb.Node_HorizontalRuleNode:
		return &ast.HorizontalRule{Symbol: n.HorizontalRuleNode.Symbol}
	case *nodepb.Node_BlockquoteNode:
		children := convertToASTNodes(n.BlockquoteNode.Children)
		return &ast.Blockquote{Children: children}
	case *nodepb.Node_OrderedListNode:
		children := convertToASTNodes(n.OrderedListNode.Children)
		return &ast.OrderedList{Number: n.OrderedListNode.Number, Indent: int(n.OrderedListNode.Indent), Children: children}
	case *nodepb.Node_UnorderedListNode:
		children := convertToASTNodes(n.UnorderedListNode.Children)
		return &ast.UnorderedList{Symbol: n.UnorderedListNode.Symbol, Indent: int(n.UnorderedListNode.Indent), Children: children}
	case *nodepb.Node_TaskListNode:
		children := convertToASTNodes(n.TaskListNode.Children)
		return &ast.TaskList{Symbol: n.TaskListNode.Symbol, Indent: int(n.TaskListNode.Indent), Complete: n.TaskListNode.Complete, Children: children}
	case *nodepb.Node_MathBlockNode:
		return &ast.MathBlock{Content: n.MathBlockNode.Content}
	case *nodepb.Node_TableNode:
		return convertTableToASTNode(node)
	case *nodepb.Node_EmbeddedContentNode:
		return &ast.EmbeddedContent{ResourceName: n.EmbeddedContentNode.ResourceName, Params: n.EmbeddedContentNode.Params}
	case *nodepb.Node_TextNode:
		return &ast.Text{Content: n.TextNode.Content}
	case *nodepb.Node_BoldNode:
		children := convertToASTNodes(n.BoldNode.Children)
		return &ast.Bold{Symbol: n.BoldNode.Symbol, Children: children}
	case *nodepb.Node_ItalicNode:
		return &ast.Italic{Symbol: n.ItalicNode.Symbol, Content: n.ItalicNode.Content}
	case *nodepb.Node_BoldItalicNode:
		return &ast.BoldItalic{Symbol: n.BoldItalicNode.Symbol, Content: n.BoldItalicNode.Content}
	case *nodepb.Node_CodeNode:
		return &ast.Code{Content: n.CodeNode.Content}
	case *nodepb.Node_ImageNode:
		return &ast.Image{AltText: n.ImageNode.AltText, URL: n.ImageNode.Url}
	case *nodepb.Node_LinkNode:
		return &ast.Link{Text: n.LinkNode.Text, URL: n.LinkNode.Url}
	case *nodepb.Node_AutoLinkNode:
		return &ast.AutoLink{URL: n.AutoLinkNode.Url, IsRawText: n.AutoLinkNode.IsRawText}
	case *nodepb.Node_TagNode:
		return &ast.Tag{Content: n.TagNode.Content}
	case *nodepb.Node_StrikethroughNode:
		return &ast.Strikethrough{Content: n.StrikethroughNode.Content}
	case *nodepb.Node_EscapingCharacterNode:
		return &ast.EscapingCharacter{Symbol: n.EscapingCharacterNode.Symbol}
	case *nodepb.Node_MathNode:
		return &ast.Math{Content: n.MathNode.Content}
	case *nodepb.Node_HighlightNode:
		return &ast.Highlight{Content: n.HighlightNode.Content}
	case *nodepb.Node_SubscriptNode:
		return &ast.Subscript{Content: n.SubscriptNode.Content}
	case *nodepb.Node_SuperscriptNode:
		return &ast.Superscript{Content: n.SuperscriptNode.Content}
	case *nodepb.Node_ReferencedContentNode:
		return &ast.ReferencedContent{ResourceName: n.ReferencedContentNode.ResourceName, Params: n.ReferencedContentNode.Params}
	default:
		return &ast.Text{}
	}
}

func convertTableToASTNode(node *nodepb.Node) *ast.Table {
	table := &ast.Table{
		Header:    node.GetTableNode().Header,
		Delimiter: node.GetTableNode().Delimiter,
	}
	for _, row := range node.GetTableNode().Rows {
		table.Rows = append(table.Rows, row.Cells)
	}
	return table
}
