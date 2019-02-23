package site

import (
	"encoding/json"
	"net/url"
)

// Site represent Web-site structure
type Site struct {
	EntryPage  *Page `json:"entry_page"`
	TotalPages int   `json:"total_pages"`
}

// Page represent web-site page structure with own URL
// and slice of the links - pointers to other pages
type Page struct {
	Url   *url.URL `json:"url"`
	Links []*Page  `json:"links,omitempty"`
}

// MarshalJSON corrects Json marshaling
// for page structure type
func (p Page) MarshalJSON() ([]byte, error) {
	basicSite := struct {
		Url   string  `json:"url"`
		Links []*Page `json:"links"`
	}{
		Url:   p.Url.String(),
		Links: p.Links,
	}

	return json.Marshal(basicSite)
}
