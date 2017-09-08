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
	testCases := []struct {
		rfFilters     RFTestFilters
		expectedQuery url.Values
	}{
		// Empty query
		{
			rfFilters:     RFTestFilters{},
			expectedQuery: url.Values{"page": []string{"1"}, "page_size": []string{"50"}},
		},

		// Non-empty query
		{
			rfFilters: RFTestFilters{
				Tags:          []string{"foo", "bar"},
				SiteID:        123,
				SmartFolderID: 321,
				RunGroupID:    237,
			},
			expectedQuery: url.Values{
				"page":            []string{"1"},
				"page_size":       []string{"50"},
				"tags":            []string{"foo", "bar"},
				"site_id":         []string{"123"},
				"smart_folder_id": []string{"321"},
				"run_group_id":    []string{"237"},
			},
		},

		// Filter by feature and run group
		{
			rfFilters: RFTestFilters{
				FeatureID:  123,
				RunGroupID: 75,
			},
			expectedQuery: url.Values{
				"page":         []string{"1"},
				"page_size":    []string{"50"},
				"feature_id":   []string{"123"},
				"run_group_id": []string{"75"},
			},
		},
	}

	for _, tc := range testCases {
		setup()
		defer cleanup()

		// Empty query
		mux.HandleFunc("/tests", func(w http.ResponseWriter, r *http.Request) {
			receivedQuery := r.URL.Query()
			if !reflect.DeepEqual(tc.expectedQuery, receivedQuery) {
				t.Errorf("Unexpected query sent to Rainforest API. Got %v, want %v", receivedQuery, tc.expectedQuery)
			}

			w.Header().Add("X-Total-Pages", "1")
			w.Write([]byte("[]"))
		})

		_, err := client.GetTests(&tc.rfFilters)
		if err != nil {
			t.Error(err.Error())
		}

		_, err = client.GetTests(&tc.rfFilters)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestGetTestsMultiplePages(t *testing.T) {
	setup()
	defer cleanup()

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

	rfFilters := RFTestFilters{}

	_, err := client.GetTests(&rfFilters)
	if err != nil {
		t.Error(err.Error())
	}

	// Test deserialization
	cleanup()
	setup()
	mux.HandleFunc("/tests", func(w http.ResponseWriter, r *http.Request) {
		test := []*RFTest{
			{
				TestID: 123,
				RFMLID: "123",
			},
		}
		b, err := json.Marshal(test)
		if err != nil {
			t.Fatal("Error marshalling test:", err)
		}
		w.Header().Add("X-Total-Pages", "1")
		w.Write(b)
	})

	tests, err := client.GetTests(&RFTestFilters{})
	if len(tests) != 1 {
		t.Error("Invalid number of tests returned:", len(tests))
	}

	got := tests[0]
	want := &RFTest{
		TestID: 123,
		RFMLID: "123",
	}
	if got.TestID != want.TestID || got.RFMLID != want.RFMLID {
		t.Errorf("test (%v) deserialized incorrectly", got)
	}
	if !got.Execute {
		t.Error("GetTests didn't set execute: true by default")
	}
}

func TestGetTest(t *testing.T) {
	setup()
	defer cleanup()

	mux.HandleFunc("/tests/123", func(w http.ResponseWriter, r *http.Request) {
		test := &RFTest{
			TestID: 123,
			RFMLID: "123",
			Title:  "A test",
		}
		b, err := json.Marshal(test)
		if err != nil {
			t.Fatal("Error marshalling test:", err)
		}
		w.Write(b)
	})

	test, err := client.GetTest(123)
	if err != nil {
		t.Error("Error fetching test:", err)
	}
	if test.TestID != 123 || test.RFMLID != "123" || test.Title != "A test" {
		t.Errorf("test %v was unmarshalled incorrectly", test)
	}
	if !test.Execute {
		t.Error("GetTest didn't set execute: true by default")
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

	// Remote file download reference
	test.Steps = []interface{}{
		RFTestStep{
			Action:   "{{ file.download(1235, foobar, testing.csv) }}",
			Response: "nothing",
		},
	}

	if test.HasUploadableFiles() {
		t.Error("Test only has remote uploadable files")
	}

	// Remote screenshot reference
	test.Steps = []interface{}{
		RFTestStep{
			Action:   "{{ file.screenshot(3332, hkBde5) }}",
			Response: "nothing",
		},
	}

	if test.HasUploadableFiles() {
		t.Error("Test only has remote uploadable files")
	}

	// Remote and local reference
	test.Steps = []interface{}{
		RFTestStep{
			Action:   "{{ file.download(1235, foobar, testing.csv) }} {{ file.download(./local/path.csv) }}",
			Response: "nothing",
		},
	}

	if !test.HasUploadableFiles() {
		t.Error("Test should have an uploadable file")
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
