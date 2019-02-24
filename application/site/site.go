package site

import (
	"errors"
	// "net/url"
	"strings"
	"sync"

	"github.com/andskur/web-crawler/utils"
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

// var mu = sync.RWMutex{}

// Site represent Web-site structure
type Site struct {
	Url        *Url          `json:"url" xml:"url"`
	TotalPages int           `json:"total_pages" xml:"total_pages"`
	PageTree   *Page         `json:"tree,omitempty" xml:"tree,omitempty"`
	HashMap    PagesHashMap  `json:"map,omitempty" xml:"map,omitempty"`
	Mu         *sync.RWMutex `json:"-" xml:"-"`
}

// NewSite create new site from given target Url
func NewSite(targetUrl string) (*Site, error) {
	entryPage, err := ParseRequestURI(targetUrl)
	if err != nil {
		return nil, err
	}
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
	uri, err := s.validatePage(link, parent)
	if err != nil {
		return nil, err
	}
	parent.TotalLinks++

	page := NewPage(uri)

	s.Mu.Lock()
	s.HashMap[parent.Url.String()] = append(s.HashMap[parent.Url.String()], page.Url.String())
	s.Mu.Unlock()

	parent.Links = append(parent.Links, page)

	return page, nil
}

// AddPageToSite add given page to current site
func (s *Site) AddPageToSite(page Page) error {
	s.Mu.Lock()
	ok := utils.InMap(page.Url.String(), s.HashMap)
	s.Mu.Unlock()
	if ok {
		return errAlreadyParsed
	}

	s.Mu.Lock()
	s.HashMap[page.Url.String()] = []string{}
	s.Mu.Unlock()

	return nil
}

// TODO damn validation... need to improve this shit

// validatePage validate page before adding to current site
// return Urls struct if valid
func (s Site) validatePage(link string, parent *Page) (*Url, error) {
	if strings.Contains(link, "email-protection") {
		return nil, errEmailProtected
	}
	if err := s.isLinkInHost(link); err != nil {
		return nil, errExternalLink
	}
	uri, err := parent.Url.ParseUrl(link)
	if err != nil {
		return nil, errParsedLink
	}
	if len(uri.Query()) > 0 {
		return nil, errQueryLink
	}
	if uri.Host != s.Url.Host {
		return nil, errExternalLink
	}

	s.Mu.Lock()
	contain := utils.InSlice(uri.String(), s.HashMap[parent.Url.String()])
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
