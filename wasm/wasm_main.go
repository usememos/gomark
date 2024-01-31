package main

import (
	"encoding/json"
	"syscall/js"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/yourselfhosted/gomark/parser"
	"github.com/yourselfhosted/gomark/parser/tokenizer"
	nodepb "github.com/yourselfhosted/gomark/proto/gen/node/v1"
	"github.com/yourselfhosted/gomark/restore"
)

// Parse converts markdown to nodes.
func Parse(this js.Value, inputs []js.Value) any {
	markdown := inputs[0].String()
	tokens := tokenizer.Tokenize(markdown)
	astNodes, err := parser.Parse(tokens)
	if err != nil {
		panic(err)
	}

	nodes := ConvertFromASTNodes(astNodes)
	data := []interface{}{}
	for _, node := range nodes {
		bytes, _ := protojson.Marshal(node)
		v := map[string]interface{}{}
		json.Unmarshal(bytes, &v)
		data = append(data, v)
	}
	return data
}

// Restore converts nodes to markdown.
func Restore(this js.Value, inputs []js.Value) any {
	astNodes, ok := convertJSValueToInterface(inputs[0]).([]interface{})
	if !ok {
		return nil
	}

	nodes := []*nodepb.Node{}
	for _, n := range astNodes {
		bytes, _ := json.Marshal(n)
		v1Node := &nodepb.Node{}
		protojson.Unmarshal(bytes, v1Node)
		nodes = append(nodes, v1Node)
	}
	content := restore.Restore(convertToASTNodes(nodes))
	return content
}

// convertJSValueToInterface converts a js.Value to a Go interface{}.
func convertJSValueToInterface(value js.Value) interface{} {
	switch value.Type() {
	case js.TypeString:
		return value.String()
	case js.TypeNumber:
		return value.Float()
	case js.TypeBoolean:
		return value.Bool()
	case js.TypeObject:
		if value.InstanceOf(js.Global().Get("Array")) {
			length := value.Length()
			array := make([]interface{}, length)
			for i := 0; i < length; i++ {
				array[i] = convertJSValueToInterface(value.Index(i))
			}
			return array
		} else {
			obj := make(map[string]interface{})
			keys := js.Global().Get("Object").Call("keys", value)
			for i := 0; i < keys.Length(); i++ {
				key := keys.Index(i).String()
				obj[key] = convertJSValueToInterface(value.Get(key))
			}
			return obj
		}
	default:
		return nil
	}
}

func registerCallbacks() {
	js.Global().Set("parse", js.FuncOf(Parse))
	js.Global().Set("restore", js.FuncOf(Restore))
}

func main() {
	registerCallbacks()

	select {} // block forever
}
