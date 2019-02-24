package xml

import (
	"encoding/xml"
	"io/ioutil"
)

// WriterJson represent Xml implementation of the IWriter interface
type WriterXml struct{}

// WriteTo writes providing data to given file
func (WriterXml) WriteTo(data interface{}, fileName string) error {
	xmlFormat, err := xml.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fileName, xmlFormat, 0644); err != nil {
		return err
	}
	return nil
}
