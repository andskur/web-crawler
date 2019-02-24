package main

import (
	"errors"
	"flag"
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

	logrus.Info(app.Site.TotalPages)
	logrus.Info(app.Duration)

	if err := app.WriteOutput(); err != nil {
		logrus.Fatal(err)
	}
}

// getTarget parse target URL from command lines argument
func getTarget() (target string) {
	if len(os.Args) < 2 {
		logrus.Fatal(usage)
		os.Exit(1)
	}
	target = os.Args[1]
	if target == "" {
		logrus.Fatal(errNoTarget)
	}
	return
}
