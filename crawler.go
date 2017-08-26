package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"errors"
	"time"
)

// Page represents an HTML page with the resources it depends on and links to
type Page struct {
	Url string
	Links []string
	Assets []string
	Images []string
}

// String method on the Page in order to print the content in a human readable way
func (p Page) String() string {
	return fmt.Sprintf("\n\tUrl: %s\n\t\tLinks: %s\n\t\tAssets: %s\n", p.Url, p.Links, p.Assets)
}

// GetPage returns a new Page structure opportunely populated given a resource and a domain
func GetPage(url string, domain string) (Page, error) {

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	r, err := client.Get(url)

	if err != nil {
		fmt.Errorf("error while getting url: %s", err)
		return Page{}, err
	}

	defer r.Body.Close()

	if (r.StatusCode != http.StatusOK) {
		return Page{}, errors.New(fmt.Sprintf("Page not reachable because of error [%s]", r.Status))
	}

	pBody, _ := ioutil.ReadAll(r.Body)
	var page = Page{
		Url: url,
	}

	links, err := Parse(Links, domain, pBody)
	if err != nil {
		fmt.Errorf("error while parsing links: %s", err)
		return page, err
	}

	page.Links = links

	assets, err := Parse(Assets, domain, pBody)
	if err != nil {
		fmt.Errorf("error while parsing assets: %s", err)
		return page, err
	}

	page.Assets = assets

	images, err := Parse(Images, domain, pBody)
	if err != nil {
		fmt.Errorf("error while parsing images: %s", err)
		return page, err
	}

	page.Images = images

	return page, nil
}
