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
