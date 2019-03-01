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
	XMLName    xml.Name      `json:"-" xml:"site"`
	Url        *Url          `json:"url" xml:"url"`                       // basic site Url
	TotalPages int           `json:"total_pages" xml:"total_pages"`       // total counts site page
	PageTree   *Page         `json:"tree,omitempty" xml:"tree,omitempty"` // site page tree
	HashMap    PagesHashMap  `json:"map,omitempty" xml:"map,omitempty"`   // site hash page map
	Mu         *sync.RWMutex `json:"-" xml:"-"`                           // global Read/Write mutex variable for threadsafe operations with maps
}

// NewSite create new site from given target Url
func NewSite(entryPage *Url) *Site {
	return &Site{
		Url:      entryPage,
		PageTree: NewPage(entryPage),
		HashMap:  make(map[string][]string),
		Mu:       &sync.RWMutex{},
	}
}

// AddPageToSite add given page to current site
func (s *Site) AddPageToSite(page *Page) error {
	// check if page already in main hash map
	s.Mu.Lock()
	ok := inMap(page.Url.String(), s.HashMap)
	s.Mu.Unlock()
	if ok {
		return errAlreadyParsed
	}

	// add page to main hash map
	s.Mu.Lock()
	s.HashMap[page.Url.String()] = []string{}
	s.Mu.Unlock()

	return nil
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
