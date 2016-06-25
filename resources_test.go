package main

import (
	"os"
	"testing"
)

var fakeResourceGetter = webResGetter{
	getFolders: func() [][]string {
		return [][]string{{"Folders"}}
	},
	getSites: func() [][]string {
		return [][]string{{"Sites"}}
	},
	getBrowsers: func() [][]string {
		return [][]string{{"Browsers"}}
	},
}

func TestFetchResource(t *testing.T) {
	old := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	web = fakeResourceGetter
	var testCases = []struct {
		expectedResource string
	}{
		{expectedResource: "Folders"},
		{expectedResource: "Sites"},
		{expectedResource: "Browsers"},
	}
	for _, tcase := range testCases {
		chosenResource := fetchResource(tcase.expectedResource)
		if chosenResource[0][0] != tcase.expectedResource {
			t.Error("fetchResource did not choose the correct resource")
		}
	}
	os.Stdout = old
}
