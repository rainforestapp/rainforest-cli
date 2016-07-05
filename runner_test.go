package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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
	tempOsArgs := os.Args
	os.Args = []string{"rainforest-cli", "-tags=foo,bar", "run"}
	main()
	os.Args = tempOsArgs
}

// func TestRunBySmartID(t *testing.T) {
// 	sitesResp := "Post Request Successful"
// 	ts := newTestServer("/runs.json?tags=foo%2Cbar", sitesResp, 200, t)
// 	defer ts.Close()
// 	baseURL = ts.URL
// 	tempOsArgs := os.Args
// 	os.Args = []string{"rainforest-cli", "-tags=foo,bar", "run"}
// 	main()
// 	os.Args = tempOsArgs
// }
