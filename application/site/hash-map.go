package site

import (
	"encoding/xml"
	"sort"
)

// PagesHashMap represent Pages Hash Map structure type
type PagesHashMap map[string][]string

// hashPage represent PagesHashMap formatter for XML marshaling
type hashPage struct {
	XMLName xml.Name  `xml:"page"`
	Url     string    `xml:"url"`
	Links   *[]string `xml:"links>url,omitempty"`
}

// MarshalXML corrects XML marshaling
// for Page Hash Map structure type
func (p PagesHashMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// FIXME need to hide empty list fields
	var pages []hashPage
	for url, links := range p {
		page := hashPage{Url: url}
		var lks []string
		for _, link := range links {
			lks = append(lks, link)
		}
		page.Links = &lks
		pages = append(pages, page)
	}

	sort.Slice(pages, func(i, j int) bool {
		return len(pages[i].Url) < len(pages[j].Url)
	})

	err := e.EncodeElement(struct {
		XMLName xml.Name   `xml:"page"`
		Pages   []hashPage `xml:"pages"`
	}{
		Pages: pages}, start)
	if err != nil {
		return err
	}
	return nil
}
