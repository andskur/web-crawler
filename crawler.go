package main

import (
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/sirupsen/logrus"
)

var target = "https://monzo.com/"

var mu sync.RWMutex

func main() {
	crawler, err := NewCrawler(target)
	if err != nil {
		logrus.Fatal(err)
	}

	crawler.StartCrawling()

	//spew.Dump(crawler.HashMap)
	spew.Dump(crawler.Site)
	logrus.Info(crawler.Site.TotalPages)
	logrus.Info(crawler.TotalDelay)
}

type Crawler struct {
	TargetUrl  *url.URL
	Site       *Site
	HashMap    map[string][]string
	TotalDelay time.Duration
	wg         sync.WaitGroup
}

type Site struct {
	EntryPage  *Page
	TotalPages int
}

type Page struct {
	Url   *url.URL
	Links []*Page
}

//NewCrawler creates new Crawler structure instance
func NewCrawler(targetUrl string) (*Crawler, error) {
	crawler := &Crawler{
		HashMap: make(map[string][]string),
	}
	crawler.HashMap[targetUrl] = []string{}

	formatUrl, err := url.Parse(target)
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

//StartCrawling starting crawling
func (c *Crawler) StartCrawling() {
	started := time.Now()

	c.wg.Add(1)
	go c.CrawlPage(c.Site.EntryPage)
	c.wg.Wait()

	c.TotalDelay = time.Since(started)
	c.Site.TotalPages = len(c.HashMap)
}

func (c *Crawler) CrawlPage(page *Page) error {
	defer c.wg.Done()
	logrus.Infof("Start crawl %s", page.Url.String())

	resp, err := http.Get(page.Url.String())
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()

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

						mu.Lock()
						c.HashMap[page.Url.String()] = append(c.HashMap[page.Url.String()], childUrl.String())
						mu.Unlock()

						mu.Lock()
						_, ok := c.HashMap[childUrl.String()]
						mu.Unlock()
						if ok {
							continue
						}

						mu.Lock()
						c.HashMap[childUrl.String()] = []string{}
						mu.Unlock()

						childPage := &Page{Url: childUrl}

						page.Links = append(page.Links, childPage)

						c.wg.Add(1)
						go c.CrawlPage(childPage)

						break
					}
				}
			}
		}
	}
}

//isLinkValid check if given link is valid for parsing
func isLinkValid(link, host string) bool {
	if (strings.HasPrefix(link, "/") || strings.Contains(link, host)) && !strings.Contains(link, "email-protection") {
		return true
	}
	return false
}

func removeAnchor(s string) string {
	if idx := strings.Index(s, "/#"); idx != -1 {
		return s[:idx]
	}
	return s
}
