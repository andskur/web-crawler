package site

import (
	"encoding/json"
	"encoding/xml"
	"net/url"
)

// Url represent Urls structure
// which inherit basic stdlib url.Url
type Url struct {
	*url.URL
}

// Parse parses a URL in the context of the receiver. The provided URL
// may be relative or absolute. Parse returns nil, err on parse
// failure, otherwise its return value is the same as ResolveReference.
func (u *Url) ParseUrl(ref string) (*Url, error) {
	uri, err := u.Parse(ref)
	if err != nil {
		return nil, err
	}
	return &Url{uri}, nil
}

// ParseRequestURI parses rawUrl into a URL structure. It assumes that
// rawUrl was received in an HTTP request, so the rawUrl is interpreted
// only as an absolute URI or an absolute path.
// The string rawUrl is assumed not to have a #fragment suffix.
// (Web browsers strip #fragment before sending the URL to a web server.)
func ParseRequestURI(rawUrl string) (*Url, error) {
	uri, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return nil, err
	}
	return &Url{uri}, err
}

// MarshalJSON provide corrects Url Json marshaling
func (u Url) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// MarshalXML provide corrects Url Xml marshaling
func (u Url) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(u.String(), start)
}
