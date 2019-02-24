package application

import (
	"errors"
	"fmt"
	"os"

	"github.com/andskur/web-crawler/application/writer"

	"github.com/andskur/web-crawler/application/crawler"
)

var errInvalidMapType = errors.New("Invalid sitemap type\nSupported types:\n\t hash - Hash Map\n\t tree - page tree")

// Application represent Crawler Application structure
type Application struct {
	*crawler.Crawler
	writer.IWriter
	*Config
	Output interface{}
}

// Config represent Crawler Application config
type Config struct {
	Filename     string
	MapType      string
	OutputFormat writer.Format
}

func NewConfig(mapType, outputFormat string) (*Config, error) {
	format, err := writer.ParseFormats(outputFormat)
	if err != nil {
		return nil, err
	}
	return &Config{MapType: mapType, OutputFormat: format}, nil
}

// NewApplication create new Web Crawler Application instance with
// from given configuration parameters
func NewApplication(target, fileName, mapType, outputFormat string) (*Application, error) {

	cfg, err := NewConfig(mapType, outputFormat)
	if err != nil {
		return nil, err
	}

	app := &Application{Config: cfg}

	if err := app.initWriter(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := app.initCrawler(target); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

// initWriter initialize application Crawler instance
func (a *Application) initCrawler(target string) error {
	craw, err := crawler.NewCrawler(target)
	if err != nil {
		return err
	}
	a.Crawler = craw
	return nil
}

// initWriter initialize application Output Writer instance
func (a *Application) initWriter() error {
	wrt, err := writer.NewWriter(a.OutputFormat)
	if err != nil {
		return err
	}
	a.IWriter = wrt
	return nil
}

// formatFilename format filename to correct value
func (c *Config) formatFilename(name string) {
	c.Filename = fmt.Sprintf("%s.%s", name, c.OutputFormat)
}
