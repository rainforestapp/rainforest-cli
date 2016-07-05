package main

import (
	"os"
	"testing"
)

// func newTestServer(path, resp string, statusCode int, t *testing.T) *httptest.Server {
// 	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.URL.Path != path {
// 			t.Errorf("fetchRource hit wrong endpoint (wanted %v but got %v)", path, r.URL.Path)
// 		}
// 		w.WriteHeader(statusCode)
// 		w.Write([]byte(resp))
// 	}))
// }

func TestRunByTags(t *testing.T) {
	sitesResp := "Post Request Successful"
	ts := newTestServer("/runs.json?tags=foo%2Cbar", sitesResp, 200, t)
	defer ts.Close()
	baseURL = ts.URL
	tempOsArgs := os.Args
	os.Args = []string{"rainforest-cli", "-tags=foo,bar", "run"}
	main()
	os.Args = tempOsArgs

}
