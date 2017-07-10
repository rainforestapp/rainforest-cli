package main

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
)

func regexMatchOut(pattern string, t *testing.T) {
	matched, err := regexp.Match(pattern, tablesOut.(*bytes.Buffer).Bytes())
	if err != nil {
		t.Error("Error with pattern match:", err)
	}
	if !matched {
		t.Errorf("Printed out %v, want %v", tablesOut, pattern)
	}
}

func TestPrintResourceTable(t *testing.T) {
	tablesOut = &bytes.Buffer{}
	defer func() {
		tablesOut = os.Stdout
	}()

	testBody := [][]string{{"1337", "Dyer"}, {"42", "Situation"}}
	printResourceTable([]string{"Column 1", "Column 2"}, testBody)
	regexMatchOut(`\| +COLUMN 1 +\| +COLUMN 2 +\|`, t)
	regexMatchOut(`\| +1337 +\| +Dyer +\|`, t)
}

type testResourceAPI struct {
	Folders   []rainforest.Folder
	Browsers  []rainforest.Browser
	Sites     []rainforest.Site
	RunGroups []rainforest.RunGroup
}

func (api testResourceAPI) GetFolders() ([]rainforest.Folder, error) {
	return api.Folders, nil
}

func (api testResourceAPI) GetBrowsers() ([]rainforest.Browser, error) {
	return api.Browsers, nil
}

func (api testResourceAPI) GetSites() ([]rainforest.Site, error) {
	return api.Sites, nil
}

func (api testResourceAPI) GetRunGroups() ([]rainforest.RunGroup, error) {
	return api.RunGroups, nil
}

func TestPrintFolders(t *testing.T) {
	tablesOut = &bytes.Buffer{}
	defer func() {
		tablesOut = os.Stdout
	}()

	testAPI := testResourceAPI{
		Folders: []rainforest.Folder{
			{ID: 123, Title: "First Folder Title"},
			{ID: 456, Title: "Second Folder Title"},
		},
	}

	printFolders(testAPI)
	regexMatchOut(`\| +FOLDER ID +\| +FOLDER NAME +\|`, t)
	regexMatchOut(`\| +123 +\| +First Folder Title +\|`, t)
	regexMatchOut(`\| +456 +\| +Second Folder Title +\|`, t)
}

func TestPrintBrowsers(t *testing.T) {
	tablesOut = &bytes.Buffer{}
	defer func() {
		tablesOut = os.Stdout
	}()

	testAPI := testResourceAPI{
		Browsers: []rainforest.Browser{
			{Name: "chrome", Description: "Google Chrome"},
			{Name: "firefox", Description: "Mozilla Firefox"},
		},
	}

	printBrowsers(testAPI)
	regexMatchOut(`\| +BROWSER ID +\| +BROWSER NAME +\|`, t)
	regexMatchOut(`\| +chrome +\| +Google Chrome +\|`, t)
	regexMatchOut(`\| +firefox +\| +Mozilla Firefox +\|`, t)
}

func TestPrintSites(t *testing.T) {
	tablesOut = &bytes.Buffer{}
	defer func() {
		tablesOut = os.Stdout
	}()

	testAPI := testResourceAPI{
		Sites: []rainforest.Site{
			{ID: 123, Name: "My favorite site", Category: "site"},
			{ID: 456, Name: "My favorite app URL", Category: "ios"},
			{ID: 789, Name: "Site with unknown platform", Category: "unknown_platform"},
		},
	}

	printSites(testAPI)
	regexMatchOut(`\| +SITE ID +\| +SITE NAME +\| +CATEGORY +\|`, t)
	regexMatchOut(`\| +123 +\| +My favorite site +\| +Site +\|`, t)
	regexMatchOut(`\| +456 +\| +My favorite app URL +\| +iOS +\|`, t)
	regexMatchOut(`\| +789 +\| +Site with unknown platform +\| +unknown_platform +\|`, t)
}

func TestPrintRunGroups(t *testing.T) {
	tablesOut = &bytes.Buffer{}
	defer func() {
		tablesOut = os.Stdout
	}()

	testAPI := testResourceAPI{
		RunGroups: []rainforest.RunGroup{
			{ID: 59, Title: "Run Group Race"},
			{ID: 245, Title: "Run Group Marathon"},
		},
	}

	printRunGroups(testAPI)
	regexMatchOut(`\| +RUN GROUP ID +\| +RUN GROUP NAME +\|`, t)
	regexMatchOut(`\| +59 +\| +Run Group Race +\|`, t)
	regexMatchOut(`\| +245 +\| +Run Group Marathon +\|`, t)
}
