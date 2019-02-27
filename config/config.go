package config

import (
	"fmt"

	"github.com/andskur/web-crawler/application/site"
	"github.com/andskur/web-crawler/application/writer"
)

// Config represent Crawler Application config
type Config struct {
	Target   *site.Url     // target web site page
	Filename string        // name of file for output write
	MapType  string        // type od sitemap, Page tree or Hash map
	Output   writer.Format // output format, Json or Xml
	Verbose  bool          // verbose mode
}

// NewConfig create new config instance from given parameters
func NewConfig(target, fileName, mapType, outputFormat string, verbose bool) (*Config, error) {
	cfg := &Config{
		MapType: mapType,
		Verbose: verbose,
	}

	// set target url
	if err := cfg.setTarget(target); err != nil {
		return nil, err
	}

	// set output format
	if err := cfg.setOutput(outputFormat); err != nil {
		return nil, err
	}

	// set file name
	cfg.setFileName(fileName)

	return cfg, nil
}

// setOutput create target Web Site Url from string
// and set it to current Config instance
func (c *Config) setTarget(target string) (err error) {
	c.Target, err = site.ParseRequestURI(target)
	return
}

// setTarget parse output format from given string
// and set it to current Config instance
func (c *Config) setOutput(output string) (err error) {
	c.Output, err = writer.ParseFormats(output)
	return
}

// setTarget set filename to current Config instance
func (c *Config) setFileName(fileName string) {
	switch fileName {
	case "":
		c.formatFilename(c.Target.Host)
	default:
		c.formatFilename(fileName)
	}
}

// formatFilename format filename to correct value
func (c *Config) formatFilename(name string) {
	c.Filename = fmt.Sprintf("%s.%s", name, c.Output)
}
