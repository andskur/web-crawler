package writer

import (
	"fmt"
	"testing"
)

func TestFormat_String(t *testing.T) {
	tests := []struct {
		name string
		w    Format
		want string
	}{
		{"getJson", JSON, "json"},
		{"getXml", XML, "xml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.String(); got != tt.want {
				fmt.Println(got)
				t.Errorf("Format.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFormats(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Format
		wantErr bool
	}{
		{"getJson", args{"json"}, JSON, false},
		{"getXml", args{"xml"}, XML, false},
		{"invalid", args{"invalid"}, unsupported, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFormats(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFormats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFormats() = %v, want %v", got, tt.want)
			}
		})
	}
}
