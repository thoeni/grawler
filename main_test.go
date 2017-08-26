package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"sync"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
)

type domainTest struct {
	url				string
	expectedBaseURL string
	expectedDomain 	string
}

var testData = []domainTest {
	{
		url: "https://thoeni.io/",
		expectedBaseURL: "https://thoeni.io",
		expectedDomain: "thoeni.io",
	},
	{
		url: "https://thoeni.io/blog",
		expectedBaseURL: "https://thoeni.io",
		expectedDomain: "thoeni.io",
	},
	{
		url: "https://thoeni.io/blog/",
		expectedBaseURL: "https://thoeni.io",
		expectedDomain: "thoeni.io",
	},
	{
		url: "http://thoeni.io/blog/post",
		expectedBaseURL: "http://thoeni.io",
		expectedDomain: "thoeni.io",
	},
	{
		url: "http://thoeni.io/index.html",
		expectedBaseURL: "http://thoeni.io",
		expectedDomain: "thoeni.io",
	},
}

var expectedSitemap map[string]Page
var testJson []byte


func TestMain(m *testing.M) {
	f, err := ioutil.ReadFile("./test-data/thoeni.io.0.json")
	if err != nil {
		fmt.Printf("Cannot open file because of: %v\n", err)
		return
	}

	testJson = f

	expectedSitemap = make(map[string]Page)
	if err := json.Unmarshal(testJson, &expectedSitemap); err != nil {
		fmt.Printf("Cannot unmarshal file because of: %v\n", err)
		return
	}

	m.Run()
}

func TestExtractDomainSuccess(t *testing.T) {

	for _, tt := range testData {
		siteDetails, err := extractDomain(tt.url)

		if err != nil {
			t.Fail()
		}

		actualBaseURL := siteDetails.baseURL
		actualDomain := siteDetails.domain

		assert.Equal(t, tt.expectedBaseURL, actualBaseURL)
		assert.Equal(t, tt.expectedDomain, actualDomain)
	}
}

func TestCallExampleSite(t *testing.T) {
	links = make(chan resource)
	var wgr sync.WaitGroup
	wgr.Add(1)
	baseURL := "https://thoeni.io/"
	domain := "thoeni.io"
	level := 0
	go func() {
		for l := range links {
			assert.Contains(t, expectedSitemap[baseURL].Links, l.Url)
		}
	}()

	go fetch(baseURL, domain, resource{baseURL, 0}, &level, &wgr)

	wgr.Wait()
	close(links)

	assert.Equal(t, baseURL, sitemap[baseURL].Url)
	assert.Equal(t, 23, len(sitemap[baseURL].Links))
}

func TestSaveToFile(t *testing.T) {
	saveToFile(expectedSitemap, "TestFlushToFile.json", false)

	actualFile, err := ioutil.ReadFile("TestFlushToFile.json")
	if err != nil {
		fmt.Printf("Cannot open file because of: %v\n", err)
		t.Fail()
		return
	}

	assert.Equal(t, testJson, actualFile)

	os.Remove("TestFlushToFile.json")
}