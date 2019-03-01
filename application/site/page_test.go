package site

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewPage(t *testing.T) {
	validUrl, err := ParseRequestURI("https://monzo.com")
	if err != nil {
		return
	}

	type args struct {
		url *Url
	}
	tests := []struct {
		name string
		args args
		want *Page
	}{
		{"validPage", args{validUrl}, &Page{Url: validUrl, Logger: logrus.WithField("page", validUrl.String())}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPage(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPage_inPage(t *testing.T) {
	page := getTestPage()

	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Found", args{"https://monzo.com/blog"}, true},
		{"Found2", args{"https://monzo.com/blog/haha"}, true},
		{"UnFound", args{"https://monzo.com/blog/lalala"}, false},
		{"ParentPage", args{"https://monzo.com"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := page.inPage(tt.args.s); got != tt.want {
				t.Errorf("Page.inPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPage_validateLink(t *testing.T) {
	page := getTestPage()
	type args struct {
		link string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"validPath", args{"/about/site"}, false},
		{"validUrl", args{"https://monzo.com/faq"}, false},
		{"invalidPath", args{"about/contact"}, true},
		{"externalUrl", args{"https://twitter.com/lala"}, true},
		{"email-protected", args{"https://monzo.com/faq/1/email-protection"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := page.validateLink(tt.args.link); (err != nil) != tt.wantErr {
				t.Errorf("Page.validateLink() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPage_validateUrl(t *testing.T) {
	page := getTestPage()

	validUrl, _ := ParseRequestURI("https://monzo.com/news")
	withQuery, _ := validUrl.ParseUrl("/banking/thebest?utm=blabla")
	external, _ := ParseRequestURI("https://facebook.com/https://monzo.com/asd")
	invalidScheme, _ := ParseRequestURI("http://monzo.com/news")
	alreadyInSlice, _ := ParseRequestURI("http://monzo.com/blog")
	alreadyInSlice2, _ := ParseRequestURI("http://monzo.com/blog/haha")

	type args struct {
		url *Url
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{validUrl}, false},
		{"withQuery", args{withQuery}, true},
		{"external", args{external}, true},
		{"invalidScheme", args{invalidScheme}, true},
		{"alreadyInSlice", args{alreadyInSlice}, true},
		{"alreadyInSlice2", args{alreadyInSlice2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := page.validateUrl(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Page.validateUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPage_AddSubPage(t *testing.T) {
	page := getTestPage()

	validLink := "https://monzo.com/news"
	validUrl, _ := ParseRequestURI(validLink)
	validPage := NewPage(validUrl)

	type args struct {
		link string
	}
	tests := []struct {
		name    string
		args    args
		want    *Page
		wantErr bool
	}{
		{"success", args{validLink}, validPage, false},
		{"unsuccess", args{"https://monzo.com/blog"}, nil, true},
		{"external", args{"https://facebook.com"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := page.AddSubPage(tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("Page.AddSubPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Page.AddSubPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestPage() (page *Page) {
	url, _ := ParseRequestURI("https://monzo.com")
	subUrl, _ := url.ParseUrl("/blog")
	subUrl2, _ := url.ParseUrl("/blog/haha")

	page = NewPage(url)
	page.Links = append(page.Links, NewPage(subUrl), NewPage(subUrl2))
	return
}
