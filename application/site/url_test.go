package site

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestUrl_ParseUrl(t *testing.T) {

	validUrl, err := url.Parse("https://monzo.com")
	if err != nil {
		return
	}

	validSubUrl, err := validUrl.Parse("/blog/monzo")
	if err != nil {
		return
	}

	type fields struct {
		URL *url.URL
	}
	type args struct {
		ref string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Url
		wantErr bool
	}{
		{"validUrl", fields{nil}, args{"https://monzo.com"}, &Url{validUrl}, false},
		{"validSubUrl", fields{validUrl}, args{"/blog/monzo"}, &Url{validSubUrl}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Url{
				URL: tt.fields.URL,
			}
			got, err := u.ParseUrl(tt.args.ref)
			fmt.Println(got.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("Url.ParseUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Url.ParseUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRequestURI(t *testing.T) {

	validUrl, err := url.Parse("https://monzo.com")
	if err != nil {
		return
	}

	type args struct {
		rawUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    *Url
		wantErr bool
	}{
		{"validUrl", args{"https://monzo.com"}, &Url{validUrl}, false},
		{"invalidUrl", args{"blog/test"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRequestURI(tt.args.rawUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequestURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRequestURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
