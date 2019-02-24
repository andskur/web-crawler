package application

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/andskur/web-crawler/application/crawler"
	"github.com/andskur/web-crawler/application/writer"
	"github.com/andskur/web-crawler/config"
)

var errInvalidMapType = errors.New("Invalid sitemap type\nSupported types:\n\t hash - Hash Map\n\t tree - page tree")

// Application represent Crawler Application structure
type Application struct {
	*config.Config
	*crawler.Crawler
	Writer writer.IWriter
	Output interface{}
}

// NewApplication create new Web Crawler Application instance with
// from given configuration parameters
func NewApplication(target, fileName, mapType, outputFormat string) (*Application, error) {
	cfg, err := config.NewConfig(mapType, outputFormat)
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

	logrus.SetFormatter(&logrus.TextFormatter{
		// DisableColors: true,
		FullTimestamp: true,
	})

	return app, nil
}

// setFilename set output file correct name with extension
func (a *Application) setFilename(fileName string) {
	switch fileName {
	case "":
		a.FormatFilename(a.Site.Url.Host)
	default:
		a.FormatFilename(fileName)
	}
}

// setOutput set valid application output type
func (a *Application) setOutput() error {
	switch a.MapType {
	case "hash":
		a.Output = a.Site.HashMap
	case "tree":
		a.Output = a.Site
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
	a.Writer = wrt
	return nil
}
