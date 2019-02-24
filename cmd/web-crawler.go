package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/andskur/web-crawler/application"
)

// usage constant provide help message
const usage = "Usage:\n    {url} {-flags} \nExample: ./web-crawler https://monzo.com"

var (
	errNoTarget = errors.New("no target url provided")
)

func main() {
	target := getTarget()

	flagSet := flag.NewFlagSet("set", flag.ExitOnError)
	fn := flagSet.String("fn", "", "-fn {filename} filename to write output")
	mt := flagSet.String("mt", "hash", "-mt {hash || tree} sitemap type, hash map or page tree (default \"hash\")")
	of := flagSet.String("of", "json", "-of {json || xml} output format, json or xml (default \"json\")")

	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		logrus.Fatal(err)
	}

	app, err := application.NewApplication(target, *fn, *mt, *of)
	if err != nil {
		logrus.Fatal(err)
	}

	app.StartCrawling()

	logrus.Info(app.SiteTree.TotalPages)
	logrus.Info(app.TotalDelay)

	switch app.OutputFormat {
	case "json":
		if err := writeJson(app.Output, app.Filename); err != nil {
			logrus.Fatal(err)
		}
	case "xml":
		if err := writeXml(app.Output, app.Filename); err != nil {
			logrus.Fatal(err)
		}
	}

}

// getTarget parse target URL from command lines argument
func getTarget() (target string) {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		flag.PrintDefaults()
		os.Exit(1)
	}
	target = os.Args[1]
	if target == "" {
		logrus.Error(errNoTarget)
		os.Exit(1)
	}
	return
}

// writeJson write json application data to json file with providing name
func writeJson(data interface{}, fileName string) error {
	jsonFormat, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fileName, jsonFormat, 0644); err != nil {
		return err
	}
	return nil
}

// writeXml write xml application data to xml file with providing name
func writeXml(data interface{}, fileName string) error {
	xmlFormat, err := xml.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fileName, xmlFormat, 0644); err != nil {
		return err
	}
	return nil
}
