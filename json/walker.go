package json

import (
	"fmt"
	"sort"
)

type JsonVisitor interface {
	Object(value map[string]any)
	Item(key string, value any)

	Array(value []any)
	Element(index int, value any)

	String(value string)
	Number(value float64)
	Bool(value bool)
}

func WalkJson(o any, visitor JsonVisitor) {
	switch t := o.(type) {
	case bool:
		visitor.Bool(t)
	case string:
		visitor.String(t)
	case float64:
		visitor.Number(t)
	case []any:
		visitor.Array(t)
	case map[string]any:
		visitor.Object(t)
	default:
		fmt.Println(t, "is an unknown type")
	}
}

func WalkObject(obj map[string]any, visitor JsonVisitor) {
	for k, v := range(obj) {
		visitor.Item(k, v)
	}
}
func WalkObjectOrdered(obj map[string]any, visitor JsonVisitor) {
	keys := make([]string, 0, len(obj))
	for k, _ := range(obj) {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range(keys) {
		visitor.Item(k, obj[k])
	}
}

func WalkArray(obj []any, visitor JsonVisitor) {
	for i, v := range(obj) {
		visitor.Element(i, v)
	}
}

//----------------------------------------------------------------------

type JsonWalker struct {
	Visitor JsonVisitor
}

func (jw *JsonWalker) Walk(o any) {
	WalkJson(o, jw.Visitor)
}
func (jw *JsonWalker) WalkObject(obj map[string]any) {
	WalkObject(obj, jw.Visitor)
}
func (jw *JsonWalker) WalkArray(obj []any) {
	WalkArray(obj, jw.Visitor)
}
