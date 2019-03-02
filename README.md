# Web-crawler

**Web-crawler** is simple concurrency web-crawler providing crawling
all internal pages in target Web site. Output will write to Json or Xml
file an can be one og two types - Hash Map or Page Tree.

## Installation
Application uses go-modules for dependencies management,
so you don't need application inside your GOPATH.

### Dependencies
* [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) - verbose logging organisation
* [golang.org/x/net](https://godoc.org/golang.org/x/net/html) - html parsing

#####Build application (in application mani directory):
```bash
$ go build cmd/web-crawler.go
```

#####Testing:
```bash
$ go test -v ./...
```

## Usage

#### Example:
```bash
$ ./web-crawler https://monzo.com
Start crawling web site monzo.com...
Total pages: 718...
All done!
718 pages crawled at monzo.com in 12.963644303s
Hash sitemap written to monzo.com.json
```

#### Options:

```bash
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

#### Flags explanation:

##### **-fn**
Filename of file where sitemap will be written
##### **-mt** 
Sitemap type, can be **hash** Hash Map or **tree** Page Tree

*XML Hash Map example:*
```xml
<site>
 <url>https://monzo.com</url>
 <total_pages>718</total_pages>
 <map>
  <page>
   <url>https://monzo.com</url>
   <total_links>23</total_links>
   <links>
    <url>https://monzo.com/</url>
    <url>https://monzo.com/about</url>
    <url>https://monzo.com/blog</url>
    ...
   </links>
  </page>
  <page>
    <url>https://monzo.com/faq</url>
    <total_links>18</total_links>
    <links>
     <url>https://monzo.com/</url>
     <url>https://monzo.com/about</url>
     <url>https://monzo.com/blog</url>
     ...
    </links>
  </page>
  ...
 </map>
</site>
```

*JSON Page Tree example:*
```json
{
 "url": "https://monzo.com",
 "total_pages": 718,
 "tree": {
  "url": "https://monzo.com",
  "total": 23,
  "links": [
    {
      "url": "https://monzo.com/community",
      "total": 19,
      "links": [
       {
        "url": "https://monzo.com/"
       },
       {
        "url": "https://monzo.com/about"
       },
       {
        "url": "https://monzo.com/blog",
        "total": 21,
        "links": [
           {
             "url": "https://monzo.com/blog/authors/hugo-cornejo/",
             "total": 19,
              "links": [
                {
                  "url": "https://monzo.com/"
                },
                {
                  "url": "https://monzo.com/about"
                },
                {
                  "url": "https://monzo.com/blog"
                },
                ...
              ]
           },
           ...
        ]
       },
       ...
      ]
    },
    {
      "url": "https://monzo.com/about",
      "total": 19,
      "links": [
        {
         "url": "https://monzo.com/"
        },
        {
         "url": "https://monzo.com/about"
        }
        ...
      ]
    }
  ]
}
```

##### **-of** 
Output format, can be **json** or **xml**

##### **-p** 
Parralelizm mode, set goroutins limit, based of available machines CPU'S.
Makes application slower, but resource-safety.

##### **-v** 
Verbose mode
