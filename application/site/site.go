package site

import (
	"encoding/json"
	"encoding/xml"
	"net/url"
)

// SiteTree represent Web-site structure
type Site struct {
	EntryPage  *Page `json:"entry_page" xml:"entry_page"`
	TotalPages int   `json:"total_pages" xml:"total_pages"`
}

// Page represent web-site page structure with own URL
// and slice of the links - pointers to other pages
type Page struct {
	Url        *url.URL
	TotalLinks int
	Links      []*Page
}

// MarshalJSON corrects Json marshaling
// for page structure type
func (p Page) MarshalJSON() ([]byte, error) {
	page := struct {
		Url        string  `json:"url"`
		TotalLinks int     `json:"total_links"`
		Links      []*Page `json:"links,omitempty"`
	}{
		Url:        p.Url.String(),
		TotalLinks: p.TotalLinks,
		Links:      p.Links,
	}
	return json.Marshal(page)
}

// MarshalXML corrects XML marshaling
// for Page Tree structure type
func (p Page) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// FIXME need to hide empty list fields
	err := e.EncodeElement(struct {
		XMLName xml.Name `xml:"page"`
		Url     string   `xml:"url"`
		Links   []*Page  `xml:"links>page,omitempty"`
	}{
		Url:   p.Url.String(),
		Links: p.Links}, start)
	if err != nil {
		return err
	}
	return nil
}
