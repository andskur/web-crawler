package site

import (
	"net/url"
)

// Site represent Web-site structure
type Site struct {
	Url        *url.URL     `json:"url" xml:"url"`
	TotalPages int          `json:"total_pages" xml:"total_pages"`
	PageTree   *Page        `json:"tree,omitempty" xml:"tree,omitempty"`
	HashMap    PagesHashMap `json:"map,omitempty" xml:"map,omitempty"`
}

// NewSite create new site from given target Url
func NewSite(targetUrl string) (*Site, error) {
	entryPage, err := url.ParseRequestURI(targetUrl)
	if err != nil {
		return nil, err
	}
	site := &Site{
		Url: entryPage,
		PageTree: &Page{
			Url: entryPage,
		},
		HashMap: make(map[string][]string),
	}
	return site, nil
}
