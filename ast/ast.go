package ast

type NodeType string

// Block nodes.
const (
	LineBreakNode       NodeType = "LINE_BREAK"
	ParagraphNode       NodeType = "PARAGRAPH"
	CodeBlockNode       NodeType = "CODE_BLOCK"
	HeadingNode         NodeType = "HEADING"
	HorizontalRuleNode  NodeType = "HORIZONTAL_RULE"
	BlockquoteNode      NodeType = "BLOCKQUOTE"
	OrderedListNode     NodeType = "ORDERED_LIST"
	UnorderedListNode   NodeType = "UNORDERED_LIST"
	TaskListNode        NodeType = "TASK_LIST"
	MathBlockNode       NodeType = "MATH_BLOCK"
	TableNode           NodeType = "TABLE"
	EmbeddedContentNode NodeType = "EMBEDDED_CONTENT"
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
)

type Node interface {
	// Type returns a node type.
	Type() NodeType

	// Restore returns a string representation of this node.
	Restore() string

	// PrevSibling returns a previous sibling node of this node.
	PrevSibling() Node

	// NextSibling returns a next sibling node of this node.
	NextSibling() Node

	// SetPrevSibling sets a previous sibling node to this node.
	SetPrevSibling(Node)

	// SetNextSibling sets a next sibling node to this node.
	SetNextSibling(Node)
}

type BaseNode struct {
	prevSibling Node

	nextSibling Node
}

func (n *BaseNode) PrevSibling() Node {
	return n.prevSibling
}

func (n *BaseNode) NextSibling() Node {
	return n.nextSibling
}

func (n *BaseNode) SetPrevSibling(node Node) {
	n.prevSibling = node
}

func (n *BaseNode) SetNextSibling(node Node) {
	n.nextSibling = node
}

func IsBlockNode(node Node) bool {
	switch node.Type() {
	case ParagraphNode, CodeBlockNode, HeadingNode, HorizontalRuleNode, BlockquoteNode, OrderedListNode, UnorderedListNode, TaskListNode, MathBlockNode, TableNode, EmbeddedContentNode:
		return true
	default:
		return false
	}
}
