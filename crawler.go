package crawler

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

// global Read/Write mutex variable
// for threadsafe operations with maps
var mu sync.RWMutex

// Crawler represent web-crawler structure
type Crawler struct {
	TargetUrl  *url.URL
	Site       *Site
	HashMap    map[string][]string
	TotalDelay time.Duration
	wg         sync.WaitGroup
}

// NewCrawler creates new Crawler structure instance
func NewCrawler(targetUrl string) (*Crawler, error) {
	crawler := &Crawler{
		HashMap: make(map[string][]string),
	}
	crawler.HashMap[targetUrl] = []string{}

	formatUrl, err := url.ParseRequestURI(targetUrl)
	if err != nil {
		return nil, err
	}
	crawler.TargetUrl = formatUrl

	crawler.Site = &Site{
		EntryPage: &Page{
			Url: formatUrl,
		},
	}
	return crawler, nil
}

// StartCrawling starting crawling
func (c *Crawler) StartCrawling() {
	started := time.Now()

	c.wg.Add(1)
	go c.CrawlPage(c.Site.EntryPage)
	c.wg.Wait()

	c.TotalDelay = time.Since(started)
	c.Site.TotalPages = len(c.HashMap)
}

// CrawlPage crawl given site page
func (c *Crawler) CrawlPage(page *Page) error {
	defer c.wg.Done()
	logrus.Infof("Start crawl %s", page.Url.String())

	resp, err := http.Get(page.Url.String())
	if err != nil {
		return err
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		mu.Lock()
		delete(c.HashMap, page.Url.String())
		mu.Unlock()
		return errors.New("unsupported page format")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	tokens := html.NewTokenizer(resp.Body)

	for {
		switch tokens.Next() {
		case html.ErrorToken:
			return nil
		case html.StartTagToken, html.EndTagToken:
			if token := tokens.Token(); token.Data == "a" {
				for _, attr := range token.Attr {
					if link := removeAnchor(attr.Val); isLinkValid(link, c.TargetUrl.String()) && attr.Key == "href" {
						childUrl, err := page.Url.Parse(link)
						if err != nil {
							return err
						}

						if len(childUrl.Query()) > 0 {
							continue
						}

						if childUrl.Host != page.Url.Host {
							continue
						}

						mu.Lock()
						contain := inSlice(childUrl.String(), c.HashMap[page.Url.String()])
						mu.Unlock()
						if contain {
							continue
						}

						mu.Lock()
						c.HashMap[page.Url.String()] = append(c.HashMap[page.Url.String()], childUrl.String())
						mu.Unlock()

						childPage := &Page{Url: childUrl}
						page.Links = append(page.Links, childPage)

						mu.Lock()
						ok := inMap(childUrl.String(), c.HashMap)
						mu.Unlock()
						if ok {
							continue
						}

						mu.Lock()
						c.HashMap[childUrl.String()] = []string{}
						mu.Unlock()

						// c.wg.Add(1)
						// go c.CrawlPage(childPage)
					}
				}
			}
		}
	}
}

// TODO move validation to Page methods

// isLinkValid check if given link is valid for parsing
func isLinkValid(link, host string) bool {
	if (strings.HasPrefix(link, "/") || strings.Contains(link, host)) && !strings.Contains(link, "email-protection") {
		return true
	}
	return false
}

// removeAnchor remove anchor from given string link
func removeAnchor(s string) string {
	if idx := strings.Index(s, "/#"); idx != -1 {
		return s[:idx]
	}
	return s
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
	/*if len(m[s]) > 0 || len(m[s + "/"]) > 0 || len(m[strings.TrimSuffix(s, "/")]) > 0 {
		return true
	}
	return false*/
	_, ok := m[s]
	_, okSlash := m[s+"/"]
	_, okOneMore := m[strings.TrimSuffix(s, "/")]
	if okSlash || ok || okOneMore {
		return true
	}
	return false
}
