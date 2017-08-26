package main

import (
	"fmt"
	"sync"
	"strings"
	"encoding/json"
	"flag"
	"io/ioutil"
	"regexp"
	"sync/atomic"
	"time"
)

var sitemap map[string]Page = make(map[string]Page)

var m sync.Mutex
var links chan resource
var count uint64

func main() {

	url := flag.String("url", "", "URL for website to build sitemap, e.g. https://thoeni.io (required)")
	depth := flag.Int("depth", 0, "Depth level from root (-1 means unbounded)")
	graph := flag.Bool("graph", false, "Whether the generated output should feed the html graph page. Generates a `sitemap.js` file. (default false)")
	output := flag.String("output", "", "Filename for output. Empty value prints to stdout. (default empty)")
	flag.Parse()

	r, _ := regexp.Compile("(https|http):\\/\\/.*\\/*")
	if valid := r.MatchString(*url); !valid {
		fmt.Println("URL not valid. For example -url=https://thoeni.io or -url=https://thoeni.io/post/")
		return
	}

	siteDetails, err := extractDomain(*url)
	if err != nil {
		fmt.Printf("Input parameter %s is invalid", *url, err)
		return
	}

	links = make(chan resource)

	var wg sync.WaitGroup
	wg.Add(1)
	start := time.Now().UTC()

	go func() {
		for l := range links {
			wg.Add(1)
			go fetch(siteDetails.baseURL, siteDetails.domain, l, depth, &wg)
		}
	}()

	links <- resource{Url: *url, Level: 0}

	wg.Done()
	wg.Wait()

	close(links)

	if (*graph == true || *output != "") {
		if err := saveToFile(sitemap, *output, *graph); err != nil {
			fmt.Printf("Error while writing to file, %v", err)
		}
	} else {
		b, _ := json.MarshalIndent(sitemap, "", "   ")
		fmt.Println(string(b))
	}

	fmt.Printf("\nParsed %d pages in %s", count, time.Since(start))
}

type resource struct {
	Url string
	Level int
}

func fetch(baseURL string, domain string, r resource, depth *int, wg *sync.WaitGroup) {
	defer wg.Done()
	if (*depth > -1 && r.Level > *depth) {
		return
	}
	m.Lock()
	if _, visited := sitemap[r.Url]; visited {
		m.Unlock()
		return
	}
	sitemap[r.Url] = Page{Url: r.Url}
	m.Unlock()

	p, err := GetPage(r.Url, domain)
	if err != nil {
		fmt.Printf("error while retrieving the page [%s]: %v", r.Url, err)
	}

	atomic.AddUint64(&count, 1)

	m.Lock()
	sitemap[r.Url] = p
	m.Unlock()

	for _, l := range p.Links {
		// Make relative links absolute
		if strings.HasPrefix(l, "/") {
			l = baseURL +l
		}
		links <- resource{l, (r.Level+1)}
	}
}

type siteDetails struct {
	baseURL string
	domain string
}

func extractDomain(url string) (siteDetails, error) {
	r, err := regexp.Compile("(.*\\/\\/(.*?))(\\/|$)")

	if err != nil {
		return siteDetails{}, err
	}

	m := r.FindStringSubmatch(url)

	return siteDetails{
		baseURL: m[1],
		domain: m[2],
	}, nil
}

func saveToFile(s map[string]Page, filename string, jsVar bool) error {

	b, _ := json.Marshal(s)
	var err error

	if jsVar {
		err = ioutil.WriteFile("sitemap.js", []byte(fmt.Sprintf("var sitemap = %s", b)), 0644)
	}

	if filename != "" {
		err = ioutil.WriteFile(filename, b, 0644)
	}

	return err
}