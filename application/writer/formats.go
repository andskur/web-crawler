package writer

import (
	"fmt"
)

// Format is Enum that represent
// of supported writer types
type Format int

// available Writer constants
const (
	JSON Format = iota
	XML
	unsupported
)

// writers is slice of writer string representations
var formats = [...]string{
	JSON: "json",
	XML:  "xml",
}

// String return writer enum as a string
func (w Format) String() string {
	return formats[w]
}

// ParseFormats return new Format enum from given string
func ParseFormats(s string) (Format, error) {
	for i, r := range formats {
		if s == r {
			return Format(i), nil
		}
	}
	return unsupported, fmt.Errorf("invalid Writer Format value %q", s)
}
