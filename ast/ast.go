package ast

type NodeType string

// Block nodes.
const (
	LineBreakNode       NodeType = "LineBreak"
	ParagraphNode       NodeType = "Paragraph"
	CodeBlockNode       NodeType = "CodeBlock"
	HeadingNode         NodeType = "Heading"
	HorizontalRuleNode  NodeType = "HorizontalRule"
	BlockquoteNode      NodeType = "Blockquote"
	OrderedListNode     NodeType = "OrderedList"
	UnorderedListNode   NodeType = "UnorderedList"
	TaskListNode        NodeType = "TaskList"
	MathBlockNode       NodeType = "MathBlock"
	TableNode           NodeType = "Table"
	EmbeddedContentNode NodeType = "EmbeddedContent"
)

// Inline nodes.
const (
	TextNode              NodeType = "Text"
	BoldNode              NodeType = "Bold"
	ItalicNode            NodeType = "Italic"
	BoldItalicNode        NodeType = "BoldItalic"
	CodeNode              NodeType = "Code"
	ImageNode             NodeType = "Image"
	LinkNode              NodeType = "Link"
	AutoLinkNode          NodeType = "AutoLink"
	TagNode               NodeType = "Tag"
	StrikethroughNode     NodeType = "Strikethrough"
	EscapingCharacterNode NodeType = "EscapingCharacter"
	MathNode              NodeType = "Math"
	HighlightNode         NodeType = "Highlight"
	SubscriptNode         NodeType = "Subscript"
	SuperscriptNode       NodeType = "Superscript"
	ReferencedContentNode NodeType = "ReferencedContent"
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
	Type NodeType `json:"type"`

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
