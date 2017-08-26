package main

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"bytes"
)

func TestGetLinks(t *testing.T) {
	testPage, err := os.OpenFile("test-data/thoeni.io.html", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatalf("Cannot open file because of: %s", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(testPage)

	r, err := Parse(Links, "thoeni.io", buf.Bytes())
	if err != nil {
		t.Fail()
	}

	assert.Len(t, r, 23)
	assert.Contains(t,  r, "https://thoeni.io/page/resources")
}

func TestGetAssets(t *testing.T) {
	testPage, err := os.OpenFile("test-data/thoeni.io.html", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatalf("Cannot open file because of: %s", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(testPage)

	r, err := Parse(Assets, "thoeni.io", buf.Bytes())
	if err != nil {
		t.Fail()
	}

	assert.Len(t, r, 5)
	assert.Contains(t, r, "https://thoeni.io/css/main.css")
}
