package crawler

import (
	"fmt"
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
	Site     *site.Site     // web site for crawling
	Duration time.Duration  // total crawling duration
	Verbose  bool           // verbose mode
	wg       sync.WaitGroup // crawler waitgroup
}

// NewCrawler creates new Crawler structure instance
func NewCrawler(targetUrl *site.Url, verbose bool) (*Crawler, error) {
	crawler := &Crawler{
		Site:    site.NewSite(targetUrl),
		Verbose: verbose,
	}

	return crawler, nil
}

// StartCrawling starting crawling
func (c *Crawler) StartCrawling() {
	// print Crawler result after its execution
	defer c.PrintResult()

	// calculate total duration
	defer c.calcDuration(time.Now())

	fmt.Printf("Start crawling web site %s...\n", c.Site.Url.Host)

	// start crawling site pages
	c.wg.Add(1)
	go func() {
		if err := c.CrawlPage(c.Site.PageTree); err != nil && c.Verbose {
			logrus.Error(err)
		}
	}()

	done := make(chan struct{})

	// if verbose disabled - print total pages count concurrency
	if !c.Verbose {
		go c.printTotal(done)
	}

	// waiting finish crawling of all site pages
	c.wg.Wait()
	close(done)
}

// CrawlPage crawl given site page
func (c *Crawler) CrawlPage(page *site.Page) error {
	defer c.wg.Done()

	// increase total site pages count
	c.Site.TotalPages++

	if c.Verbose {
		page.Logger.Info("Start page crawling...")
	}

	// http request too new crawling page
	resp, err := http.Get(page.Url.String())
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			if c.Verbose {
				page.Logger.Error(err)
			}
			return
		}
	}()

	// FIXME need to find better way for check page format
	// check response format, need only tex/html for next crawling
	if contentType := resp.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "text/html") {
		// if page is not text/html - delete it from site Hash Map and decrease total site pages count
		c.Site.Mu.Lock()
		delete(c.Site.HashMap, page.Url.String())
		c.Site.Mu.Unlock()
		c.Site.TotalPages--
		return fmt.Errorf("unsupported page format - %s", contentType)
	}

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
						childPage, err := page.AddSubPage(link)
						if err != nil {
							// TODO temporarily disabling, need to implement logging levels
							/*if c.Verbose {
								page.Logger.WithField("link", link).Error(err)
							}*/
							continue
						}

						// add child page to parent links slice
						c.Site.Mu.Lock()
						c.Site.HashMap[page.Url.String()] = append(c.Site.HashMap[page.Url.String()], childPage.Url.String())
						c.Site.Mu.Unlock()

						// validate and add page to site
						if err := c.Site.AddPageToSite(*childPage); err != nil {
							// TODO temporarily disabling, need to implement logging levels
							/*if c.Verbose {
								childPage.Logger.Error(err)
							}*/
							continue
						}

						// start crawl child page
						c.wg.Add(1)
						go func() {
							if err := c.CrawlPage(childPage); err != nil && c.Verbose {
								childPage.Logger.Error(err)
							}
						}()
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

// printTotal concurrently print total crawled site pages
func (c *Crawler) printTotal(done chan struct{}) {
	for {
		select {
		case <-done:
			goto Finish
		default:
			fmt.Printf("\rTotal pages: %d...", c.Site.TotalPages)
		}
	}
Finish:
	fmt.Println("\nAll done!")
}

// PrintResult print Crawler results
func (c *Crawler) PrintResult() {
	fmt.Printf("%d pages crawled at %s in %s\n", c.Site.TotalPages, c.Site.Url.Host, c.Duration)
}

// removeAnchor remove anchor from given string link
func removeAnchor(s string) string {
	if idx := strings.Index(s, "/#"); idx != -1 {
		return s[:idx]
	}
	return s
}
