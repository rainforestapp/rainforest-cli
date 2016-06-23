package main

import "testing"

type fakeResGetter struct{}

func (g fakeResGetter) getFolders() (tableData [][]string) {
	tableData = [][]string{{"Folders"}}
	return
}
func (g fakeResGetter) getSites() (tableData [][]string) {
	tableData = [][]string{{"Sites"}}
	return
}
func (g fakeResGetter) getBrowsers() (tableData [][]string) {
	tableData = [][]string{{"Browsers"}}
	return
}

var fakeRes fakeResGetter

func TestFetchResource(t *testing.T) {
	var testCases = []struct {
		expectedResource string
	}{
		{expectedResource: "Folders"},
		{expectedResource: "Sites"},
		{expectedResource: "Browsers"},
	}
	for _, tcase := range testCases {
		chosenResource := fetchResource(tcase.expectedResource, fakeRes)
		if chosenResource[0][0] != tcase.expectedResource {
			t.Error("fetchResource did not choose the correct resource")
		}
	}
}
