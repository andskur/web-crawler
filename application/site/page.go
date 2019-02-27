package site

import (
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
