package crawler

import (
	"testing"

	"github.com/andskur/web-crawler/application/site"
)

func TestCrawler_CrawlPage(t *testing.T) {
	c, _ := NewCrawler(getTestSite().Url, false, make(chan int, 10000))
	type args struct {
		page *site.Page
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"successCrawling", args{c.Site.PageTree}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.wg.Add(1)
			if err := c.CrawlPage(tt.args.page); (err != nil) != tt.wantErr {
				t.Errorf("Crawler.CrawlPage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_removeAnchor(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"anchorFullLink", args{"https://monzo.com/blog/#gotobotton"}, "https://monzo.com/blog"},
		{"anchorPath", args{"/blog/#gotobotton"}, "/blog"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeAnchor(tt.args.s); got != tt.want {
				t.Errorf("removeAnchor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestSite() *site.Site {
	url, _ := site.ParseRequestURI("https://ya.ru")
	site := site.NewSite(url)
	return site
}
