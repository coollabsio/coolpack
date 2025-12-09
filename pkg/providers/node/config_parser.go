package node

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

// ConfigParser parses JavaScript/TypeScript config files using tree-sitter
type ConfigParser struct {
	tsParser *sitter.Parser
	jsParser *sitter.Parser
}

// NewConfigParser creates a new config parser
func NewConfigParser() *ConfigParser {
	tsParser := sitter.NewParser()
	tsParser.SetLanguage(typescript.GetLanguage())

	jsParser := sitter.NewParser()
	jsParser.SetLanguage(javascript.GetLanguage())

	return &ConfigParser{
		tsParser: tsParser,
		jsParser: jsParser,
	}
}

// ParseTS parses TypeScript source code and returns the root node
func (p *ConfigParser) ParseTS(source []byte) (*sitter.Node, error) {
	tree, err := p.tsParser.ParseCtx(context.Background(), nil, source)
	if err != nil {
		return nil, err
	}
	return tree.RootNode(), nil
}

// ParseJS parses JavaScript source code and returns the root node
func (p *ConfigParser) ParseJS(source []byte) (*sitter.Node, error) {
	tree, err := p.jsParser.ParseCtx(context.Background(), nil, source)
	if err != nil {
		return nil, err
	}
	return tree.RootNode(), nil
}

// FindPropertyValue searches for a property with the given name in an object
// and returns its string value if found
func FindPropertyValue(node *sitter.Node, source []byte, propertyName string) string {
	if node == nil {
		return ""
	}

	// Recursively search the tree
	return findPropertyInNode(node, source, propertyName)
}

// FindNestedPropertyValue searches for a nested property path (e.g., "server.preset")
// and returns its string value if found
func FindNestedPropertyValue(node *sitter.Node, source []byte, path ...string) string {
	if node == nil || len(path) == 0 {
		return ""
	}

	// Find the first property in the path
	objectNode := findPropertyObjectNode(node, source, path[0])
	if objectNode == nil {
		return ""
	}

	// If there's only one property in the path, return its value
	if len(path) == 1 {
		return trimQuotes(getNodeText(objectNode, source))
	}

	// Otherwise, continue searching in the nested object
	return FindNestedPropertyValue(objectNode, source, path[1:]...)
}

// findPropertyObjectNode finds a property and returns its value node (for nested lookups)
func findPropertyObjectNode(node *sitter.Node, source []byte, propertyName string) *sitter.Node {
	if node == nil {
		return nil
	}

	nodeType := node.Type()

	if nodeType == "pair" || nodeType == "property_assignment" {
		keyNode := node.ChildByFieldName("key")
		valueNode := node.ChildByFieldName("value")

		if keyNode != nil && valueNode != nil {
			key := getNodeText(keyNode, source)
			key = trimQuotes(key)

			if key == propertyName {
				return valueNode
			}
		}
	}

	// Recurse into children
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if result := findPropertyObjectNode(child, source, propertyName); result != nil {
			return result
		}
	}

	return nil
}

func findPropertyInNode(node *sitter.Node, source []byte, propertyName string) string {
	if node == nil {
		return ""
	}

	// Check if this is a property assignment (pair in tree-sitter terms)
	nodeType := node.Type()

	if nodeType == "pair" || nodeType == "property_assignment" {
		// Get the key and value
		keyNode := node.ChildByFieldName("key")
		valueNode := node.ChildByFieldName("value")

		if keyNode != nil && valueNode != nil {
			key := getNodeText(keyNode, source)
			// Remove quotes if present
			key = trimQuotes(key)

			if key == propertyName {
				value := getNodeText(valueNode, source)
				return trimQuotes(value)
			}
		}
	}

	// Also handle shorthand property patterns like { output: 'export' }
	if nodeType == "property_identifier" || nodeType == "identifier" {
		text := getNodeText(node, source)
		if text == propertyName {
			// Check if sibling is the value
			sibling := node.NextSibling()
			if sibling != nil {
				return trimQuotes(getNodeText(sibling, source))
			}
		}
	}

	// Recurse into children
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if result := findPropertyInNode(child, source, propertyName); result != "" {
			return result
		}
	}

	return ""
}

func getNodeText(node *sitter.Node, source []byte) string {
	if node == nil {
		return ""
	}
	start := node.StartByte()
	end := node.EndByte()
	if int(end) > len(source) {
		end = uint32(len(source))
	}
	return string(source[start:end])
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
