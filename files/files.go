package files

import (
	"encoding/json"
)

type FileReader struct {
	readFile func(string) ([]byte, error)
}

func New(readFile func(string) ([]byte, error)) *FileReader {
	return &FileReader{readFile: readFile}
}

func (r FileReader) JSONToStruct(path string, cache interface{}) error {
	file, err := r.readFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, cache)
}
