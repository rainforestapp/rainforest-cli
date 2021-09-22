package main

import (
	"bytes"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
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
	Folders      []rainforest.Folder
	Browsers     []rainforest.Browser
	Sites        []rainforest.Site
	Environments []rainforest.Environment
	Features     []rainforest.Feature
	RunGroups    []rainforest.RunGroup
	Junit        string
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

func (api testResourceAPI) GetEnvironments() ([]rainforest.Environment, error) {
	return api.Environments, nil
}

func (api testResourceAPI) GetFeatures() ([]rainforest.Feature, error) {
	return api.Features, nil
}

func (api testResourceAPI) GetRunGroups() ([]rainforest.RunGroup, error) {
	return api.RunGroups, nil
}

func (api testResourceAPI) GetRunJunit(run_id int) (*string, error) {
	return &api.Junit, nil
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

func TestPrintEnvironments(t *testing.T) {
	tablesOut = &bytes.Buffer{}
	defer func() {
		tablesOut = os.Stdout
	}()

	testAPI := testResourceAPI{
		Environments: []rainforest.Environment{
			{ID: 123, Name: "QA"},
			{ID: 456, Name: "Staging 1"},
			{ID: 789, Name: "Staging 2"},
		},
	}

	printEnvironments(testAPI)
	regexMatchOut(`\| +ENVIRONMENT ID +\| +ENVIRONMENT NAME +\|`, t)
	regexMatchOut(`\| +123 +\| +QA +\|`, t)
	regexMatchOut(`\| +456 +\| +Staging 1 +\|`, t)
	regexMatchOut(`\| +789 +\| +Staging 2 +\|`, t)
}

func TestPrintFeatures(t *testing.T) {
	tablesOut = &bytes.Buffer{}
	defer func() {
		tablesOut = os.Stdout
	}()

	testAPI := testResourceAPI{
		Features: []rainforest.Feature{
			{ID: 123, Title: "My favorite feature"},
			{ID: 456, Title: "My least favorite feature"},
			{ID: 789, Title: "An OK feature"},
		},
	}

	printFeatures(testAPI)
	regexMatchOut(`\| +FEATURE ID +\| +FEATURE TITLE +\|`, t)
	regexMatchOut(`\| +123 +\| +My favorite feature +\|`, t)
	regexMatchOut(`\| +456 +\| +My least favorite feature +\|`, t)
	regexMatchOut(`\| +789 +\| +An OK feature +\|`, t)
}

func TestPrintRunGroups(t *testing.T) {
	tablesOut = &bytes.Buffer{}
	defer func() {
		tablesOut = os.Stdout
	}()

	testAPI := testResourceAPI{
		RunGroups: []rainforest.RunGroup{
			{ID: 123, Title: "My favorite run group"},
			{ID: 456, Title: "My least favorite run group"},
			{ID: 789, Title: "An OK run group"},
		},
	}

	printRunGroups(testAPI)
	regexMatchOut(`\| +RUN GROUP ID +\| +RUN GROUP TITLE +\|`, t)
	regexMatchOut(`\| +123 +\| +My favorite run group +\|`, t)
	regexMatchOut(`\| +456 +\| +My least favorite run group +\|`, t)
	regexMatchOut(`\| +789 +\| +An OK run group +\|`, t)
}

func TestPostRunJUnitReport(t *testing.T) {
	// returns nil with no junit setting enabled
	fakeContext := newFakeContext(map[string]interface{}{"token": "test"}, cli.Args{"1"})
	err := postRunJUnitReport(fakeContext, 1)

	if err != nil {
		t.Errorf("postRunJUnitReport returned %+v", err)
	}
}

func TestWriteJunit(t *testing.T) {
	fakeContext := newFakeContext(map[string]interface{}{"junit-file": "junit.xml"}, cli.Args{"1"})
	testAPI := testResourceAPI{
		Junit: "<xml>hai</xml>",
	}
	err := writeJunit(fakeContext, testAPI)

	if err != nil {
		t.Errorf("writeJunit returned %+v", err)
	}

	data, _ := os.ReadFile("junit.xml")
	if !reflect.DeepEqual(testAPI.Junit, string(data)) {
		t.Errorf("writeJunit wrote %+v, want %+v", data, testAPI.Junit)
	}

	fakeContext = newFakeContext(map[string]interface{}{"junit-file": ""}, cli.Args{"1"})
	err = writeJunit(fakeContext, testAPI)
	expected := "JUnit output file not specified"
	if !reflect.DeepEqual(expected, err.Error()) {
		t.Errorf("writeJunit should have errored: expected '%v', got '%v'", expected, err.Error())
	}

	fakeContext = newFakeContext(map[string]interface{}{}, cli.Args{"1"})
	err = writeJunit(fakeContext, testAPI)
	expected = "JUnit output file not specified"
	if !reflect.DeepEqual(expected, err.Error()) {
		t.Errorf("writeJunit should have errored: expected '%v', got '%v'", expected, err.Error())
	}

	fakeContext = newFakeContext(map[string]interface{}{"junit-file": "junit.xml"}, cli.Args{})
	err = writeJunit(fakeContext, testAPI)
	expected = "No run ID argument found."
	if !reflect.DeepEqual(expected, err.Error()) {
		t.Errorf("writeJunit should have errored: expected '%v', got '%v'", expected, err.Error())
	}
}

func TestAugmentJunitFileName(t *testing.T) {
	testCases := []struct {
		OriginalName string
		RerunAttempt uint
		Want         string
	}{
		{
			OriginalName: "output",
			RerunAttempt: 1,
			Want:         "output.1",
		},
		{
			OriginalName: "output.xml",
			RerunAttempt: 1,
			Want:         "output.xml.1",
		},
		{
			OriginalName: "output.xml",
			RerunAttempt: 2,
			Want:         "output.xml.2",
		},
		{
			OriginalName: "output.xml.1",
			RerunAttempt: 3,
			Want:         "output.xml.1.3",
		},
		{
			OriginalName: "some.output.xml.1.2.3",
			RerunAttempt: 2,
			Want:         "some.output.xml.1.2.3.2",
		},
		{
			OriginalName: "output.html",
			RerunAttempt: 1,
			Want:         "output.html.1",
		},
		{
			OriginalName: "output.1",
			RerunAttempt: 2,
			Want:         "output.1.2",
		},
	}
	for _, testCase := range testCases {
		got := augmentJunitFileName(testCase.OriginalName, testCase.RerunAttempt)
		if got != testCase.Want {
			t.Errorf("Wrong name, wanted '%v', got '%v'", testCase.Want, got)
		}
	}
}
