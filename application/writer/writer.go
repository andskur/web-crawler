package writer

import (
	"errors"

	"github.com/andskur/web-crawler/application/writer/json"
	"github.com/andskur/web-crawler/application/writer/xml"
)

var ErrUnsupportedWriter = errors.New("unsupported writer type")

// IWriter represent writer data to file interface
type IWriter interface {
	WriteTo(data interface{}, fileName string) error
}

// NewWriter create new writer instance
func NewWriter(wtype Format) (wrt IWriter, err error) {
	switch wtype {
	case JSON:
		wrt = json.WriterJson{}
	case XML:
		wrt = xml.WriterXml{}
	default:
		err = ErrUnsupportedWriter
	}
	return
}
