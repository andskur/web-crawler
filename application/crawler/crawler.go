package crawler

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"

	"github.com/andskur/web-crawler/application/site"
)

// Crawler represent web-crawler structure
type Crawler struct {
	Site     *site.Site
	Duration time.Duration
	wg       sync.WaitGroup
}

// NewCrawler creates new Crawler structure instance
func NewCrawler(targetUrl *site.Url) (*Crawler, error) {
	crawSite, err := site.NewSite(targetUrl)
	if err != nil {
		return nil, err
	}
	crawler := &Crawler{Site: crawSite}

	return crawler, nil
}

// StartCrawling starting crawling
func (c *Crawler) StartCrawling() {
	defer c.calcDuration(time.Now())

	c.wg.Add(1)
	go c.CrawlPage(c.Site.PageTree)
	c.wg.Wait()

	c.Site.TotalPages = len(c.Site.HashMap)
}

// CrawlPage crawl given site page
func (c *Crawler) CrawlPage(page *site.Page) error {
	defer c.wg.Done()
	logrus.Infof("Start crawl %s", page.Url.String())

	// http request too new crawling page
	resp, err := http.Get(page.Url.String())
	if err != nil {
		return err
	}

	// check response format, need only tex/html for next crawling
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		c.Site.Mu.Lock()
		delete(c.Site.HashMap, page.Url.String())
		c.Site.Mu.Unlock()
		return errors.New("unsupported page format")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	// parse html body
	tokens := html.NewTokenizer(resp.Body)

	// find valid html tags
	for {
		switch tokens.Next() {
		case html.ErrorToken:
			return nil
		case html.StartTagToken, html.EndTagToken:
			// we need only <a> html tag
			if token := tokens.Token(); token.Data == "a" {
				for _, attr := range token.Attr {
					// and only "href" attribute and remove anchor
					if link := removeAnchor(attr.Val); attr.Key == "href" {
						// validate and add child page to parent page
						childPage, err := c.Site.AddPageToParent(link, page)
						if err != nil {
							// logrus.Error(err)
							continue
						}

						// validate and add page to site
						if err := c.Site.AddPageToSite(*childPage); err != nil {
							// logrus.Error(err)
							continue
						}

						// start crawl child page
						c.wg.Add(1)
						go c.CrawlPage(childPage)
					}
				}
			}
		}
	}
}

// duration calculate total Crawler execution time
func (c *Crawler) calcDuration(invocation time.Time) {
	c.Duration = time.Since(invocation)
}

// removeAnchor remove anchor from given string link
func removeAnchor(s string) string {
	if idx := strings.Index(s, "/#"); idx != -1 {
		return s[:idx]
	}
	return s
}
