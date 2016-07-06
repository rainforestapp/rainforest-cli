package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestPostServer(check func([]byte)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		check(body)
		w.WriteHeader(200)
		w.Write([]byte("Success"))
	}))
}

func TestRunByTags(t *testing.T) {
	checkBody := func(body []byte) {
		if string(body) != `{"tags":["foo","bar"]}` {
			t.Errorf(`expected {"tags":["foo","bar"]}, got %v`, string(body))
		}
	}
	ts := newTestPostServer(checkBody)
	defer ts.Close()
	baseURL = ts.URL
	smartFolderID = 0
	siteID = 0
	tags = "foo,bar"
	createRun()
}

func TestRunBySmartFolder(t *testing.T) {
	checkBody := func(body []byte) {
		if string(body) != `{"tags":[""],"smart_folder_id":700}` {
			t.Errorf(`expected {"tags":[""],"smart_folder_id":700}, got %v`, string(body))
		}
	}
	ts := newTestPostServer(checkBody)
	defer ts.Close()
	baseURL = ts.URL
	smartFolderID = 700
	siteID = 0
	tags = ""
	createRun()
}

func TestRunBySiteId(t *testing.T) {
	checkBody := func(body []byte) {
		if string(body) != `{"tags":[""],"site_id":800}` {
			t.Errorf(`expected {"tags":[""],"site_id":800}, got %v`, string(body))
		}
	}
	ts := newTestPostServer(checkBody)
	defer ts.Close()
	baseURL = ts.URL
	smartFolderID = 0
	siteID = 800
	tags = ""
	createRun()
}

func TestRunByTestID(t *testing.T) {
	checkBody := func(body []byte) {
		if string(body) != `{"tests":["1","3","4","7"]}` {
			t.Errorf(`expected {"tests":["1","3","4","7"]}, got %v`, string(body))
		}
	}
	ts := newTestPostServer(checkBody)
	defer ts.Close()
	baseURL = ts.URL
	smartFolderID = 0
	siteID = 0
	testIDs = "1,3,4,7"
	createRun()
}
