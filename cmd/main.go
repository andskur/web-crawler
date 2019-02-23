package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/andskur/web-crawler"
	"github.com/sirupsen/logrus"
)

// usage constant provide help message
const usage = "Usage:\n    {url} -{flags} \nExample: ./web-crawler https://monzo.com\n"

var errNoTarget = errors.New("no target provided")

func main() {
	target := getTarget()

	flagSet := flag.NewFlagSet("set", flag.ExitOnError)
	fn := flagSet.String("filename", "", "Filename to write output. Example -output sitemap")
	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		logrus.Fatal(err)
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

	var fileName string
	switch *fn {
	case "":
		fileName = craw.Site.EntryPage.Url.Host + ".json"
	default:
		fileName = *fn + ".json"
	}

	if err := ioutil.WriteFile(fileName, jsonFormat, 0644); err != nil {
		logrus.Fatal(err)
	}
}

// getTarget parse target URL from command lines argument
func getTarget() (target string) {
	if len(os.Args) < 2 {
		fmt.Print(usage)
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
