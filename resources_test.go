package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

func newTestServer(path, resp string, statusCode int, t *testing.T) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			t.Errorf("fetchRource hit wrong endpoint (wanted %v but got %v)", path, r.URL.Path)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))
	}))
	return ts
}

func TestPrintSites(t *testing.T) {
	sitesResp := `[{"id": 1337, "name": "Dyer"}]`
	ts := newTestServer("/sites.json", sitesResp, 200, t)
	defer ts.Close()
	baseURL = ts.URL

	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()

	printSites()
	pattern := `\| +SITE ID +\| +SITE DESCRIPTION +\|`
	matched, err := regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}

	pattern = `\| +1337 +\| +Dyer +\|`
	matched, err = regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}

	if !matched {
		t.Logf("Table didn't match properly:")
		t.Logf("%v\n", out)
		t.Errorf("should have matched %v", pattern)
	}
}

func TestPrintSitesApiError(t *testing.T) {
	sitesResp := `{"error": "This is a bad thing"}`
	ts := newTestServer("/sites.json", sitesResp, 200, t)
	defer ts.Close()
	baseURL = ts.URL
	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()
	if os.Getenv("BE_CRASHER") == "1" {
		printSites()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestPrintSitesApiError")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want status 1", err)

	printSites()
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

	printFolders()

	pattern := `\| +FOLDER ID +\| +FOLDER DESCRIPTION +\|`
	matched, err := regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}

	pattern = `\| +707 +\| +Foo +\|`
	matched, err = regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
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

	printBrowsers()

	pattern := `\| +BROWSER ID +\| +BROWSER DESCRIPTION +\|`
	matched, err := regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}

	pattern = `\| +firefox +\| +Mozilla Firefox +\|`
	matched, err = regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}

	if !matched {
		t.Logf("Table didn't match properly:")
		t.Logf("%v\n", out)
		t.Errorf("should have matched %v", pattern)
	}
}
