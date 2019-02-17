package main

import (
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

	//spew.Dump(crawler.SiteMap)
	logrus.Info(len(crawler.SiteMap))
	logrus.Info(crawler.TotalDelay)
}

type Crawler struct {
	TargetUrl  *url.URL
	SiteMap    map[string][]string
	TotalDelay time.Duration
	wg         sync.WaitGroup
}

//NewCrawler creates new Crawler structure instance
func NewCrawler(targetUrl string) (*Crawler, error) {
	crawler := &Crawler{
		SiteMap: make(map[string][]string),
	}
	crawler.SiteMap[targetUrl] = []string{}

	formatUrl, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	crawler.TargetUrl = formatUrl

	return crawler, nil
}

//StartCrawling starting crawling
func (c *Crawler) StartCrawling() {
	started := time.Now()

	c.wg.Add(1)
	go c.CrawlUrl(*c.TargetUrl)
	c.wg.Wait()

	c.TotalDelay = time.Since(started)
}

func (c *Crawler) CrawlUrl(parentUrl url.URL) error {
	defer c.wg.Done()
	logrus.Infof("Start crawl %s\n", parentUrl.String())

	resp, err := http.Get(parentUrl.String())
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()

	tokens := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokens.Next()

		switch tokenType {
		case html.ErrorToken:
			return nil
		case html.StartTagToken, html.EndTagToken:
			token := tokens.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if link := removeAnchor(attr.Val); isLinkValid(link, c.TargetUrl.String()) && attr.Key == "href" {
						childUrl, err := parentUrl.Parse(link)
						if err != nil {
							return err
						}

						mu.Lock()
						c.SiteMap[parentUrl.String()] = append(c.SiteMap[parentUrl.String()], childUrl.String())
						mu.Unlock()

						mu.Lock()
						_, ok := c.SiteMap[childUrl.String()]
						mu.Unlock()
						if ok {
							continue
						}

						mu.Lock()
						c.SiteMap[childUrl.String()] = []string{}
						mu.Unlock()

						c.wg.Add(1)
						go c.CrawlUrl(*childUrl)
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
