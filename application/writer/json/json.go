package json

import (
	"encoding/json"
	"io/ioutil"
)

// WriterJson represent Json implementation of the IWriter interface
type WriterJson struct{}

// WriteTo writes providing data to given file
func (WriterJson) WriteTo(data interface{}, fileName string) error {
	jsonFormat, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fileName, jsonFormat, 0644); err != nil {
		return err
	}
	return nil
}
