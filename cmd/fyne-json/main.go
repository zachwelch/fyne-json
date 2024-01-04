package main

import (
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	fj "github.com/zachwelch/fyne-json/json"
)

func NewJsonView(jf *fj.JsonFile) fyne.CanvasObject {
	index := fj.NewJsonIndex(jf.Root)
	_ = fj.NewJsonIndex(map[string]any{
			"a": 56.,
			"b": map[string]any{
				"1": "string",
				"2": true,
			},
			"c": []any{"42", 42.},
		})

	path := widget.NewLabel(jf.Path)
	node := widget.NewLabel(".")
	top := container.NewVBox(
		container.NewHBox(widget.NewLabel("Path:"), path),
		container.NewHBox(widget.NewLabel("Node:"), node),
		widget.NewSeparator(),
	)

	t := widget.NewTree(
		func (id widget.TreeNodeID) []widget.TreeNodeID {
//			println("node: " + id)
			return index.Children[id]
		},
		func (id widget.TreeNodeID) bool {
//			println("branch: " + id)
			return index.HasChildren(id)
		},
		func (bool) fyne.CanvasObject {
			name := widget.NewLabel("---------------")
			value := widget.NewLabel("---------------")
			value.Wrapping = fyne.TextTruncate
			s := container.NewHSplit(name, value)
			s.SetOffset(0.25)
			return s
		},
		func (id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
//			println("update: " + id)
			json := ""
			if ! index.HasChildren(id) {
				json = fj.FormatJson(index.Items[id])
			}
			nids := strings.Split(id, ".")
			nid := nids[len(nids) - 1]
			s := o.(*container.Split)
			s.Leading.(*widget.Label).SetText(nid)
			s.Trailing.(*widget.Label).SetText(json)
		},
	)
	return container.NewBorder(top, nil, nil, nil, t)
}

//----------------------------------------------------------------------

func main() {
	if len(os.Args) != 2 {
		println("usage: " + os.Args[0] + " <jsonfile>")
		os.Exit(1)
	}
	name := os.Args[1]

	jf, err := fj.LoadJsonFile(name)
	if err != nil {
		fyne.LogError(name + ":unable to load JSON file", err)
		return
	}
	//PrintJson(jf.Root)

	jv := NewJsonView(jf)

	a := app.New()
	w := a.NewWindow("JSON Viewer")

	w.Resize(fyne.Size{800,600})
	w.SetContent(jv)

	w.ShowAndRun()
}
