package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestPostServer(expectedBody string, resp string, statusCode int, t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		flagParams := string(body)
		if flagParams != expectedBody {
			t.Errorf("fetchRource hit wrong endpoint (wanted %v but got %v)", expectedBody, flagParams)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(flagParams))
	}))
}

func TestRunByTags(t *testing.T) {
	expectedBody := `{"tags":["foo","bar"]}`
	sitesResp := "Post Request Successful"
	ts := newTestPostServer(expectedBody, sitesResp, 200, t)
	defer ts.Close()
	baseURL = ts.URL
	smartFolderID = 0
	siteID = 0
	tags = "foo,bar"
	createRun()
}

func TestRunBySmartFolder(t *testing.T) {
	expectedBody := `{"tags":[""],"smart_folder_id":700}`
	sitesResp := "Post Request Successful"
	ts := newTestPostServer(expectedBody, sitesResp, 200, t)
	defer ts.Close()
	baseURL = ts.URL
	smartFolderID = 700
	siteID = 0
	tags = ""
	createRun()
}

func TestRunBySiteId(t *testing.T) {
	expectedBody := `{"tags":[""],"site_id":800}`
	sitesResp := "Post Request Successful"
	ts := newTestPostServer(expectedBody, sitesResp, 200, t)
	defer ts.Close()
	baseURL = ts.URL
	smartFolderID = 0
	siteID = 800
	tags = ""
	createRun()
}
