package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

func newTestServer(path, resp string, statusCode int, t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			t.Errorf("fetchRource hit wrong endpoint (wanted %v but got %v)", path, r.URL.Path)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))
	}))
}

func runErrorTest(resource string, t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		printSites()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestPrint"+resource+"ApiError")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want status 1", err)
}

func checkTableCorrect(pattern string, t *testing.T) {
	matched, err := regexp.Match(pattern, out.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}
	if !matched {
		t.Logf("Table didn't match properly:")
		t.Logf("%v\n", out)
		t.Errorf("should have matched %v", pattern)
	}
}

func TestPrintSites(t *testing.T) {
	sitesResp := `[{"id": 1337, "name": "Dyer"}]`
	ts := newTestServer("/sites.json", sitesResp, 200, t)
	defer ts.Close()
	baseURL = ts.URL

	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()

	printSites()
	checkTableCorrect(`\| +SITE ID +\| +SITE DESCRIPTION +\|`, t)
	checkTableCorrect(`\| +1337 +\| +Dyer +\|`, t)
}

func TestPrintSitesApiError(t *testing.T) {
	sitesResp := `{"error": "This is a bad thing"}`
	ts := newTestServer("/sites.json", sitesResp, 400, t)
	defer ts.Close()
	baseURL = ts.URL
	runErrorTest("Sites", t)
}

func TestPrintFolders(t *testing.T) {
	sitesResp := `[{"id": 707, "title": "Foo"}]`
	ts := newTestServer("/folders.json", sitesResp, 200, t)
	defer ts.Close()
	baseURL = ts.URL

	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()

	printFolders()
	checkTableCorrect(`\| +FOLDER ID +\| +FOLDER DESCRIPTION +\|`, t)
	checkTableCorrect(`\| +707 +\| +Foo +\|`, t)
}

func TestPrintFoldersApiError(t *testing.T) {
	sitesResp := `{"error": "This is a bad thing"}`
	ts := newTestServer("/folders.json", sitesResp, 600, t)
	defer ts.Close()
	baseURL = ts.URL

	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()
	runErrorTest("Folders", t)
}

func TestPrintBrowsers(t *testing.T) {
	sitesResp := `{"available_browsers": [{"name": "firefox", "description": "Mozilla Firefox"}]}`
	ts := newTestServer("/clients.json", sitesResp, 200, t)
	defer ts.Close()
	baseURL = ts.URL

	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()

	printBrowsers()
	checkTableCorrect(`\| +BROWSER ID +\| +BROWSER DESCRIPTION +\|`, t)
	checkTableCorrect(`\| +firefox +\| +Mozilla Firefox +\|`, t)
}

func TestPrintBrowsersApiError(t *testing.T) {
	sitesResp := `{"error": "This is a bad thing"}`
	ts := newTestServer("/clients.json", sitesResp, 600, t)
	defer ts.Close()
	baseURL = ts.URL
	out = &bytes.Buffer{}
	defer func() {
		out = os.Stdout
	}()
	runErrorTest("Browsers", t)
}

// package main

// import (
// 	"reflect"
// 	"testing"
// )

// func TestSitessDisplayTable(t *testing.T) {
// 	var SitesTestCases = []struct {
// 		testStruct sitesResp
// 		want       [][]string
// 	}{
// 		{
// 			testStruct: sitesResp{
// 				sites{
// 					ID:   1337,
// 					Name: "Dyer",
// 				},
// 				sites{
// 					ID:   42,
// 					Name: "Situation",
// 				},
// 			},
// 			want: [][]string{{"1337", "Dyer"}, {"42", "Situation"}},
// 		},
// 	}
// 	for _, tcase := range SitesTestCases {
// 		got := tcase.testStruct.displayTable()
// 		if !reflect.DeepEqual(tcase.want, got) {
// 			t.Log("want:")
// 			t.Logf("\t%+v", tcase.want)
// 			t.Log("got =")
// 			t.Errorf("\t%+v", got)
// 		}
// 	}

// }

// func TestBrowsersDisplayTable(t *testing.T) {
// 	var BrowsersTestCases = []struct {
// 		testStruct browsersResp
// 		want       [][]string
// 	}{
// 		{
// 			testStruct: browsersResp{
// 				AvailableBrowsers: []browser{
// 					{
// 						Name:        "firefox",
// 						Description: "Mozilla Firefox",
// 					},
// 					{
// 						Name:        "ie11",
// 						Description: "Microsoft Internet Explorer 11",
// 					},
// 				},
// 			},
// 			want: [][]string{{"firefox", "Mozilla Firefox"}, {"ie11", "Microsoft Internet Explorer 11"}},
// 		},
// 	}

// 	for _, tcase := range BrowsersTestCases {
// 		got := tcase.testStruct.displayTable()
// 		if !reflect.DeepEqual(tcase.want, got) {
// 			t.Log("want:")
// 			t.Logf("\t%+v", tcase.want)
// 			t.Log("got =")
// 			t.Errorf("\t%+v", got)
// 		}
// 	}
// }

// func TestFoldersDisplayTable(t *testing.T) {
// 	var foldersTestCases = []struct {
// 		testStruct foldersResp
// 		want       [][]string
// 	}{
// 		{
// 			testStruct: foldersResp{
// 				folder{
// 					ID:    707,
// 					Title: "The Foo Folder",
// 				},
// 				folder{
// 					ID:    708,
// 					Title: "The Baz Folder",
// 				},
// 			},
// 			want: [][]string{{"707", "The Foo Folder"}, {"708", "The Baz Folder"}},
// 		},
// 	}

// 	for _, tcase := range foldersTestCases {
// 		got := tcase.testStruct.displayTable()
// 		if !reflect.DeepEqual(tcase.want, got) {
// 			t.Log("want:")
// 			t.Logf("\t%+v", tcase.want)
// 			t.Log("got =")
// 			t.Errorf("\t%+v", got)
// 		}
// 	}

// }
