package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/andskur/web-crawler/application"
	"github.com/andskur/web-crawler/config"
)

// usage constant provide help message
const usage = "Usage:\n    {url} {-flags} \nExample: ./web-crawler https://monzo.com"

var (
	errNoTarget = errors.New("no target url provided")
)

// TODO - better way move it to separate package or use Cobra-like external tools

func main() {
	// initialize target argument and flags
	var target string
	flagSet := flag.NewFlagSet("set", flag.ExitOnError)
	fn := flagSet.String("fn", "", "-fn {filename} filename to write output")
	mt := flagSet.String("mt", "hash", "-mt {hash || tree} sitemap type, hash map or page tree (default \"hash\")")
	of := flagSet.String("of", "json", "-of {json || xml} output format, json or xml (default \"json\")")
	p := flagSet.Bool("p", false, "-p parralelizm mode")
	v := flagSet.Bool("v", false, "-v verbose mode")

	// validate arguments
	if len(os.Args) < 2 {
		fmt.Println(usage)
		flagSet.PrintDefaults()
		os.Exit(1)
	}

	// get target from command-line argument
	target = os.Args[1]
	if target == "" {
		fmt.Println(errNoTarget)
		fmt.Println(usage)
		flagSet.PrintDefaults()
		os.Exit(1)
	}

	// parse command-line flags
	err := flagSet.Parse(os.Args[2:])
	if err != nil {
		fmt.Println(err)
		fmt.Println(usage)
		flagSet.PrintDefaults()
		os.Exit(1)
	}

	// create new Application Config from cli argument and params
	cfg, err := config.NewConfig(target, *fn, *mt, *of, *v, *p)
	if err != nil {
		logrus.Fatal(err)
	}

	// create new Application with Config
	app, err := application.NewApplication(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	// start Crawling
	if err := app.StartCrawling(); err != nil {
		logrus.Fatal(err)
	}

	// format Crawler output and write it to file
	if err := app.WriteOutput(); err != nil {
		logrus.Fatal(err)
	}

}
