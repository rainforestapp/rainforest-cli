package main

import (
	"reflect"
	"testing"
)

func TestSitessDisplayTable(t *testing.T) {
	var SitesTestCases = []struct {
		testStruct sitesResp
		want       [][]string
	}{
		{
			testStruct: sitesResp{
				sites{
					ID:   1337,
					Name: "Dyer",
				},
				sites{
					ID:   42,
					Name: "Situation",
				},
			},
			want: [][]string{{"1337", "Dyer"}, {"42", "Situation"}},
		},
	}
	for _, tcase := range SitesTestCases {
		got := tcase.testStruct.displayTable()
		if !reflect.DeepEqual(tcase.want, got) {
			t.Log("want:")
			t.Logf("\t%+v", tcase.want)
			t.Log("got =")
			t.Errorf("\t%+v", got)
		}
	}

}

func TestBrowsersDisplayTable(t *testing.T) {
	var BrowsersTestCases = []struct {
		testStruct browsersResp
		want       [][]string
	}{
		{
			testStruct: browsersResp{
				AvailableBrowsers: []browser{
					{
						Name:        "firefox",
						Description: "Mozilla Firefox",
					},
					{
						Name:        "ie11",
						Description: "Microsoft Internet Explorer 11",
					},
				},
			},
			want: [][]string{{"firefox", "Mozilla Firefox"}, {"ie11", "Microsoft Internet Explorer 11"}},
		},
	}

	for _, tcase := range BrowsersTestCases {
		got := tcase.testStruct.displayTable()
		if !reflect.DeepEqual(tcase.want, got) {
			t.Log("want:")
			t.Logf("\t%+v", tcase.want)
			t.Log("got =")
			t.Errorf("\t%+v", got)
		}
	}
}

func TestFoldersDisplayTable(t *testing.T) {
	var foldersTestCases = []struct {
		testStruct foldersResp
		want       [][]string
	}{
		{
			testStruct: foldersResp{
				folder{
					ID:    707,
					Title: "The Foo Folder",
				},
				folder{
					ID:    708,
					Title: "The Baz Folder",
				},
			},
			want: [][]string{{"707", "The Foo Folder"}, {"708", "The Baz Folder"}},
		},
	}

	for _, tcase := range foldersTestCases {
		got := tcase.testStruct.displayTable()
		if !reflect.DeepEqual(tcase.want, got) {
			t.Log("want:")
			t.Logf("\t%+v", tcase.want)
			t.Log("got =")
			t.Errorf("\t%+v", got)
		}
	}

}
