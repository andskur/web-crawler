package writer

import (
	"reflect"
	"testing"

	"github.com/andskur/web-crawler/application/writer/xml"

	"github.com/andskur/web-crawler/application/writer/json"
)

func TestNewWriter(t *testing.T) {
	type args struct {
		wtype Format
	}
	tests := []struct {
		name    string
		args    args
		wantWrt IWriter
		wantErr bool
	}{
		{"getJsonWriter", args{JSON}, json.WriterJson{}, false},
		{"getXmlWriter", args{XML}, xml.WriterXml{}, false},
		{"invalidWriter", args{unsupported}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWrt, err := NewWriter(tt.args.wtype)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotWrt, tt.wantWrt) {
				t.Errorf("NewWriter() = %v, want %v", gotWrt, tt.wantWrt)
			}
		})
	}
}
