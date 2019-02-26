package site

import (
	"encoding/xml"
	"errors"
	"strings"
	"sync"
)

// Site pages validation errors
var (
	errQueryLink       = errors.New("link with query params")
	errParsedLink      = errors.New("cannot parsing link")
	errExternalLink    = errors.New("link is external")
	errAlreadyInParent = errors.New("link already in parent slice")
	errAlreadyParsed   = errors.New("link have already parsed")
	errEmailProtected  = errors.New("link is email-protected")
)

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
func NewSite(entryPage *Url) (*Site, error) {
	/*entryPage, err := ParseRequestURI(targetUrl)
	if err != nil {
		return nil, err
	}*/
	site := &Site{
		Url: entryPage,
		PageTree: &Page{
			Url: entryPage,
		},
		HashMap: make(map[string][]string),
		Mu:      &sync.RWMutex{},
	}
	return site, nil
}

// AddPageToParent validate page and add it to given parent
func (s *Site) AddPageToParent(link string, parent *Page) (*Page, error) {
	// validate page
	uri, err := s.validatePage(link, parent)
	if err != nil {
		return nil, err
	}
	// increase parent totalPage counter
	parent.TotalLinks++

	// create new page
	page := &Page{Url: uri}

	// add child page to parent links slice
	s.Mu.Lock()
	s.HashMap[parent.Url.String()] = append(s.HashMap[parent.Url.String()], page.Url.String())
	s.Mu.Unlock()

	// add child page to parent page tree
	parent.Links = append(parent.Links, page)

	return page, nil
}

// AddPageToSite add given page to current site
func (s *Site) AddPageToSite(page Page) error {
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

// TODO damn validation... need to improve this shit

// validatePage validate page before adding to current site
// return Urls struct if valid
func (s Site) validatePage(link string, parent *Page) (*Url, error) {
	// valid if link is email-protected
	if strings.Contains(link, "email-protection") {
		return nil, errEmailProtected
	}
	// valid if link belong to same host as the current site
	if err := s.isLinkInHost(link); err != nil {
		return nil, errExternalLink
	}
	// get Url from string
	uri, err := parent.Url.ParseUrl(link)
	if err != nil {
		return nil, errParsedLink
	}
	// valid if link has querystring
	if len(uri.Query()) > 0 {
		return nil, errQueryLink
	}
	// additional host validation
	if uri.Host != s.Url.Host {
		return nil, errExternalLink
	}

	// check if link already have in parent slice
	s.Mu.Lock()
	contain := inSlice(uri.String(), s.HashMap[parent.Url.String()])
	s.Mu.Unlock()
	if contain {
		return nil, errAlreadyInParent
	}
	return uri, nil
}

// isLinkInHost check if given link belong to current site
func (s Site) isLinkInHost(link string) error {
	if strings.HasPrefix(link, "/") || strings.Contains(link, s.Url.String()) {
		return nil
	}
	return errExternalLink
}

// inSlice checks if slice contain given string
func inSlice(s string, slice []string) bool {
	for _, v := range slice {
		if s == v || v+"/" == s {
			return true
		}
	}
	return false
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
