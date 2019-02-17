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
var wg sync.WaitGroup
var mu sync.RWMutex

func main() {
	entryUrl, err := url.Parse(target)
	if err != nil {
		logrus.Fatal(err)
	}

	siteMap := make(map[string]string)

	siteMap[target] = ""

	started := time.Now()
	wg.Add(1)
	go CrawlUrl(*entryUrl, siteMap)

	wg.Wait()

	timeSpent := time.Since(started)

	logrus.Info(len(siteMap))
	logrus.Info(timeSpent)

}

func CrawlUrl(url url.URL, siteMap map[string]string) error {
	defer wg.Done()

	resp, err := http.Get(url.String())
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
					if link := attr.Val; isLinkValid(link, target) && attr.Key == "href" {
						childUrl, err := url.Parse(link)
						if err != nil {
							return err
						}

						mu.Lock()
						_, ok := siteMap[childUrl.String()]
						mu.Unlock()
						if ok {
							continue
						}

						mu.Lock()
						siteMap[childUrl.String()] = ""
						mu.Unlock()

						//spew.Dump(childUrl.String())
						wg.Add(1)
						//fmt.Println(siteMap)
						go CrawlUrl(*childUrl, siteMap)
					}
				}
			}
		}
	}
}

func isLinkValid(link, host string) bool {
	if (strings.HasPrefix(link, "/") || strings.Contains(link, host)) && !strings.Contains(link, "email-protection") {
		return true
	}
	return false
}
