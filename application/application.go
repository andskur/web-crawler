package application

import (
	"errors"
	"fmt"
	"os"

	"github.com/andskur/web-crawler/application/crawler"
)

var errInvalidMapType = errors.New("Invalid sitemap type\nSupported types:\n\t hash - Hash Map\n\t tree - nested tree")

// Application represent Crawler Application structure
type Application struct {
	*crawler.Crawler
	*Config
	Output interface{}
}

// Config represent Crawler Application config
type Config struct {
	Filename string
	MapType  string
}

// NewApplication create new Web Crawler Application instance with
// from given configuration parameters
func NewApplication(target, fileName, mapType string) (*Application, error) {
	craw, err := crawler.NewCrawler(target)
	if err != nil {
		return nil, err
	}
	cfg := &Config{MapType: mapType}
	app := &Application{Crawler: craw, Config: cfg}
	switch fileName {
	case "":
		app.Filename = app.Site.EntryPage.Url.Host + ".json"
	default:
		app.Filename = fileName + ".json"
	}
	switch app.MapType {
	case "hash":
		app.Output = app.HashMap
	case "tree":
		app.Output = app.Site
	default:
		fmt.Println(errInvalidMapType)
		os.Exit(1)
	}
	return app, nil
}
