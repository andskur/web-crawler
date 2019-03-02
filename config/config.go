package config

import (
	"fmt"
	"runtime"

	"github.com/andskur/web-crawler/application/site"
	"github.com/andskur/web-crawler/application/writer"
)

// Config represent Crawler Application config
type Config struct {
	Target    *site.Url     // target web site page
	Filename  string        // name of file for output write
	MapType   string        // type od sitemap, Page tree or Hash map
	Output    writer.Format // output format, Json or Xml
	Semaphore chan int      // semaphore for Parallelization restriction
	Verbose   bool          // verbose mode
}

// NewConfig create new config instance from given parameters
func NewConfig(target, fileName, mapType, outputFormat string, verbose, parallelizm bool) (*Config, error) {
	cfg := &Config{
		MapType: mapType,
		Verbose: verbose,
	}

	// set target url
	if err := cfg.setTarget(target); err != nil {
		return nil, err
	}

	// set output format
	if err := cfg.setOutput(outputFormat); err != nil {
		return nil, err
	}

	// set file name
	cfg.setFileName(fileName)

	// set Semaphore channel
	cfg.setSemaphore(parallelizm)

	return cfg, nil
}

// setOutput create target Web Site Url from string
// and set it to current Config instance
func (c *Config) setTarget(target string) (err error) {
	c.Target, err = site.ParseRequestURI(target)
	return
}

// setTarget parse output format from given string
// and set it to current Config instance
func (c *Config) setOutput(output string) (err error) {
	c.Output, err = writer.ParseFormats(output)
	return
}

// setTarget set filename to current Config instance
func (c *Config) setFileName(fileName string) {
	switch fileName {
	case "":
		c.Filename = formatFilename(c.Target.Host, c.Output)
	default:
		c.Filename = formatFilename(fileName, c.Output)
	}
}

// setTarget set filename to current Config instance
func (c *Config) setSemaphore(parralelizm bool) {
	switch {
	case parralelizm:
		c.Semaphore = initCapacity(runtime.NumCPU())
	default:
		c.Semaphore = initCapacity(10000)
	}
}

// formatFilename format filename to correct value
func formatFilename(name string, extension writer.Format) string {
	return fmt.Sprintf("%s.%s", name, extension)
}

// fills up a channel of integers to Semaphore capacity
func initCapacity(maxOutstanding int) chan int {
	ch := make(chan int, maxOutstanding)
	for i := 0; i < maxOutstanding; i++ {
		ch <- 1
	}
	return ch
}
