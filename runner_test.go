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
			tags: "foo bar,bar foo",
			want: `{"tags":["foo bar","bar foo"]}`,
		},
		{
			tags: "foo, bar     ",
			want: `{"tags":["foo","bar"]}`,
		},
		{
			tags: "    foo,bar",
			want: `{"tags":["foo","bar"]}`,
		},
		{
			testIDs: "200,300",
			tags:    "foo,bar",
			want:    `{"tests":[200,300]}`,
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
		smartFolderID = test.smartFolderID
		siteID = test.siteID
		tags = test.tags
		testIDs = test.testIDs
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
		{
			testIDs:       "200,300",
			smartFolderID: 200,
			want:          `{"tests":[200,300]}`,
		},
	}
	for _, test := range testCases {
		checkBody := func(body []byte) {
			if string(body) != test.want {
				t.Errorf(`expected %v, got %v`, test.want, string(body))
			}
		}
		ts := newTestPostServer(checkBody)
		baseURL = ts.URL
		smartFolderID = test.smartFolderID
		siteID = test.siteID
		tags = test.tags
		testIDs = test.testIDs
		createRun()
		ts.Close()
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
		{
			testIDs: "200,300",
			siteID:  200,
			want:    `{"tests":[200,300]}`,
		},
	}
	for _, test := range testCases {
		checkBody := func(body []byte) {
			if string(body) != test.want {
				t.Errorf(`expected %v, got %v`, test.want, string(body))
			}
		}
		ts := newTestPostServer(checkBody)
		baseURL = ts.URL
		smartFolderID = test.smartFolderID
		siteID = test.siteID
		tags = test.tags
		testIDs = test.testIDs
		createRun()
		ts.Close()

	}
}

func TestRunByTestID(t *testing.T) {
	testCases := []testRequest{
		{
			testIDs: "200",
			want:    `{"tests":[200]}`,
		},
		{
			testIDs: "",
			want:    `{}`,
		},
		{
			testIDs: "200,300",
			siteID:  200,
			want:    `{"tests":[200,300]}`,
		},
		{
			testIDs: "	200,300,400	",
			siteID: 200,
			want:   `{"tests":[200,300,400]}`,
		},
		{
			testIDs: " 200, 300, 400",
			siteID:  200,
			want:    `{"tests":[200,300,400]}`,
		},
		{
			smartFolderID: 200,
			siteID:        300,
			tags:          "foo, bar",
			testIDs:       "300,500,800",
			want:          `{"tests":[300,500,800]}`,
		},
	}

	for _, test := range testCases {
		checkBody := func(body []byte) {
			if string(body) != test.want {
				t.Errorf(`expected %v, got %v`, test.want, string(body))
			}
		}
		ts := newTestPostServer(checkBody)
		baseURL = ts.URL
		smartFolderID = test.smartFolderID
		siteID = test.siteID
		tags = test.tags
		testIDs = test.testIDs
		createRun()
		ts.Close()
	}
}
