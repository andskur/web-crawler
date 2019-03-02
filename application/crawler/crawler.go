package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/andskur/web-crawler/application/site"
)

// Crawler represent web-crawler structure
type Crawler struct {
	Site      *site.Site     // web site for crawling
	Duration  time.Duration  // total crawling duration
	Semaphore chan int       // thread-blocking Semaphore channel
	Verbose   bool           // verbose mode
	wg        sync.WaitGroup // crawler WaitGroup
}

// NewCrawler creates new Crawler structure instance
func NewCrawler(targetURL *site.Url, verbose bool, semaphore chan int) (*Crawler, error) {
	crawler := &Crawler{
		Site:      site.NewSite(targetURL),
		Verbose:   verbose,
		Semaphore: semaphore,
	}
	return crawler, nil
}

// StartCrawling starting crawling
func (c *Crawler) StartCrawling() error {
	// print Crawler result after its execution
	defer c.PrintResult()

	// calculate total duration
	defer c.calcDuration(time.Now())

	fmt.Printf("Start crawling web site %s...\n", c.Site.Url.Host)

	// start crawling site pages
	c.wg.Add(1)
	<-c.Semaphore
	if err := c.CrawlPage(c.Site.PageTree); err != nil && c.Verbose {
		return err
	}

	// create "done: channel
	done := make(chan struct{})

	// if verbose disabled - print total pages count concurrency
	if !c.Verbose {
		go c.printTotal(done)
	}

	// waiting finish crawling of all site pages
	c.wg.Wait()
	close(done)
	return nil
}

// CrawlPage crawl given site page
func (c *Crawler) CrawlPage(page *site.Page) error {
	defer c.wg.Done()

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

	// TODO need to find better way for check page format
	// check response format, need only tex/html for next crawling
	if contentType := resp.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "text/html") {
		// if page is not text/html - delete it from site Hash Map
		c.Site.DeletePageFromSite(page.Url.String())
		c.Semaphore <- 1
		return fmt.Errorf("unsupported page format - %s", contentType)
	}

	// increase total site pages count
	c.Site.TotalPages++

	// parse html body
	tokens := html.NewTokenizer(resp.Body)

	c.Semaphore <- 1

	// find valid html tags
	for {
		switch tokens.Next() {
		case html.ErrorToken:
			return nil
		case html.StartTagToken:
			// we need only <a> html tag
			token := tokens.Token()
			if token.Data != "a" {
				continue
			}

			// get link from href attribute
			link, ok := getLink(token)
			if !ok {
				continue
			}

			// validate and add child page to parent page
			childPage, err := page.AddSubPage(link)
			if err != nil {
				// TODO need to implement logging levels
				/*if c.Verbose {
					page.Logger.WithField("link", link).Error(err)
				}*/
				continue
			}

			// add child page to parent links slice
			c.Site.AddPageToParent(childPage.Url.String(), page.Url.String())

			// validate and add page to site
			if err := c.Site.AddPageToSite(childPage.Url.String()); err != nil {
				// TODO need to implement logging levels
				/*if c.Verbose {
					childPage.Logger.Error(err)
				}*/
				continue
			}

		CrawlChild:
			// check if Crawler have available threads
			select {
			case <-c.Semaphore:
				// start crawl child page if we have
				c.wg.Add(1)
				go func() {
					if err := c.CrawlPage(childPage); err != nil && c.Verbose {
						childPage.Logger.Error(err)
					}
				}()
			default:
				// print warning and try again after one second
				if c.Verbose {
					childPage.Logger.Warning("Threads limit reached... WAIT")
				}
				time.Sleep(10 * time.Millisecond)
				goto CrawlChild
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

// getLink get href Link from attribute of given html tag
func getLink(token html.Token) (link string, ok bool) {
	for _, attr := range token.Attr {
		// finds"href" attribute and remove anchor
		if attr.Key == "href" {
			link = removeAnchor(attr.Val)
			ok = true
		}
	}
	return
}
