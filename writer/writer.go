package writer

import (
	"errors"

	"github.com/andskur/web-crawler/writer/json"
	"github.com/andskur/web-crawler/writer/xml"
)

var ErrUnsupportedWriter = errors.New("unsupported writer type")

// IWriter represent writer data to file interface
type IWriter interface {
	WriteTo(data interface{}, fileName string) error
}

// NewWriter create new writer instance
func NewWriter(wtype string) (wrt IWriter, err error) {
	switch wtype {
	case "json":
		wrt = json.WriterJson{}
	case "xml":
		wrt = xml.WriterXml{}
	default:
		err = ErrUnsupportedWriter
	}
	return
}
