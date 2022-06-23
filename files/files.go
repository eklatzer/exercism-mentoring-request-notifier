package files

import (
	"encoding/json"
)

//FileReader is used to read a file
type FileReader struct {
	readFile func(string) ([]byte, error)
}

//New returns an instance of FileReader
func New(readFile func(string) ([]byte, error)) *FileReader {
	return &FileReader{readFile: readFile}
}

// JSONToStruct is used to read a JSON from the given file and unmarshal it to the cache
func (r FileReader) JSONToStruct(path string, cache interface{}) error {
	file, err := r.readFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, cache)
}
