package main

import (
	"fmt"
	"regexp"
)

// ContentRegExp represents a function type that takes a string (ideally a domain) and returns a string (a regular expression)
type ContentRegExp func(string) string

// Links given a specific domain returns a regular expression that can match the pattern `<a href="[http/https]domain" />`
// or the relative path version
func Links(domain string) string {
	return fmt.Sprintf("(a (.*))href*=*\"((((https|http):\\/\\/%s)|(\\/)).*?)\"", domain)
}

// Assets given a specific domain returns a regular expression that can match the pattern `<link href="[http/https]domain" />`
// or the relative path version
func Assets(domain string) string {
	return fmt.Sprintf("(link (.*))href*=*\"((((https|http):\\/\\/%s)|(\\/)).*?)\"", domain)
}

// Images given a specific domain returns a regular expression that can match the pattern `<img src="[http/https]domain" />`
// or the relative path version
func Images(domain string) string {
	return fmt.Sprintf("(img (.*))src*=*\"((((https|http):\\/\\/%s)|(\\/)).*?)\"", domain)
}

// Parse returns an array of strings representing resources urls found in an `htmlPage` and matched against predefined regular
// expressions defined as ContentRegExp
func Parse(cre ContentRegExp, domain string, htmlPage []byte) ([]string, error) {
	urlMap := make(map[string]interface{})

	r, err := regexp.Compile(cre(domain))
	if err != nil {
		return nil, err
	}

	p := string(htmlPage)
	matches := r.FindAllStringSubmatch(p, -1)
	for _, m := range matches {
		urlMap[m[3]] = nil
	}

	var urls = make([]string, 0)
	for link := range urlMap {
		urls = append(urls, link)
	}

	return urls, nil
}
