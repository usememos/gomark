package ast

import "slices"

func IsListItemNode(node Node) bool {
	nodeType := node.Type()
	return slices.Contains([]NodeType{
		OrderedListItemNode, UnorderedListItemNode, TaskListItemNode,
	}, nodeType)
}

func GetListItemKindAndIndent(node Node) (ListKind, int) {
	switch n := node.(type) {
	case *OrderedListItem:
		return OrderedList, n.Indent
	case *UnorderedListItem:
		return UnorderedList, n.Indent
	case *TaskListItem:
		return DescrpitionList, n.Indent
	default:
		return "", 0
	}
}
