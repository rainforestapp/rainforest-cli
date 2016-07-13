package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	tempTestInBackground := runTestInBackground
	runTestInBackground = true
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
	runTestInBackground = tempTestInBackground
}

func TestRunBySmartFolder(t *testing.T) {
	tempTestInBackground := runTestInBackground
	runTestInBackground = true
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
	runTestInBackground = tempTestInBackground
}

func TestRunBySiteId(t *testing.T) {
	tempTestInBackground := runTestInBackground
	runTestInBackground = true
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
	runTestInBackground = tempTestInBackground
}

func TestRunByTestID(t *testing.T) {
	tempTestInBackground := runTestInBackground
	runTestInBackground = true
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
			testIDs: " 200,300 , 400",
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
	runTestInBackground = tempTestInBackground
}

func checkRequest(want string, r *http.Request, t *testing.T) {
	if r.URL.String() != want {
		t.Errorf("Expected %v, got %v", want, r.URL)
	}
}

func TestRunInForeground(t *testing.T) {
	var percent int
	tempTestInBackground := runTestInBackground
	runTestInBackground = false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(200)
		if r.Method == "POST" {
			checkRequest("/runs", r, t)
			w.Write([]byte(`{"id": 7000}`))
		}

		if r.Method == "GET" {
			checkRequest("/runs/7000", r, t)
			percent += 100
			if percent < 100 {
				// response := `{
				// 		"id": 78902,
				// 		"state": "in_progress",
				// 		"state_details": {
				// 			"is_final_state": ` + "false" + `
				// 		},
				// 		"result": ` + "in_progress" + `,
				// 		"current_progress": {"percent": ` + strconv.Itoa(percent) + `}}`

				response := runStatusResponse(`"in_progress"`, "false", percent)
				w.Write([]byte(response))
			} else {
				response := runStatusResponse(`"passed"`, "true", 100)
				w.Write([]byte(response))
			}
		}
	}))
	baseURL = ts.URL
	createRun()
	runTestInBackground = tempTestInBackground
}

func runStatusResponse(result string, finalState string, percent int) string {
	// response := `{
	// 		"id": 78902,
	// 		"state": "in_progress",
	// 		"state_details": {
	// 			"is_final_state": ` + finalState + `
	// 		},
	// 		"result": ` + result + `,
	// 		"current_progress": {"percent": ` + strconv.Itoa(percent) + `}}`

	response := `{
		"id": 78902,
		"state": "testing",
		"state_details": {
			"is_final_state": ` + finalState + `
		},
		"result": ` + result + `,
		"current_progress": {"percent": ` + strconv.Itoa(percent) + `}}`
	return response
}
