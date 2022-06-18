package files

import (
	"encoding/json"
	"io/ioutil"
)

func JSONToStruct(path string, cache interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, cache)
}
