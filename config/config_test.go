package config

import (
	"reflect"
	"testing"

	"github.com/andskur/web-crawler/application/writer"

	"github.com/andskur/web-crawler/application/site"
)

/*var testConfig = Config{
	Target: getValidUrl(),
	Filename: "sitemap",
	MapType ""
}*/

func TestNewConfig(t *testing.T) {
	type args struct {
		target       string
		fileName     string
		mapType      string
		outputFormat string
		verbose      bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "validHashJson",
			args: args{
				target:       "https://monzo.com",
				fileName:     "",
				mapType:      "hash",
				outputFormat: "json",
				verbose:      false,
			},
			want: &Config{
				Target:   getValidUrl(),
				Filename: "monzo.com.json",
				MapType:  "hash",
				Output:   writer.JSON,
				Verbose:  false,
			},
			wantErr: false,
		},
		{
			name: "validHashXml",
			args: args{
				target:       "https://monzo.com",
				fileName:     "",
				mapType:      "hash",
				outputFormat: "xml",
				verbose:      false,
			},
			want: &Config{
				Target:   getValidUrl(),
				Filename: "monzo.com.xml",
				MapType:  "hash",
				Output:   writer.XML,
				Verbose:  false,
			},
			wantErr: false,
		},
		{
			name: "validTreeXml",
			args: args{
				target:       "https://monzo.com",
				fileName:     "",
				mapType:      "tree",
				outputFormat: "xml",
				verbose:      false,
			},
			want: &Config{
				Target:   getValidUrl(),
				Filename: "monzo.com.xml",
				MapType:  "tree",
				Output:   writer.XML,
				Verbose:  false,
			},
			wantErr: false,
		},
		{
			name: "validTreeJson",
			args: args{
				target:       "https://monzo.com",
				fileName:     "",
				mapType:      "tree",
				outputFormat: "json",
				verbose:      false,
			},
			want: &Config{
				Target:   getValidUrl(),
				Filename: "monzo.com.json",
				MapType:  "tree",
				Output:   writer.JSON,
				Verbose:  false,
			},
			wantErr: false,
		},
		{
			name: "validVerbose",
			args: args{
				target:       "https://monzo.com",
				fileName:     "",
				mapType:      "tree",
				outputFormat: "xml",
				verbose:      true,
			},
			want: &Config{
				Target:   getValidUrl(),
				Filename: "monzo.com.xml",
				MapType:  "tree",
				Output:   writer.XML,
				Verbose:  true,
			},
			wantErr: false,
		},
		{
			name: "noStandartFileName",
			args: args{
				target:       "https://monzo.com",
				fileName:     "sitemap",
				mapType:      "hash",
				outputFormat: "json",
				verbose:      false,
			},
			want: &Config{
				Target:   getValidUrl(),
				Filename: "sitemap.json",
				MapType:  "hash",
				Output:   writer.JSON,
				Verbose:  false,
			},
			wantErr: false,
		},
		{
			name: "invalidFormat",
			args: args{
				target:       "https://monzo.com",
				fileName:     "",
				mapType:      "hash",
				outputFormat: "graph",
				verbose:      false,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.target, tt.args.fileName, tt.args.mapType, tt.args.outputFormat, tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getValidUrl() (validUrl *site.Url) {
	validUrl, _ = site.ParseRequestURI("https://monzo.com")
	return
}
