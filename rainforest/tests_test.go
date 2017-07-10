package rainforest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestGetRFMLIDs(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"

	rfmlIDs := TestIDMappings{
		{ID: 123, RFMLID: "abc"},
		{ID: 456, RFMLID: "xyz"},
	}

	mux.HandleFunc("/tests/rfml_ids", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		enc := json.NewEncoder(w)
		enc.Encode(rfmlIDs)
	})

	out, _ := client.GetRFMLIDs()

	if !reflect.DeepEqual(rfmlIDs, out) {
		t.Errorf("Response expected = %v, actual %v", rfmlIDs, out)
	}
}

func TestGetTests(t *testing.T) {
	setup()
	defer cleanup()

	// Empty query
	rfFilters := RFTestFilters{}
	expectedQuery := url.Values{"page": []string{"1"}, "page_size": []string{"50"}}
	mux.HandleFunc("/tests", func(w http.ResponseWriter, r *http.Request) {
		receivedQuery := r.URL.Query()
		if !reflect.DeepEqual(expectedQuery, receivedQuery) {
			t.Errorf("Unexpected query sent to Rainforest API. Got %v, want %v", receivedQuery, expectedQuery)
		}

		w.Header().Add("X-Total-Pages", "1")
		w.Write([]byte("[]"))
	})

	_, err := client.GetTests(&rfFilters)
	if err != nil {
		t.Error(err.Error())
	}

	// Non-empty query
	rfFilters = RFTestFilters{
		Tags:          []string{"foo", "bar"},
		SiteID:        123,
		SmartFolderID: 321,
		RunGroupID:    237,
	}
	expectedQuery = url.Values{
		"page":            []string{"1"},
		"page_size":       []string{"50"},
		"tags":            rfFilters.Tags,
		"site_id":         []string{strconv.Itoa(rfFilters.SiteID)},
		"smart_folder_id": []string{strconv.Itoa(rfFilters.SmartFolderID)},
		"run_group_id":    []string{strconv.Itoa(rfFilters.RunGroupID)},
	}

	_, err = client.GetTests(&rfFilters)
	if err != nil {
		t.Error(err.Error())
	}

	// Multiple pages of tests
	cleanup()
	setup()

	currentPage := 1
	totalPages := 5
	mux.HandleFunc("/tests", func(w http.ResponseWriter, r *http.Request) {
		if currentPage > totalPages {
			t.Errorf("Page size received is greater than total pages: %v", currentPage)
		}

		receivedQuery := r.URL.Query()
		if receivedPageSize := receivedQuery.Get("page_size"); receivedPageSize != "50" {
			t.Errorf("Unexpected page size query: %v", receivedPageSize)
		}

		if receivedPage := receivedQuery.Get("page"); receivedPage != strconv.Itoa(currentPage) {
			t.Errorf("Expected page received. Expected %v, Got %v", currentPage, receivedPage)
		}

		currentPage++

		w.Header().Add("X-Total-Pages", "1")
		w.Write([]byte("[]"))
	})

	rfFilters = RFTestFilters{}
	expectedQuery = url.Values{"page": []string{"1"}, "page_size": []string{"50"}}

	_, err = client.GetTests(&rfFilters)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestHasUploadableFiles(t *testing.T) {
	// No uploadables
	test := RFTest{
		Steps: []interface{}{
			RFTestStep{
				Action:   "nothing here",
				Response: "or here",
			},
			RFEmbeddedTest{
				RFMLID: "definitely_nothing_here",
			},
		},
	}
	if test.HasUploadableFiles() {
		t.Error("Test has no uploadable files")
	}

	// With file download
	test.Steps = []interface{}{
		RFTestStep{
			Action:   "{{ file.download(./my/path) }}",
			Response: "nothing",
		},
	}
	if !test.HasUploadableFiles() {
		t.Error("Test has uploadable files")
	}

	// With screenshot
	test.Steps = []interface{}{
		RFTestStep{
			Action:   "{{ file.screenshot(./my/path) }}",
			Response: "nothing",
		},
	}
	if !test.HasUploadableFiles() {
		t.Error("Test has uploadable files")
	}

	// With missing argument
	test.Steps = []interface{}{
		RFTestStep{
			Action:   "{{ file.download }}",
			Response: "nothing",
		},
	}
	if test.HasUploadableFiles() {
		t.Error("Test should not have any uploadable files without an argument")
	}
}

func TestUpdateTest(t *testing.T) {
	// Test just the required attributes
	rfTest := RFTest{
		TestID:   123,
		RFMLID:   "an_rfml_id",
		Title:    "a title",
		StartURI: "/",
	}
	rfTest.PrepareToUploadFromRFML(TestIDMappings{})

	setup()
	defer cleanup()

	var data []byte
	var err error
	var bodyStr string
	mux.HandleFunc("/tests/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Incorrect HTTP method - expected PUT, got %v", r.Method)
			return
		}

		data, err = ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		bodyStr = string(data)
		if !strings.Contains(bodyStr, "\"id\":123") {
			t.Errorf("Correct test ID not received. Got: %v", bodyStr)
		} else if !strings.Contains(bodyStr, "\"rfml_id\":\"an_rfml_id\"") {
			t.Errorf("Unexpected RFML ID received. Expected: \"an_rfml_id\", Got: %v", bodyStr)
		} else if !strings.Contains(bodyStr, "\"title\":\"a title\"") {
			t.Errorf("Unexpected title received. Expected: \"a title\", Got:%v", bodyStr)
		} else if !strings.Contains(bodyStr, "\"start_uri\":\"/\"") {
			t.Errorf("Unexpected start URI received. Expected: \"/\", Got:%v", bodyStr)
		} else if !strings.Contains(bodyStr, "\"browsers\":[]") {
			t.Errorf("Unexpected browsers parameter received. Got:%v", bodyStr)
		} else if !strings.Contains(bodyStr, "\"tags\":[]") {
			t.Errorf("Unexpected tags parameter received. Got:%v", bodyStr)
		}
	})

	err = client.UpdateTest(&rfTest)
	if err != nil {
		t.Error(err.Error())
	}

	// // With extra attributes
	rfTest.Browsers = []string{"chrome", "firefox"}
	rfTest.Tags = []string{"foo", "bar"}
	rfTest.mapBrowsers()

	cleanup()
	setup()

	mux.HandleFunc("/tests/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Incorrect HTTP method - expected PUT, got %v", r.Method)
			return
		}

		var data []byte
		data, err = ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		bodyStr = string(data)

		if !strings.Contains(bodyStr, "\"name\":\"chrome\"") || !strings.Contains(bodyStr, "\"name\":\"firefox\"") {
			t.Errorf("Unexpected browsers received. Expected: \"chrome\", \"firefox\", Got:%v", bodyStr)
		} else if !strings.Contains(bodyStr, "\"tags\":[\"foo\",\"bar\"]") {
			t.Errorf("Unexpected browsers received. Expected: \"foo\", \"bar\", Got:%v", bodyStr)
		}
	})

	err = client.UpdateTest(&rfTest)
	if err != nil {
		t.Error(err.Error())
	}

	// Empty browsers and tags list
	rfTest.Browsers = []string{}
	rfTest.Tags = []string{}
	rfTest.mapBrowsers()

	cleanup()
	setup()

	mux.HandleFunc("/tests/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Incorrect HTTP method - expected PUT, got %v", r.Method)
			return
		}

		var data []byte
		data, err = ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		bodyStr = string(data)

		if !strings.Contains(bodyStr, "\"browsers\":[]") {
			t.Errorf("Unexpected browsers received. Expected: [], Got: %v", bodyStr)
		} else if !strings.Contains(bodyStr, "\"tags\":[]") {
			t.Errorf("Unexpected tags received. Expected: [], Got: %v", bodyStr)
		}
	})

	err = client.UpdateTest(&rfTest)
	if err != nil {
		t.Error(err.Error())
	}
}
