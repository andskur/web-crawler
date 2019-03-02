package site

import (
	"encoding/xml"
	"errors"
	"strings"
	"sync"
)

var errAlreadyParsed = errors.New("page have already parsed")

// Site represent Web-site structure
type Site struct {
	XMLName    xml.Name     `json:"-" xml:"site"`
	Url        *Url         `json:"url" xml:"url"`                       // basic site Url
	TotalPages int          `json:"total_pages" xml:"total_pages"`       // total counts site page
	PageTree   *Page        `json:"tree,omitempty" xml:"tree,omitempty"` // site page tree
	HashMap    PagesHashMap `json:"map,omitempty" xml:"map,omitempty"`   // site hash page map
	mu         *sync.Mutex  `json:"-" xml:"-"`                           // mutex variable for threadsafe operations with maps
}

// NewSite create new site from given target Url
func NewSite(entryPage *Url) *Site {
	return &Site{
		Url:      entryPage,
		PageTree: NewPage(entryPage),
		HashMap:  make(map[string][]string),
		mu:       &sync.Mutex{},
	}
}

// TODO better way - move it to Page's methods

// AddPageToParent add Child Page to Parent Page slice in HashMap
func (s *Site) AddPageToParent(child, parent string) {
	s.mu.Lock()
	s.HashMap[parent] = append(s.HashMap[parent], child)
	s.mu.Unlock()
}

// AddPageToSite validate and add given page to current site
func (s *Site) AddPageToSite(page string) error {
	// check if page already in main hash map
	s.mu.Lock()
	ok := inMap(page, s.HashMap)
	s.mu.Unlock()
	if ok {
		return errAlreadyParsed
	}

	// add page to main hash map
	s.mu.Lock()
	s.HashMap[page] = []string{}
	s.mu.Unlock()

	return nil
}

// DeletePageFromSite delete given page from Site
func (s *Site) DeletePageFromSite(page string) {
	s.mu.Lock()
	delete(s.HashMap, page)
	s.mu.Unlock()
}

// TODO need refactoring

// inMap check if map contain given link
func inMap(s string, m map[string][]string) bool {
	_, ok := m[s]
	_, okSlash := m[s+"/"]
	_, okOneMore := m[strings.TrimSuffix(s, "/")]
	if okSlash || ok || okOneMore {
		return true
	}
	return false
}
