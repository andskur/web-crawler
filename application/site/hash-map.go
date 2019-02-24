package site

import (
	"encoding/json"
	"encoding/xml"
	"sort"
)

// PagesHashMap represent Pages Hash Map structure type
type PagesHashMap map[string][]string

// MarshalJSON correct formatted JSON marshaling
// for Page Hash Map structure type
func (p PagesHashMap) MarshalJSON() ([]byte, error) {
	pages := p.mapToHashPages()
	return json.Marshal(pages)
}

// MarshalXML correct formatted XML marshaling
// for Page Hash Map structure type
func (p PagesHashMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// FIXME need to hide empty list fields
	pages := p.mapToHashPages()

	err := e.EncodeElement(struct {
		XMLName xml.Name    `xml:"page"`
		Pages   *[]hashPage `xml:"pages"`
	}{
		Pages: pages}, start)
	if err != nil {
		return err
	}
	return nil
}

// hashPage represent PagesHashMap formatter for XML and JSON marshaling
type hashPage struct {
	XMLName    xml.Name  `json:"-" xml:"page"`
	Url        string    `json:"url "xml:"url"`
	TotalLinks int       `json:"total_links" xml:"total_links"`
	Links      *[]string `json:"links" xml:"links>url,omitempty"`
}

// mapToHashPages create slice of hashPage from PagesHashMap
func (p PagesHashMap) mapToHashPages() *[]hashPage {
	var pages []hashPage
	for url, links := range p {
		page := hashPage{Url: url}
		var lks []string
		for _, link := range links {
			lks = append(lks, link)
		}
		page.TotalLinks = len(lks)
		page.Links = &lks
		pages = append(pages, page)
	}
	sort.Slice(pages, func(i, j int) bool {
		return len(pages[i].Url) < len(pages[j].Url)
	})
	return &pages
}
