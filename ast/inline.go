package ast

import "fmt"

type BaseInline struct {
	BaseNode
}

func NewBaseInline(nodeType NodeType) BaseInline {
	return BaseInline{
		BaseNode: BaseNode{
			Type: nodeType,
		},
	}
}

type Text struct {
	BaseInline

	Content string `json:"content"`
}

func (*Text) Type() NodeType {
	return TextNode
}

func (n *Text) Restore() string {
	return n.Content
}

type Bold struct {
	BaseInline

	// Symbol is "*" or "_".
	Symbol   string `json:"symbol"`
	Children []Node `json:"children"`
}

func (*Bold) Type() NodeType {
	return BoldNode
}

func (n *Bold) Restore() string {
	symbol := n.Symbol + n.Symbol
	children := ""
	for _, child := range n.Children {
		children += child.Restore()
	}
	return fmt.Sprintf("%s%s%s", symbol, children, symbol)
}

type Italic struct {
	BaseInline

	// Symbol is "*" or "_".
	Symbol  string `json:"symbol"`
	Content string `json:"content"`
}

func (*Italic) Type() NodeType {
	return ItalicNode
}

func (n *Italic) Restore() string {
	return fmt.Sprintf("%s%s%s", n.Symbol, n.Content, n.Symbol)
}

type BoldItalic struct {
	BaseInline

	// Symbol is "*" or "_".
	Symbol  string `json:"symbol"`
	Content string `json:"content"`
}

func (*BoldItalic) Type() NodeType {
	return BoldItalicNode
}

func (n *BoldItalic) Restore() string {
	symbol := n.Symbol + n.Symbol + n.Symbol
	return fmt.Sprintf("%s%s%s", symbol, n.Content, symbol)
}

type Code struct {
	BaseInline

	Content string `json:"content"`
}

func (*Code) Type() NodeType {
	return CodeNode
}

func (n *Code) Restore() string {
	return fmt.Sprintf("`%s`", n.Content)
}

type Image struct {
	BaseInline

	AltText string `json:"altText"`
	URL     string `json:"url"`
}

func (*Image) Type() NodeType {
	return ImageNode
}

func (n *Image) Restore() string {
	return fmt.Sprintf("![%s](%s)", n.AltText, n.URL)
}

type Link struct {
	BaseInline

	Text string `json:"text"`
	URL  string `json:"url"`
}

func (*Link) Type() NodeType {
	return LinkNode
}

func (n *Link) Restore() string {
	return fmt.Sprintf("[%s](%s)", n.Text, n.URL)
}

type AutoLink struct {
	BaseInline

	URL       string `json:"url"`
	IsRawText bool   `json:"isRawText"`
}

func (*AutoLink) Type() NodeType {
	return AutoLinkNode
}

func (n *AutoLink) Restore() string {
	if n.IsRawText {
		return n.URL
	}
	return fmt.Sprintf("<%s>", n.URL)
}

type Tag struct {
	BaseInline

	Content string `json:"content"`
}

func (*Tag) Type() NodeType {
	return TagNode
}

func (n *Tag) Restore() string {
	return fmt.Sprintf("#%s", n.Content)
}

type Strikethrough struct {
	BaseInline

	Content string `json:"content"`
}

func (*Strikethrough) Type() NodeType {
	return StrikethroughNode
}

func (n *Strikethrough) Restore() string {
	return fmt.Sprintf("~~%s~~", n.Content)
}

type EscapingCharacter struct {
	BaseInline

	Symbol string `json:"symbol"`
}

func (*EscapingCharacter) Type() NodeType {
	return EscapingCharacterNode
}

func (n *EscapingCharacter) Restore() string {
	return fmt.Sprintf("\\%s", n.Symbol)
}

type Math struct {
	BaseInline

	Content string `json:"content"`
}

func (*Math) Type() NodeType {
	return MathNode
}

func (n *Math) Restore() string {
	return fmt.Sprintf("$%s$", n.Content)
}

type Highlight struct {
	BaseInline

	Content string `json:"content"`
}

func (*Highlight) Type() NodeType {
	return HighlightNode
}

func (n *Highlight) Restore() string {
	return fmt.Sprintf("==%s==", n.Content)
}

type Subscript struct {
	BaseInline

	Content string `json:"content"`
}

func (*Subscript) Type() NodeType {
	return SubscriptNode
}

func (n *Subscript) Restore() string {
	return fmt.Sprintf("~%s~", n.Content)
}

type Superscript struct {
	BaseInline

	Content string `json:"content"`
}

func (*Superscript) Type() NodeType {
	return SuperscriptNode
}

func (n *Superscript) Restore() string {
	return fmt.Sprintf("^%s^", n.Content)
}

type ReferencedContent struct {
	BaseInline

	ResourceName string `json:"resourceName"`
	Params       string `json:"params"`
}

func (*ReferencedContent) Type() NodeType {
	return ReferencedContentNode
}

func (n *ReferencedContent) Restore() string {
	params := ""
	if n.Params != "" {
		params = fmt.Sprintf("?%s", n.Params)
	}
	result := fmt.Sprintf("[[%s%s]]", n.ResourceName, params)
	return result
}
