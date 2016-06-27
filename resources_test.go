package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
)

func TestPrintSites(t *testing.T) {
	siteResp := `[{"id": 1337, "name": "Dyer"}]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sites.json" {
			t.Errorf("fetchRource hit wrong endpoint (wanted /sites.json but got %v)", r.URL.Path)
		}
		w.Write([]byte(siteResp))
	}))
	defer ts.Close()
	baseURL = ts.URL

	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()

	fetchResource("Sites")

	pattern := `\| +1337 +\| +Dyer +\|`
	matched, err := regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}

	if !matched {
		t.Logf("Table didn't match properly:")
		t.Logf("%v\n", out)
		t.Errorf("should have matched %v", pattern)
	}
}

func TestPrintFolders(t *testing.T) {
	siteResp := `[{"id": 707, "title": "Foo"}]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/folders.json" {
			t.Errorf("fetchRource hit wrong endpoint (wanted /Folders.json but got %v)", r.URL.Path)
		}
		w.Write([]byte(siteResp))
	}))
	defer ts.Close()
	baseURL = ts.URL

	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()

	fetchResource("Folders")

	pattern := `\| +707 +\| +Foo +\|`
	matched, err := regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}

	if !matched {
		t.Logf("Table didn't match properly:")
		t.Logf("%v\n", out)
		t.Errorf("should have matched %v", pattern)
	}
}

func TestBrowsersFolders(t *testing.T) {
	siteResp := `{"available_browsers": [{"name": "firefox", "description": "Mozilla Firefox"}]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/clients.json" {
			t.Errorf("fetchRource hit wrong endpoint (wanted /clients.json but got %v)", r.URL.Path)
		}
		w.Write([]byte(siteResp))
	}))
	defer ts.Close()
	baseURL = ts.URL

	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()

	fetchResource("Browsers")

	pattern := `\| +firefox +\| +Mozilla Firefox +\|`
	matched, err := regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}

	if !matched {
		t.Logf("Table didn't match properly:")
		t.Logf("%v\n", out)
		t.Errorf("should have matched %v", pattern)
	}
}
