package json

import (
	"os"
	"encoding/json"
)

type JsonFile struct {
	Path string
	Root any
}

func LoadJsonFile(path string) (*JsonFile, error) {
	data, err := os.ReadFile(path)
	if err != nil { return nil, err }

	var v any
	err = json.Unmarshal(data, &v)
	return &JsonFile{path, v}, err
}
