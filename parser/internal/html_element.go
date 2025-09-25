package internal

import (
	"strings"

	"github.com/usememos/gomark/ast"
	"github.com/usememos/gomark/parser/tokenizer"
)

type HTMLElementParser struct{}

func NewHTMLElementParser() *HTMLElementParser {
	return &HTMLElementParser{}
}

// ElementType defines the parsing behavior for different HTML elements.
type ElementType int

const (
	SelfClosingElement ElementType = iota // br, img
	SimpleTextElement                     // kbd (text content only)
	ContainerElement                      // small, mark (nested markdown content)
)

// Phase 1 supported HTML elements with their parsing behavior.
var supportedElements = map[string]ElementType{
	"br":    SelfClosingElement,
	"img":   SelfClosingElement,
	"kbd":   SimpleTextElement,
	"small": ContainerElement,
	"mark":  ContainerElement,
}

func (*HTMLElementParser) Match(tokens []*tokenizer.Token) (ast.Node, int) {
	if len(tokens) < 3 || tokens[0].Type != tokenizer.LessThan {
		return nil, 0
	}

	// Extract tag name (tokens[1] should contain the tag name)
	tagName := strings.ToLower(tokenizer.Stringify([]*tokenizer.Token{tokens[1]}))
	elementType, supported := supportedElements[tagName]
	if !supported {
		return nil, 0
	}

	switch elementType {
	case SelfClosingElement:
		return parseSelfClosingElement(tokens, tagName)
	case SimpleTextElement:
		return parseSimpleTextElement(tokens, tagName)
	case ContainerElement:
		return parseContainerElement(tokens, tagName)
	default:
		return nil, 0
	}
}

// Parse self-closing elements: <br>, <img src="..." alt="...">.
func parseSelfClosingElement(tokens []*tokenizer.Token, tagName string) (ast.Node, int) {
	// Find the closing >
	greaterThanIndex := tokenizer.FindUnescaped(tokens, tokenizer.GreaterThan)
	if greaterThanIndex == -1 {
		return nil, 0
	}

	// Parse attributes if present
	attributeTokens := tokens[2:greaterThanIndex]
	attributes := parseAttributes(attributeTokens)

	// Validate required attributes for img
	if tagName == "img" {
		if _, hasSrc := attributes["src"]; !hasSrc {
			return nil, 0 // img requires src attribute
		}
		if _, hasAlt := attributes["alt"]; !hasAlt {
			return nil, 0 // img requires alt attribute
		}
	}

	return &ast.HTMLElement{
		TagName:       tagName,
		Attributes:    attributes,
		IsSelfClosing: true,
	}, greaterThanIndex + 1
}

// Parse simple text elements: <kbd>text content</kbd>.
func parseSimpleTextElement(tokens []*tokenizer.Token, tagName string) (ast.Node, int) {
	// Find opening tag close
	openCloseIndex := tokenizer.FindUnescaped(tokens, tokenizer.GreaterThan)
	if openCloseIndex == -1 {
		return nil, 0
	}

	// Find closing tag
	closingTagStart := findClosingTag(tokens[openCloseIndex+1:], tagName)
	if closingTagStart == -1 {
		return nil, 0
	}

	// Parse opening tag attributes
	attributeTokens := tokens[2:openCloseIndex]
	attributes := parseAttributes(attributeTokens)

	// Extract content between opening and closing tags
	contentTokens := tokens[openCloseIndex+1 : openCloseIndex+1+closingTagStart]
	content := strings.TrimSpace(tokenizer.Stringify(contentTokens))

	// Find the end of closing tag
	closingTagTokens := tokens[openCloseIndex+1+closingTagStart:]
	closingTagEnd := findClosingTagEnd(closingTagTokens, tagName)
	if closingTagEnd == -1 {
		return nil, 0
	}

	return &ast.HTMLElement{
		TagName:    tagName,
		Attributes: attributes,
		Children:   []ast.Node{&ast.Text{Content: content}},
	}, openCloseIndex + 1 + closingTagStart + closingTagEnd + 1
}

// Parse container elements: <small>nested **markdown** content</small>.
func parseContainerElement(tokens []*tokenizer.Token, tagName string) (ast.Node, int) {
	// Find opening tag close
	openCloseIndex := tokenizer.FindUnescaped(tokens, tokenizer.GreaterThan)
	if openCloseIndex == -1 {
		return nil, 0
	}

	// Find closing tag
	closingTagStart := findClosingTag(tokens[openCloseIndex+1:], tagName)
	if closingTagStart == -1 {
		return nil, 0
	}

	// Parse opening tag attributes
	attributeTokens := tokens[2:openCloseIndex]
	attributes := parseAttributes(attributeTokens)

	// Extract content between opening and closing tags
	contentTokens := tokens[openCloseIndex+1 : openCloseIndex+1+closingTagStart]

	// Recursively parse inner content for nested markdown
	children, err := ParseInline(contentTokens)
	if err != nil || len(children) == 0 {
		// Fallback to plain text if parsing fails
		content := strings.TrimSpace(tokenizer.Stringify(contentTokens))
		children = []ast.Node{&ast.Text{Content: content}}
	}

	// Find the end of closing tag
	closingTagTokens := tokens[openCloseIndex+1+closingTagStart:]
	closingTagEnd := findClosingTagEnd(closingTagTokens, tagName)
	if closingTagEnd == -1 {
		return nil, 0
	}

	return &ast.HTMLElement{
		TagName:    tagName,
		Attributes: attributes,
		Children:   children,
	}, openCloseIndex + 1 + closingTagStart + closingTagEnd + 1
}

// parseAttributes extracts key-value pairs from attribute tokens.
func parseAttributes(tokens []*tokenizer.Token) map[string]string {
	attributes := make(map[string]string)

	i := 0
	for i < len(tokens) {
		// Skip whitespace
		for i < len(tokens) && tokens[i].Type == tokenizer.Space {
			i++
		}
		if i >= len(tokens) {
			break
		}

		// Get attribute name
		if tokens[i].Type != tokenizer.Text {
			i++
			continue
		}
		attrName := strings.ToLower(tokens[i].Value)
		i++

		// Skip whitespace and find =
		for i < len(tokens) && tokens[i].Type == tokenizer.Space {
			i++
		}
		if i >= len(tokens) || tokens[i].Type != tokenizer.EqualSign {
			// Attribute without value, skip
			continue
		}
		i++ // skip =

		// Skip whitespace
		for i < len(tokens) && tokens[i].Type == tokenizer.Space {
			i++
		}
		if i >= len(tokens) {
			break
		}

		// Parse attribute value
		var attrValue string
		if tokens[i].Type == tokenizer.Apostrophe {
			// Handle single quoted values 'value'
			i++ // skip opening quote
			valueTokens := []*tokenizer.Token{}
			for i < len(tokens) && tokens[i].Type != tokenizer.Apostrophe {
				valueTokens = append(valueTokens, tokens[i])
				i++
			}
			if i < len(tokens) {
				i++ // skip closing quote
			}
			attrValue = tokenizer.Stringify(valueTokens)
		} else if tokens[i].Type == tokenizer.Text && strings.HasPrefix(tokens[i].Value, "\"") {
			// Handle double quoted values "value" - collect tokens until we find end quote
			valueTokens := []*tokenizer.Token{}

			// First token starts with quote
			firstValue := tokens[i].Value
			if strings.HasSuffix(firstValue, "\"") && len(firstValue) > 1 {
				// Single token contains the whole quoted value: "value"
				attrValue = strings.Trim(firstValue, "\"")
				i++
			} else {
				// Multi-token quoted value: "some text"
				valueTokens = append(valueTokens, &tokenizer.Token{Type: tokenizer.Text, Value: strings.TrimPrefix(firstValue, "\"")})
				i++

				// Collect tokens until we find the closing quote
				for i < len(tokens) {
					if tokens[i].Type == tokenizer.Text && strings.HasSuffix(tokens[i].Value, "\"") {
						// Found closing quote
						endValue := strings.TrimSuffix(tokens[i].Value, "\"")
						if endValue != "" {
							valueTokens = append(valueTokens, &tokenizer.Token{Type: tokens[i].Type, Value: endValue})
						}
						i++
						break
					}
					valueTokens = append(valueTokens, tokens[i])
					i++
				}
				attrValue = tokenizer.Stringify(valueTokens)
			}
		} else {
			// Unquoted value
			attrValue = tokens[i].Value
			i++
		}

		attributes[attrName] = attrValue
	}

	return attributes
}

// findClosingTag finds the start position of closing tag </tagName>.
func findClosingTag(tokens []*tokenizer.Token, tagName string) int {
	for i := 0; i < len(tokens)-3; i++ {
		if tokens[i].Type == tokenizer.LessThan &&
			tokens[i+1].Type == tokenizer.Slash &&
			strings.ToLower(tokenizer.Stringify([]*tokenizer.Token{tokens[i+2]})) == tagName {
			return i
		}
	}
	return -1
}

// findClosingTagEnd finds the end position of closing tag (after >).
func findClosingTagEnd(tokens []*tokenizer.Token, tagName string) int {
	if len(tokens) < 4 {
		return -1
	}
	// Expect: <, /, tagName, >
	if tokens[0].Type == tokenizer.LessThan &&
		tokens[1].Type == tokenizer.Slash &&
		strings.ToLower(tokenizer.Stringify([]*tokenizer.Token{tokens[2]})) == tagName &&
		tokens[3].Type == tokenizer.GreaterThan {
		return 3
	}
	return -1
}
