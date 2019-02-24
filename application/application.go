package application

import (
	"errors"
	"fmt"
	"os"

	"github.com/andskur/web-crawler/application/crawler"
)

var errInvalidMapType = errors.New("Invalid sitemap type\nSupported types:\n\t hash - Hash Map\n\t tree - page tree")

// Application represent Crawler Application structure
type Application struct {
	*crawler.Crawler
	*Config
	Output interface{}
}

// Config represent Crawler Application config
type Config struct {
	Filename     string
	MapType      string
	OutputFormat string
}

// NewApplication create new Web Crawler Application instance with
// from given configuration parameters
func NewApplication(target, fileName, mapType, outputFormat string) (*Application, error) {
	craw, err := crawler.NewCrawler(target)
	if err != nil {
		return nil, err
	}
	cfg := &Config{MapType: mapType, OutputFormat: outputFormat}
	app := &Application{Crawler: craw, Config: cfg}

	app.setFilename(fileName)

	if err := app.setOutput(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return app, nil
}

// setFilename set output file correct name with extension
func (a *Application) setFilename(fileName string) {
	switch fileName {
	case "":
		a.formatFilename(a.SiteTree.EntryPage.Url.Host)
	default:
		a.formatFilename(fileName)
	}
}

// setOutput set valid application output type
func (a *Application) setOutput() error {
	switch a.MapType {
	case "hash":
		a.Output = a.HashMap
	case "tree":
		a.Output = a.SiteTree
	default:
		return errInvalidMapType
	}
	return nil
}

// formatFilename format filename to correct value
func (c *Config) formatFilename(name string) {
	c.Filename = fmt.Sprintf("%s.%s", name, c.OutputFormat)
}
