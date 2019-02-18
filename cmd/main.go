package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/andskur/web-crawler"
	"github.com/sirupsen/logrus"
)

var target = "https://monzo.com/"

func main() {
	crawler, err := crawler.NewCrawler(target)
	if err != nil {
		logrus.Fatal(err)
	}

	crawler.StartCrawling()

	logrus.Info(crawler.Site.TotalPages)
	logrus.Info(crawler.TotalDelay)

	jsonFormat, err := json.MarshalIndent(crawler.HashMap, "", " ")
	if err != nil {
		logrus.Fatal(err)
	}

	if err := ioutil.WriteFile("test.json", jsonFormat, 0644); err != nil {
		logrus.Fatal(err)
	}
}
