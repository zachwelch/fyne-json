package main

import (
	"errors"
	"fmt"
	"os"

	fj "github.com/zachwelch/fyne-json/json"
)

func main() {
	args := os.Args
//	cmdname := args[0]
	n_args := len(args)
	if n_args < 2 {
		panic(errors.New("missing filter argument"))
	}
	if n_args < 3 {
		println("reading from stdin...")
		args = append(args, "/dev/stdin")
	}

	filter := args[1]
	if filter == "." {
		filter = ""
	}

	for _, name := range(args[2:]) {
		jf, err := fj.LoadJsonFile(name)
		if err != nil { panic(err) }

		index := fj.NewJsonIndex(jf.Root)
		item := index.Items[filter]
		if item != nil {
			fj.PrintJson(item)
			fmt.Println()
		}
	}
}
