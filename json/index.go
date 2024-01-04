package json

import (
	"fmt"
	"sort"
	"strings"
)

type JsonChildren map[string][]string

type JsonIndex struct {
	Ids []string
	Items map[string]any
	Children JsonChildren
}

func NewJsonIndex(o any) *JsonIndex {
	index := &JsonIndex{
		Ids: []string{},
		Items: make(map[string]any),
		Children: make(JsonChildren),
	}
	indexer := NewJsonIndexer(index)
	WalkJson(o, indexer)
	return indexer.Index
}

func (ji *JsonIndex) Add(path string, obj any) {
	ji.Ids = append(ji.Ids, path)
	ji.Items[path] = obj
}
func (ji *JsonIndex) AddChildren(path string, children []string) {
	ji.Children[path] = append(ji.Children[path], children...)
}
func (ji *JsonIndex) HasChildren(path string) bool {
	_, ok := ji.Children[path]
	return ok
}

//----------------------------------------------------------------------

type JsonIndexer struct {
	JsonWalker
	Index *JsonIndex
	Paths []string
}

func NewJsonIndexer(index *JsonIndex) *JsonIndexer {
	indexer := &JsonIndexer{
		Index: index,
		Paths: []string{""},
	}
	indexer.Visitor = indexer
	return indexer
}

func (ji *JsonIndexer) WithPath(path string, f func()) {
	old := ji.Paths
	ji.Paths = append(ji.Paths, path)
	f()
	ji.Paths = old
}

func (ji *JsonIndexer) Path() string {
	if len(ji.Paths) == 0 {
		return "."
	}
	return strings.Join(ji.Paths, ".")
}

func (ji *JsonIndexer) ChildPath(child string) string {
	var path string
	ji.WithPath(child, func () { path = ji.Path() })
	return path
}

func (ji *JsonIndexer) Add(o any) {
	ji.Index.Add(ji.Path(), o)
}
func (ji *JsonIndexer) AddChildren(children []string) {
	sort.Strings(children)
	ji.Index.AddChildren(ji.Path(), children)
}
func (ji *JsonIndexer) Object(o map[string]any) {
	ji.Add(o)

	names := make([]string, 0, len(o))
	for name, _ := range(o) {
		names = append(names, ji.ChildPath(name))
	}
	ji.AddChildren(names)

	ji.WalkObject(o)
}
func (ji *JsonIndexer) Item(key string, o any) {
	ji.WithPath(key, func () { WalkJson(o, ji) })
}
func (ji *JsonIndexer) Array(o []any) {
	ji.Add(o)

	names := make([]string, len(o))
	for i, _ := range(o) {
		names[i] = ji.ChildPath(fmt.Sprintf("[%d]", i))
	}
	ji.AddChildren(names)

	ji.WalkArray(o)
}
func (ji *JsonIndexer) Element(index int, o any) {
	path := fmt.Sprintf("[%d]", index)
	ji.WithPath(path, func () { WalkJson(o, ji) })
}
func (ji *JsonIndexer) String(o string) { ji.Add(o) }
func (ji *JsonIndexer) Number(o float64) { ji.Add(o) }
func (ji *JsonIndexer) Bool(o bool) { ji.Add(o) }
