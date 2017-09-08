package rainforest

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
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
			},
			expectedQuery: url.Values{
				"page":            []string{"1"},
				"page_size":       []string{"50"},
				"tags":            []string{"foo", "bar"},
				"site_id":         []string{"123"},
				"smart_folder_id": []string{"321"},
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
