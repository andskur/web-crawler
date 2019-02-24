package config

import (
	"fmt"

	"github.com/andskur/web-crawler/application/writer"
)

// Config represent Crawler Application config
type Config struct {
	Filename     string
	MapType      string
	OutputFormat writer.Format
}

// NewConfig create new config instance from given parameters
func NewConfig(mapType, outputFormat string) (*Config, error) {
	format, err := writer.ParseFormats(outputFormat)
	if err != nil {
		return nil, err
	}
	return &Config{MapType: mapType, OutputFormat: format}, nil
}

// formatFilename format filename to correct value
func (c *Config) FormatFilename(name string) {
	c.Filename = fmt.Sprintf("%s.%s", name, c.OutputFormat)
}
