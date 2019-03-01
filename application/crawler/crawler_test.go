package crawler

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/andskur/web-crawler/application/site"
)

func TestNewCrawler(t *testing.T) {
	type args struct {
		targetUrl *site.Url
		verbose   bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Crawler
		wantErr bool
	}{
		{"validCralwer", args{getTestSite().Url, false}, &Crawler{Site: getTestSite(), Verbose: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCrawler(tt.args.targetUrl, tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCrawler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCrawler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrawler_CrawlPage(t *testing.T) {
	c, _ := NewCrawler(getTestSite().Url, false)
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

func TestCrawler_calcDuration(t *testing.T) {
	type fields struct {
		Site     *site.Site
		Duration time.Duration
		Verbose  bool
		wg       sync.WaitGroup
	}
	type args struct {
		invocation time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				Site:     tt.fields.Site,
				Duration: tt.fields.Duration,
				Verbose:  tt.fields.Verbose,
				wg:       tt.fields.wg,
			}
			c.calcDuration(tt.args.invocation)
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
