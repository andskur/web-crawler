package application

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/andskur/web-crawler/application/crawler"
	"github.com/andskur/web-crawler/application/writer"
	"github.com/andskur/web-crawler/config"
)

var errInvalidMapType = errors.New("Invalid sitemap type\nSupported types:\n\t hash - Hash Map\n\t tree - page tree")

// Application represent Crawler Application structure
type Application struct {
	*config.Config                  // configuration params
	*crawler.Crawler                // web crawler instance
	Writer           writer.IWriter // output writer instance
}

// NewApplication create new Web Crawler Application instance with
// from given configuration parameters
func NewApplication(cfg *config.Config) (*Application, error) {
	app := &Application{Config: cfg}
	if err := app.initApp(); err != nil {
		return nil, err
	}

	return app, nil
}

// initApp initialize all necessary Application instances
func (a *Application) initApp() error {
	// init Writer
	if err := a.initWriter(); err != nil {
		return err
	}

	// init Crawler
	if err := a.initCrawler(); err != nil {
		return err
	}

	// init Logger
	a.initLogger()

	return nil
}

// initWriter initialize Application Crawler instance
func (a *Application) initCrawler() (err error) {
	a.Crawler, err = crawler.NewCrawler(a.Target, a.Config.Verbose, a.Config.Semaphore)
	return
}

// initWriter initialize Application Output Writer instance
func (a *Application) initWriter() (err error) {
	a.Writer, err = writer.NewWriter(a.Output)
	return
}

// initLogger initialize Application logger formatter
func (a *Application) initLogger() {
	formatter := &logrus.TextFormatter{}
	switch {
	case a.Config.Verbose:
		formatter.FullTimestamp = true
	default:
		formatter.FullTimestamp = false
		formatter.DisableColors = true
	}
	logrus.SetFormatter(formatter)
}

// WriteOutput write Application output to file
func (a *Application) WriteOutput() error {
	if err := a.formatOutput(); err != nil {
		return err
	}
	if err := a.Writer.WriteTo(a.Site, a.Filename); err != nil {
		return err
	}

	fmt.Printf("%s sitemap written to %s\n", strings.Title(a.MapType), a.Filename)

	return nil
}

// TODO don't like output methods... Best way - interface and enum, like IWriter

// FormatOutput format application output after execution
func (a *Application) formatOutput() error {
	switch a.MapType {
	case "hash":
		a.Site.PageTree = nil
	case "tree":
		a.Site.HashMap = nil
	default:
		return errInvalidMapType
	}
	return nil
}
