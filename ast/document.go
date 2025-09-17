package ast

// Document is the root node that stores the entire AST along with optional metadata.
type Document struct {
	BaseNode

	// Children holds the top-level block nodes in source order.
	Children []Node
	// Metadata is an open-ended store for caller-defined properties.
	Metadata map[string]any
}

func (*Document) Type() NodeType {
	return DocumentNode
}

func (d *Document) Restore() string {
	var result string
	for _, child := range d.Children {
		result += child.Restore()
	}
	return result
}

// CloneMetadata returns a shallow copy of the metadata map.
func (d *Document) CloneMetadata() map[string]any {
	if len(d.Metadata) == 0 {
		return nil
	}
	metadataCopy := make(map[string]any, len(d.Metadata))
	for k, v := range d.Metadata {
		metadataCopy[k] = v
	}
	return metadataCopy
}
