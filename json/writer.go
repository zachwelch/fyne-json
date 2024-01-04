package json

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type writer_progress struct {
	i, n int
}
// update progress after processing item; returns true if it was the last
func (wp *writer_progress) update() bool {
	if wp.i < wp.n {
		wp.i += 1
	} else {
		println("warning: writer progress overflowed")
	}
	return wp.i == wp.n
}

type JsonWriter struct {
	JsonWalker
	Writer io.Writer
	stack []writer_progress
}

func (jw *JsonWriter) progress_push(n int) {
	olen := len(jw.stack)
	stack := make([]writer_progress, olen + 1)
	copy(stack, jw.stack)
	stack[olen] = writer_progress{0, n}
	jw.stack = stack
}
func (jw *JsonWriter) progress_pop() {
	n := len(jw.stack)
	if n > 0 {
		jw.stack = jw.stack[0:n-1]
	} else {
		println("Warning: writer stack underflow")
	}
}
func (jw *JsonWriter) progress() *writer_progress {
	last := len(jw.stack) - 1
	if last < 0 {
		return nil
	}
	return &jw.stack[last]
}

func (jw *JsonWriter) separator() {
	if ! jw.progress().update() {
		fmt.Fprint(jw.Writer, ",")
	}
}
func (jw *JsonWriter) Object(o map[string]any) {
	jw.progress_push(len(o))
	fmt.Fprint(jw.Writer, "{")
	jw.WalkObject(o)
	fmt.Fprint(jw.Writer, "}")
	jw.progress_pop()
}
func (jw *JsonWriter) Item(key string, value any) {
	jw.String(key)
	fmt.Fprint(jw.Writer, ":")
	jw.Walk(value)
	jw.separator()
}
func (jw *JsonWriter) Array(o []any) {
	jw.progress_push(len(o))
	fmt.Fprint(jw.Writer, "[")
	jw.WalkArray(o)
	fmt.Fprint(jw.Writer, "]")
}
func (jw *JsonWriter) Element(index int, value any) {
	jw.Walk(value)
	jw.separator()
}
func (jw *JsonWriter) Number(o float64) {
	fmt.Fprintf(jw.Writer, "%f", o)
}
func (jw *JsonWriter) String(o string) {
	fmt.Fprintf(jw.Writer, "\"%s\"", o)
}
func (jw *JsonWriter) Bool(o bool) {
	var t string
	if o { t = "true" } else { t = "false" }
	fmt.Fprintf(jw.Writer, "%s", t)
}

//----------------------------------------------------------------------

type JsonBuilder struct { JsonWriter }

func FormatJson(o any) string {
	sb := &strings.Builder{}

	jb := &JsonBuilder{}
	jb.Visitor = jb
	jb.Writer = sb
	jb.Walk(o)

	return sb.String()
}

//----------------------------------------------------------------------

type JsonPrinter struct { JsonWriter }

func PrintJson(o any) {
	printer := &JsonPrinter{}
	printer.Visitor = printer
	printer.Writer = os.Stdout
	printer.Walk(o)
}
