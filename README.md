# Web-crawler

**Web-crawler** is simple concurrency web-crawler providing crawling
all internal pages in target Web site. Output will write to Json or Xml
file an can be one og two types - Hash Map or Page Tree.

## Installation
Application uses go-modules for dependencies management,
so you don't need application inside your GOPATH.

#####Build application (in application mani directory):
```
$ go build cmd/web-crawler.go
```

#####Testing:
```
$ go test -v ./...
```

## Usage

##### Example:
```
$ ./web-crawler https://monzo.com
```

##### Options:

```
Usage:
    {url} {-flags}
Example: ./web-crawler https://monzo.com
  -fn string
    	-fn {filename} filename to write output
  -mt string
    	-mt {hash || tree} sitemap type, hash map or page tree (default "hash") (default "hash")
  -of string
    	-of {json || xml} output format, json or xml (default "json") (default "json")
  -p  	-p parralelizm mode
  -v	-v verbose mode
```

