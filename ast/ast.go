package ast

type NodeType string

// Block nodes.
const (
	LineBreakNode         NodeType = "LINE_BREAK"
	ParagraphNode         NodeType = "PARAGRAPH"
	CodeBlockNode         NodeType = "CODE_BLOCK"
	HeadingNode           NodeType = "HEADING"
	HorizontalRuleNode    NodeType = "HORIZONTAL_RULE"
	BlockquoteNode        NodeType = "BLOCKQUOTE"
	ListNode              NodeType = "LIST"
	ListItemNode          NodeType = "LIST_ITEM"
	OrderedListItemNode   NodeType = "ORDERED_LIST_ITEM"
	UnorderedListItemNode NodeType = "UNORDERED_LIST_ITEM"
	TaskListItemNode      NodeType = "TASK_LIST_ITEM"
	MathBlockNode         NodeType = "MATH_BLOCK"
	TableNode             NodeType = "TABLE"
	EmbeddedContentNode   NodeType = "EMBEDDED_CONTENT"
)

// Inline nodes.
const (
	TextNode              NodeType = "TEXT"
	BoldNode              NodeType = "BOLD"
	ItalicNode            NodeType = "ITALIC"
	BoldItalicNode        NodeType = "BOLD_ITALIC"
	CodeNode              NodeType = "CODE"
	ImageNode             NodeType = "IMAGE"
	LinkNode              NodeType = "LINK"
	AutoLinkNode          NodeType = "AUTO_LINK"
	TagNode               NodeType = "TAG"
	StrikethroughNode     NodeType = "STRIKETHROUGH"
	EscapingCharacterNode NodeType = "ESCAPING_CHARACTER"
	MathNode              NodeType = "MATH"
	HighlightNode         NodeType = "HIGHLIGHT"
	SubscriptNode         NodeType = "SUBSCRIPT"
	SuperscriptNode       NodeType = "SUPERSCRIPT"
	ReferencedContentNode NodeType = "REFERENCED_CONTENT"
	SpoilerNode           NodeType = "SPOILER"
	HTMLElementNode       NodeType = "HTML_ELEMENT"
)

type Node interface {
	// Type returns a node type.
	Type() NodeType

	// Restore returns a string representation of this node.
	Restore() string
}

type BaseNode struct {
}

func IsBlockNode(node Node) bool {
	switch node.Type() {
	case ParagraphNode, CodeBlockNode, HeadingNode, HorizontalRuleNode, BlockquoteNode, ListNode, ListItemNode, OrderedListItemNode, UnorderedListItemNode, TaskListItemNode, TableNode, EmbeddedContentNode:
		return true
	default:
		return false
	}
}
