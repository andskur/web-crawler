package main

import (
	"encoding/json"
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
	mt := flagSet.String("mt", "hash", "-mt {hash || tree} sitemap type, hash map or nested tree (default \"hash\")")

	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		logrus.Fatal(err)
	}

	app, err := application.NewApplication(target, *fn, *mt)
	if err != nil {
		logrus.Fatal(err)
	}

	app.StartCrawling()

	logrus.Info(app.Site.TotalPages)
	logrus.Info(app.TotalDelay)

	jsonFormat, err := json.MarshalIndent(app.Output, "", " ")
	if err != nil {
		logrus.Fatal(err)
	}

	if err := ioutil.WriteFile(app.Filename, jsonFormat, 0644); err != nil {
		logrus.Fatal(err)
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
