package site

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// Page represent web-site page structure with own URL
// and slice of the links - pointers to other pages
type Page struct {
	Url        *Url          `json:"url" xml:"url"`                              // Page Url
	TotalLinks int           `json:"total,omitempty" xml:"total,omitempty"`      // Total valid links in page
	Links      []*Page       `json:"links,omitempty" xml:"links>page,omitempty"` // Slice of valid pages links in current Page
	Logger     *logrus.Entry `json:"-" xml:"-"`                                  // Page logger with necessary fields
}

// NewPage create new Page structure instance
func NewPage(url *Url) *Page {
	logger := logrus.WithField("page", url.String())
	return &Page{Url: url, Logger: logger}
}

// AddSubPage validate and create Child Page of current Parent page
// Return Child page after successes result
func (p *Page) AddSubPage(link string) (*Page, error) {
	// valid given link string
	if err := p.validateLink(link); err != nil {
		return nil, err
	}

	// get Url from string
	url, err := p.Url.ParseUrl(link)
	if err != nil {
		return nil, errParsedLink
	}

	// validate received Url
	if err := p.validateUrl(url); err != nil {
		return nil, err
	}

	// increase parent totalPage counter
	p.TotalLinks++

	// create new page
	page := NewPage(url)

	// add child page to parent page tree
	p.Links = append(p.Links, page)

	return page, nil
}

// validateUrl validate if given Url is valid
// with to Child Page of current Parent Page
func (p Page) validateUrl(url *Url) error {
	// valid if link has querystring
	if len(url.Query()) > 0 {
		return errQueryLink
	}

	// additional host validation
	if url.Host != p.Url.Host {
		return errExternalLink
	}

	// remove duplicate http & https
	if url.Scheme != p.Url.Scheme {
		return errAlreadyParsed
	}

	// check if parent page already have this sub page
	contain := p.inPage(url.String())
	if contain {
		return errAlreadyInParent
	}
	return nil
}

// validateLink validate if given link string
// can be parsed to Child Page of current Parent Page
func (p Page) validateLink(link string) error {
	// valid if link is email-protected
	if strings.Contains(link, "email-protection") {
		return errEmailProtected
	}

	// check if given link belong to current site
	if !strings.HasPrefix(link, "/") && !strings.Contains(link, p.Url.Host) {
		return errExternalLink
	}
	return nil
}

// inPage checks if Page contain given Child Page
func (p Page) inPage(s string) bool {
	for _, v := range p.Links {
		if s == v.Url.String() || v.Url.String()+"/" == s {
			return true
		}
	}
	return false
}
