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
		Pairs: []TestIDPair{
			{ID: 123, RFMLID: "abc"},
			{ID: 456, RFMLID: "xyz"},
		},
	}

	mux.HandleFunc("/tests/rfml_ids", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		enc := json.NewEncoder(w)
		enc.Encode(rfmlIDs.Pairs)
	})

	out, err := client.GetRFMLIDs()
	if err != nil {
		t.Error(err.Error())
	} else if !reflect.DeepEqual(&rfmlIDs, out) {
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
		var data []byte
		data, err = json.Marshal(test)
		if err != nil {
			t.Fatal("Error marshalling test:", err)
		}
		w.Header().Add("X-Total-Pages", "1")
		w.Write(data)
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

	requestCount := 0
	mux.HandleFunc("/tests/123", func(w http.ResponseWriter, r *http.Request) {
		requestCount += 1

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
	} else if test.TestID != 123 || test.RFMLID != "123" || test.Title != "A test" {
		t.Errorf("test %v was unmarshalled incorrectly", test)
	} else if !test.Execute {
		t.Error("GetTest didn't set execute: true by default")
	}

	test2, err := client.GetTest(123)
	if err != nil {
		t.Error("Error fetching test:", err.Error())
	} else if requestCount > 1 {
		t.Error("Request was not cached")
	} else if test != test2 {
		t.Error("Expected to receive cached test, got another test")
	}
}

func TestPrepareToWriteAsRFML(t *testing.T) {
	client = NewClient("foo", false)

	// No test elements
	test := RFTest{}
	err := test.PrepareToWriteAsRFML(client, false)
	if err != nil {
		t.Error(err.Error())
	}

	if test.Steps != nil {
		t.Error("Unexpected steps appeared in test")
	}

	// Just step elements
	test = RFTest{
		Elements: []testElement{
			{
				Redirect: false,
				Type:     "step",
				Details: testElementDetails{
					Action:   "first action",
					Response: "first response",
				},
			},
			{
				Redirect: true,
				Type:     "step",
				Details: testElementDetails{
					Action:   "second action",
					Response: "second response",
				},
			},
		},
	}

	err = test.PrepareToWriteAsRFML(client, false)
	if err != nil {
		t.Error(err.Error())
	}
	expectedSteps := []interface{}{
		RFTestStep{Action: "first action", Response: "first response", Redirect: false},
		RFTestStep{Action: "second action", Response: "second response", Redirect: true},
	}
	for i, expectedStep := range expectedSteps {
		actualStep := test.Steps[i]
		if !reflect.DeepEqual(expectedStep, actualStep) {
			t.Errorf("Unexpected step.\nExpected:\n%v\nGot:\n%v", expectedStep, actualStep)
		}
	}

	// Step and test elements, all embedded
	test = RFTest{
		Elements: []testElement{
			{
				Redirect: false,
				Type:     "step",
				Details: testElementDetails{
					Action:   "first action",
					Response: "first response",
				},
			},
			{
				Redirect: true,
				Type:     "test",
				Details:  testElementDetails{ID: 778899},
			},
			{
				Redirect: true,
				Type:     "step",
				Details: testElementDetails{
					Action:   "second action",
					Response: "second response",
				},
			},
		},
	}

	setup()
	defer cleanup()

	mux.HandleFunc("/tests/rfml_ids", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Incorrect HTTP verb. Expected GET. Got %v", r.Method)
		}

		idPairs := []TestIDPair{{ID: 778899, RFMLID: "an_rfml_id"}}
		var b []byte
		b, err = json.Marshal(&idPairs)
		if err != nil {
			t.Fatal("Error marshalling test IDs: ", err.Error())
		}
		w.Write(b)
	})

	err = test.PrepareToWriteAsRFML(client, true)
	if err != nil {
		t.Error(err.Error())
	}
	expectedSteps = []interface{}{
		RFTestStep{Action: "first action", Response: "first response", Redirect: false},
		RFEmbeddedTest{RFMLID: "an_rfml_id", Redirect: true},
		RFTestStep{Action: "second action", Response: "second response", Redirect: true},
	}
	for i, expectedStep := range expectedSteps {
		actualStep := test.Steps[i]
		if !reflect.DeepEqual(expectedStep, actualStep) {
			t.Errorf("Unexpected step.\nExpected:\n%v\nGot:\n%v", expectedStep, actualStep)
		}
	}

	// Step and deeply embedded test elements, all flattened
	cleanup()
	setup()
	mux.HandleFunc("/tests/rfml_ids", func(http.ResponseWriter, *http.Request) {
		t.Fatal("This request shouldn't have been made!")
	})
	mux.HandleFunc("/tests/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Incorrect HTTP verb. Expected GET. Got %v", r.Method)
		}

		var returnedTest RFTest
		if r.URL.Path == "/tests/778899" {
			returnedTest = RFTest{
				Elements: []testElement{
					{
						Redirect: true,
						Type:     "step",
						Details: testElementDetails{
							Action:   "action in first embedded test",
							Response: "response in first embedded test",
						},
					},
					{
						Redirect: false,
						Type:     "test",
						Details:  testElementDetails{ID: 543543},
					},
				},
			}
		} else if r.URL.Path == "/tests/543543" {
			returnedTest = RFTest{
				Elements: []testElement{
					{
						Redirect: true,
						Type:     "step",
						Details: testElementDetails{
							Action:   "action in second embedded test",
							Response: "response in second embedded test",
						},
					},
				},
			}
		} else {
			t.Error("Unexpected path: ", r.URL.Path)
		}

		var b []byte
		b, err = json.Marshal(&returnedTest)
		if err != nil {
			t.Fatal("Error marshalling test: ", err.Error())
		}
		w.Write(b)
	})

	err = test.PrepareToWriteAsRFML(client, false)
	if err != nil {
		t.Error(err.Error())
	}
	expectedSteps = []interface{}{
		RFTestStep{Action: "first action", Response: "first response", Redirect: false},
		RFTestStep{Action: "action in first embedded test", Response: "response in first embedded test", Redirect: true},
		RFTestStep{Action: "action in second embedded test", Response: "response in second embedded test", Redirect: true},
		RFTestStep{Action: "second action", Response: "second response", Redirect: true},
	}
	if !reflect.DeepEqual(expectedSteps, test.Steps) {
		t.Errorf("Unexpected steps.\nExpected:\n%v\nGot:\n%v", expectedSteps, test.Steps)
	}

	// Browsers
	test = RFTest{
		BrowsersMap: []map[string]interface{}{
			{"state": "enabled", "name": "some_browser"},
			{"state": "disabled", "name": "disabled_browser"},
			{"state": "enabled", "name": "another_browser"},
		},
	}

	err = test.PrepareToWriteAsRFML(client, true)
	if err != nil {
		t.Error(err.Error())
	}
	expectedBrowsers := []string{"some_browser", "another_browser"}
	if !reflect.DeepEqual(expectedBrowsers, test.Browsers) {
		t.Errorf("Unexpected browsers.\nExpected:\n%v\nGot:\n%v", expectedBrowsers, test.Browsers)
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
	rfTest.PrepareToUploadFromRFML(&TestIDMappings{})

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
		} else if strings.Contains(bodyStr, "\"feature_id\"") {
			t.Errorf("Unexpected parameter found: \"feature_id\" in:\n%v", bodyStr)
		}
	})

	err = client.UpdateTest(&rfTest)
	if err != nil {
		t.Error(err.Error())
	}

	// // With extra attributes
	rfTest.Browsers = []string{"chrome", "firefox"}
	rfTest.Tags = []string{"foo", "bar"}
	rfTest.FeatureID = 909
	rfTest.State = "disabled"

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
			t.Errorf("Expected browsers not received. Expected: \"chrome\", \"firefox\", Got: %v", bodyStr)
		}
		if !strings.Contains(bodyStr, "\"tags\":[\"foo\",\"bar\"]") {
			t.Errorf("Expected tags not received. Expected: \"foo\", \"bar\", Got: %v", bodyStr)
		}
		if !strings.Contains(bodyStr, "\"feature_id\":909") {
			t.Errorf("Expected feature ID not received. Expected: 909, Got %v", bodyStr)
		}
		if !strings.Contains(bodyStr, "\"state\":\"disabled\"") {
			t.Errorf("Expected state to be disabled, Got %v", bodyStr)
		}
	})

	err = client.UpdateTest(&rfTest)
	if err != nil {
		t.Error(err.Error())
	}

	// Deleted feature ID, empty browsers and tags list
	rfTest.FeatureID = deleteFeature
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
		} else if !strings.Contains(bodyStr, "\"feature_id\":null") {
			t.Errorf("Unexpected folder ID received. Expected: null, Got: %v", bodyStr)
		}
	})

	err = client.UpdateTest(&rfTest)
	if err != nil {
		t.Error(err.Error())
	}
}
