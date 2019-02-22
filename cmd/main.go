package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/andskur/web-crawler"
	"github.com/sirupsen/logrus"
)

// usage constant provide help message
const usage = "Usage:\n    {url}\nExample: ./web-crawler https://monzo.com\n"

var errNoTarget = errors.New("no target provided")

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	target := os.Args[1]

	if target == "" {
		logrus.Error(errNoTarget)
		os.Exit(1)
	}

	craw, err := crawler.NewCrawler(target)
	if err != nil {
		logrus.Fatal(err)
	}

	craw.StartCrawling()

	logrus.Info(craw.Site.TotalPages)
	logrus.Info(craw.TotalDelay)

	jsonFormat, err := json.MarshalIndent(craw.Site, "", " ")
	if err != nil {
		logrus.Fatal(err)
	}

	if err := ioutil.WriteFile("test.json", jsonFormat, 0644); err != nil {
		logrus.Fatal(err)
	}
}
