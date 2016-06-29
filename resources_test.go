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

func runErrorTest(resource string, t *testing.T) {
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
}

func checkTableCorrect(pattern string, t *testing.T) {
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
	checkTableCorrect(`\| +SITE ID +\| +SITE DESCRIPTION +\|`, t)
	checkTableCorrect(`\| +1337 +\| +Dyer +\|`, t)
}

func TestPrintSitesApiError(t *testing.T) {
	sitesResp := `{"error": "This is a bad thing"}`
	ts := newTestServer("/sites.json", sitesResp, 400, t)
	defer ts.Close()
	baseURL = ts.URL
	runErrorTest("Sites", t)
}

func TestPrintFolders(t *testing.T) {
	sitesResp := `[{"id": 707, "title": "Foo"}]`
	ts := newTestServer("/folders.json", sitesResp, 200, t)
	defer ts.Close()
	baseURL = ts.URL

	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()

	printFolders()
	checkTableCorrect(`\| +FOLDER ID +\| +FOLDER DESCRIPTION +\|`, t)
	checkTableCorrect(`\| +707 +\| +Foo +\|`, t)
}

func TestPrintFoldersApiError(t *testing.T) {
	sitesResp := `{"error": "This is a bad thing"}`
	ts := newTestServer("/folders.json", sitesResp, 600, t)
	defer ts.Close()
	baseURL = ts.URL

	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()
	runErrorTest("Folders", t)
}

func TestPrintBrowsers(t *testing.T) {
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
	checkTableCorrect(`\| +BROWSER ID +\| +BROWSER DESCRIPTION +\|`, t)
	checkTableCorrect(`\| +firefox +\| +Mozilla Firefox +\|`, t)
}

func TestPrintBrowsersApiError(t *testing.T) {
	sitesResp := `{"error": "This is a bad thing"}`
	ts := newTestServer("/clients.json", sitesResp, 600, t)
	defer ts.Close()
	baseURL = ts.URL
	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()
	runErrorTest("Browsrs", t)
}
