package site

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewSite(t *testing.T) {
	url, _ := ParseRequestURI("https://monzo.com")
	type args struct {
		entryPage *Url
	}
	tests := []struct {
		name string
		args args
		want *Site
	}{
		{"validSite", args{url}, &Site{
			Url:      url,
			PageTree: NewPage(url),
			HashMap:  make(map[string][]string),
			Mu:       &sync.RWMutex{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSite(tt.args.entryPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSite_AddPageToSite(t *testing.T) {
	site := getTestSite()

	validUrl, _ := ParseRequestURI("https://monzo.com/news")
	validPage := NewPage(validUrl)

	existUrl, _ := ParseRequestURI("https://monzo.com/news")
	existPage := NewPage(existUrl)

	type args struct {
		page *Page
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{validPage}, false},
		{"unsuccess", args{existPage}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := site.AddPageToSite(tt.args.page); (err != nil) != tt.wantErr {
				t.Errorf("Site.AddPageToSite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inMap(t *testing.T) {
	site := getTestSite()

	type args struct {
		s string
		m map[string][]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"true", args{"https://monzo.com/blog/haha", site.HashMap}, true},
		{"false", args{"https://monzo.com/blog/bebe", site.HashMap}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inMap(tt.args.s, tt.args.m); got != tt.want {
				t.Errorf("inMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestSite() *Site {
	url, _ := ParseRequestURI("https://monzo.com")
	subUrl, _ := url.ParseUrl("https://monzo.com/blog")
	subSubUrl, _ := url.ParseUrl("https://monzo.com/blog/haha")
	site := NewSite(url)
	site.HashMap[subUrl.String()] = []string{subSubUrl.String()}
	site.HashMap[subSubUrl.String()] = []string{}
	return site
}
