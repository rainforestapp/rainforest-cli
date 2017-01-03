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
	setup()
	defer cleanup()

	// Empty query
	rfFilters := RFTestFilters{}
	expectedQuery := url.Values{}
	mux.HandleFunc("/tests", func(w http.ResponseWriter, r *http.Request) {
		receivedQuery := r.URL.Query()
		if !reflect.DeepEqual(expectedQuery, receivedQuery) {
			t.Errorf("Unexpected query sent to Rainforest API. Got %v, want %v", receivedQuery, expectedQuery)
		}

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
	}
	expectedQuery = url.Values{
		"tags":            rfFilters.Tags,
		"site_id":         []string{strconv.Itoa(rfFilters.SiteID)},
		"smart_folder_id": []string{strconv.Itoa(rfFilters.SmartFolderID)},
	}

	_, err = client.GetTests(&rfFilters)
	if err != nil {
		t.Error(err.Error())
	}
}
