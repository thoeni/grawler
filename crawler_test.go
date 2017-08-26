package main

import (
	"testing"
	"io/ioutil"
	"net/http/httptest"
	"net/http"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestGetPage(t *testing.T) {
	localResp, _ := ioutil.ReadFile("test-data/thoeni.io.html")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, string(localResp))
	}))
	defer ts.Close()

	p, err := GetPage(ts.URL, "thoeni.io")
	if err != nil {
		t.Error("Failed to get page")
	}

	assert.Equal(t, ts.URL, p.Url)
	assert.Len(t, p.Links, 23)
	assert.Len(t, p.Assets, 5)
}

func TestGetPage404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "")
	}))
	defer ts.Close()

	_, err := GetPage(ts.URL, "thoeni.io")

	assert.EqualError(t, err, "Page not reachable because of error [404 Not Found]")
}