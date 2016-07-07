package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRequest struct {
	smartFolderID int
	siteID        int
	tags          string
	testIDs       string
	want          string
	got           string
}

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
	testCases := []testRequest{
		{
			tags: "foo,bar",
			want: `{"tags":["foo","bar"]}`,
		},
		{
			tags: "foo",
			want: `{"tags":["foo"]}`,
		},
		{
			tags: "foo,bar     ",
			want: `{"tags":["foo","bar"]}`,
		},
		{
			tags: "    foo,bar",
			want: `{"tags":["foo","bar"]}`,
		},
	}
	for _, test := range testCases {
		checkBody := func(body []byte) {
			if string(body) != test.want {
				t.Errorf(`expected %v got %v`, test.want, string(body))
			}
		}
		ts := newTestPostServer(checkBody)
		baseURL = ts.URL
		smartFolderID = 0
		siteID = 0
		tags = test.tags
		createRun()
		ts.Close()
	}
}

func TestRunBySmartFolder(t *testing.T) {
	testCases := []testRequest{
		{
			smartFolderID: 0,
			want:          `{}`,
		},
		{
			smartFolderID: 200,
			want:          `{"smart_folder_id":200}`,
		},
	}
	for _, test := range testCases {
		checkBody := func(body []byte) {
			if string(body) != test.want {
				t.Errorf(`expected %v, got %v`, test.want, string(body))
			}
		}
		ts := newTestPostServer(checkBody)
		defer ts.Close()
		baseURL = ts.URL
		smartFolderID = test.smartFolderID
		siteID = 0
		tags = ""
		createRun()
	}
}

func TestRunBySiteId(t *testing.T) {
	testCases := []testRequest{
		{
			siteID: 0,
			want:   `{}`,
		},
		{
			siteID: 200,
			want:   `{"site_id":200}`,
		},
	}
	for _, test := range testCases {
		checkBody := func(body []byte) {
			if string(body) != test.want {
				t.Errorf(`expected %v, got %v`, test.want, string(body))
			}
		}
		ts := newTestPostServer(checkBody)
		defer ts.Close()
		baseURL = ts.URL
		smartFolderID = 0
		siteID = test.siteID
		tags = ""
		createRun()
	}
}

func TestRunByTestID(t *testing.T) {
	checkBody := func(body []byte) {
		if string(body) != `{"tests":"1,3,4,7"}` {
			t.Errorf(`expected {"tests":"1,3,4,7"}, got %v`, string(body))
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
