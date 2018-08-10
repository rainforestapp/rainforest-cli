package rainforest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const uploadableRegex = `{{ *file\.(download|screenshot)\(([^\)]+)\) *}}`

// TestIDPair is a type representing RF tests that contain the test definitions.
type TestIDPair struct {
	ID     int    `json:"id"`
	RFMLID string `json:"rfml_id"`
}

// TestIDFinder finds corresponding test IDs and RFML IDs.
type TestIDFinder struct {
	// TODO: Unexport this
	Pairs      []TestIDPair
	idToRFMLID map[int]string
	rfmlIDtoID map[string]int
}

// GetRFMLID finds the corresponding RFML ID for a test ID
func (finder TestIDFinder) GetRFMLID(testID int) (string, error) {
	if finder.idToRFMLID == nil {
		// NOTE: Assuming pairs is already populated
		finder.idToRFMLID = make(map[int]string, len(finder.Pairs))
		for _, pair := range finder.Pairs {
			finder.idToRFMLID[pair.ID] = pair.RFMLID
		}
	}

	rfmlID, ok := finder.idToRFMLID[testID]
	if !ok {
		return "", fmt.Errorf("Unable to find test with ID %v", testID)
	}

	return rfmlID, nil
}

// MapRFMLIDtoID creates a map from RFML IDs to IDs
func (finder TestIDFinder) GetTestID(rfmlID string) (int, error) {
	if finder.rfmlIDtoID == nil {
		// NOTE: Assuming pairs is already populated
		finder.rfmlIDtoID = make(map[string]int, len(finder.Pairs))
		for _, pair := range finder.Pairs {
			finder.rfmlIDtoID[pair.RFMLID] = pair.ID
		}
	}

	testID, ok := finder.rfmlIDtoID[rfmlID]
	if !ok {
		return 0, fmt.Errorf("Unable to find test with RFML ID %v", rfmlID)
	}

	return testID, nil
}

// FeatureIDInt is a wrapper the int type used in the FeatureID field of RFTest
// that implements the json.Marshaler interface.
type FeatureIDInt int

const deleteFeature = -1

// MarshalJSON treats an FeatureIDInt value of -1 as a special case so that it
// is marshalledd into a `null` JSON value.
func (id *FeatureIDInt) MarshalJSON() ([]byte, error) {
	intVal := int(*id)
	if intVal == deleteFeature {
		return json.Marshal(nil)
	}

	// If not -1, marshal like any other int
	return json.Marshal(intVal)
}

// RFTest is a struct representing the Rainforest Test with its settings and steps
type RFTest struct {
	TestID      int                      `json:"id"`
	RFMLID      string                   `json:"rfml_id"`
	Source      string                   `json:"source"`
	Title       string                   `json:"title,omitempty"`
	State       string                   `json:"state,omitempty"`
	StartURI    string                   `json:"start_uri"`
	SiteID      int                      `json:"site_id,omitempty"`
	Description string                   `json:"description,omitempty"`
	Tags        []string                 `json:"tags"`
	BrowsersMap []map[string]interface{} `json:"browsers"`
	Elements    []testElement            `json:"elements,omitempty"`
	FeatureID   FeatureIDInt             `json:"feature_id,omitempty"`

	// Browsers and Steps are helper fields
	Browsers []string      `json:"-"`
	Steps    []interface{} `json:"-"`
	// RFMLPath is a helper field for keeping track of the filepath to the
	// test's RFML file.
	RFMLPath string `json:"-"`

	// Execute is a non-API field that specifies whether the test should be
	// executed or just uploaded (e.g. for embedded tests). It defaults to
	// true when reading from RFML.
	Execute bool `json:"-"`
}

// testElement is one of the helpers to construct the proper JSON test sturcture
type testElement struct {
	Redirect bool               `json:"redirection"`
	Type     string             `json:"type"`
	Details  testElementDetails `json:"element"`
}

// testElementDetails is one of the helpers to construct the proper JSON test sturcture
type testElementDetails struct {
	ID       int    `json:"id,omitempty"`
	Action   string `json:"action,omitempty"`
	Response string `json:"response,omitempty"`
}

// mapBrowsers fills the browsers field with format recognized by the API
func (t *RFTest) mapBrowsers() {
	t.BrowsersMap = make([]map[string]interface{}, len(t.Browsers))
	for i, browser := range t.Browsers {
		mappedBrowser := map[string]interface{}{
			"state": "enabled",
			"name":  browser,
		}
		t.BrowsersMap[i] = mappedBrowser
	}
}

// unmapBrowsers parses browsers from the API format to internal go one
func (t *RFTest) unmapBrowsers() {
	t.Browsers = []string{}
	for _, browserMap := range t.BrowsersMap {
		if browserMap["state"] == "enabled" {
			t.Browsers = append(t.Browsers, browserMap["name"].(string))
		}
	}
}

// marshallElements converts go rfml structs into format understood by the API
func (t *RFTest) marshallElements(finder TestIDFinder) error {
	// if there are no steps skip marshalling
	if len(t.Steps) == 0 {
		return nil
	}
	t.Elements = make([]testElement, len(t.Steps))
	for i, step := range t.Steps {
		switch castStep := step.(type) {
		case RFTestStep:
			stepElementDetails := testElementDetails{Action: castStep.Action, Response: castStep.Response}
			stepElement := testElement{Redirect: castStep.Redirect, Type: "step", Details: stepElementDetails}
			t.Elements[i] = stepElement
		case RFEmbeddedTest:
			embeddedID, err := finder.GetTestID(castStep.RFMLID)
			if err != nil {
				return err
			}
			embeddedElementDetails := testElementDetails{ID: embeddedID}
			embeddedElement := testElement{Redirect: castStep.Redirect, Type: "test", Details: embeddedElementDetails}
			t.Elements[i] = embeddedElement
		}
	}
	return nil
}

// unmarshalElements converts API elements format into RFML go structs
func (t *RFTest) unmarshalElements(finder TestIDFinder) error {
	if len(t.Elements) == 0 {
		return nil
	}
	t.Steps = make([]interface{}, len(t.Elements))
	for i, element := range t.Elements {
		switch element.Type {
		case "step":
			step := RFTestStep{Action: element.Details.Action, Response: element.Details.Response, Redirect: element.Redirect}
			t.Steps[i] = step
		case "test":
			rfmlID, err := finder.GetRFMLID(element.Details.ID)
			if err != nil {
				return err
			}
			embedd := RFEmbeddedTest{RFMLID: rfmlID, Redirect: element.Redirect}
			t.Steps[i] = embedd
		}
	}
	return nil
}

// PrepareToUploadFromRFML uses different helper methods to prepare struct for API upload
func (t *RFTest) PrepareToUploadFromRFML(finder TestIDFinder) error {
	t.Source = "rainforest-cli"
	if t.StartURI == "" {
		t.StartURI = "/"
	}

	// prevent []string(nil) as a value for Tags
	if len(t.Tags) == 0 {
		t.Tags = []string{}
	}

	t.mapBrowsers()
	err := t.marshallElements(finder)
	if err != nil {
		return err
	}
	return nil
}

// PrepareToWriteAsRFML uses different helper methods to prepare struct for translation to RFML
func (t *RFTest) PrepareToWriteAsRFML(finder TestIDFinder) error {
	err := t.unmarshalElements(finder)
	if err != nil {
		return err
	}
	t.unmapBrowsers()
	return nil
}

// HasUploadableFiles returns true if test has embedded files in the format {{ file.screenshot(path/to/file) }}
// or {{ file.download(path/to/file) }}. It returns false otherwise.
func (t *RFTest) HasUploadableFiles() bool {
	for _, step := range t.Steps {
		s, ok := step.(RFTestStep)
		if ok && s.hasUploadableFiles() {
			return true
		}
	}

	return false
}

// RFTestStep contains single Rainforest step
type RFTestStep struct {
	Action   string
	Response string
	Redirect bool
}

func (s *RFTestStep) hasUploadableFiles() bool {
	return len(s.embeddedFilesInAction()) > 0 || len(s.embeddedFilesInResponse()) > 0
}

func (s *RFTestStep) embeddedFilesInAction() []embeddedFile {
	return findEmbeddedFiles(s.Action)
}

func (s *RFTestStep) embeddedFilesInResponse() []embeddedFile {
	return findEmbeddedFiles(s.Response)
}

// uploadable contains the information of an embedded step variables
type embeddedFile struct {
	// text is the entire step variable text. eg: "{{ file.screenshot(path/to/file) }}"
	text string
	// the step variable used. Either "screenshot" or "download"
	stepVar string
	// the path argument to the step variable
	path string
}

// findEmbeddedFiles looks through a string and parses out embedded step variables
// and returns a slice of uploadables
func findEmbeddedFiles(s string) []embeddedFile {
	reg := regexp.MustCompile(uploadableRegex)
	matches := reg.FindAllStringSubmatch(s, -1)

	uploadables := []embeddedFile{}

	for _, match := range matches {
		// If there are multiple arguments, check that a file ID and a signature
		// exist. If so, ignore, as they are remote file references and can be
		// uploaded without any string replacements.
		if parameters := strings.Split(match[2], ","); len(parameters) > 1 {
			if _, err := strconv.Atoi(parameters[0]); err != nil {
				continue
			}

			if sig := parameters[1]; len(sig) != 6 {
				continue
			}
		}

		uploadable := embeddedFile{
			text:    match[0],
			stepVar: match[1],
			path:    match[2],
		}

		uploadables = append(uploadables, uploadable)
	}

	return uploadables
}

// RFEmbeddedTest contains an embedded test details
type RFEmbeddedTest struct {
	RFMLID   string
	Redirect bool
}

// RFTestFilters are used to translate test filters to a proper query string
type RFTestFilters struct {
	Tags          []string
	SiteID        int
	SmartFolderID int
	FeatureID     int
	RunGroupID    int
}

func (f *RFTestFilters) toQuery() string {
	// Empty slices are ignored by the Encode function, so no need to check for
	// presence here.
	v := url.Values{"tags": f.Tags}
	if f.SiteID > 0 {
		v.Add("site_id", strconv.Itoa(f.SiteID))
	}
	if f.SmartFolderID > 0 {
		v.Add("smart_folder_id", strconv.Itoa(f.SmartFolderID))
	}
	if f.FeatureID > 0 {
		v.Add("feature_id", strconv.Itoa(f.FeatureID))
	}
	if f.RunGroupID > 0 {
		v.Add("run_group_id", strconv.Itoa(f.RunGroupID))
	}

	return v.Encode()
}

// GetRFMLIDs returns all tests IDs and RFML IDs to properly map tests to their IDs
// for uploading and deleting.
func (c *Client) GetRFMLIDs() (*TestIDFinder, error) {
	// Prepare request
	req, err := c.NewRequest("GET", "tests/rfml_ids", nil)
	if err != nil {
		return nil, err
	}

	// Send request and process response
	var testResp []TestIDPair
	_, err = c.Do(req, &testResp)
	if err != nil {
		return nil, err
	}

	// FIXME: Gross
	finder := new(TestIDFinder)
	finder.Pairs = testResp

	return finder, nil
}

// GetTests returns all tests that are optionally filtered by RFTestFilters
func (c *Client) GetTests(params *RFTestFilters) ([]RFTest, error) {
	tests := []RFTest{}
	page := 1

	for {
		testsURL := "tests?page_size=50&page=" + strconv.Itoa(page)

		queryString := params.toQuery()
		if queryString != "" {
			testsURL = testsURL + "&" + queryString
		}

		req, err := c.NewRequest("GET", testsURL, nil)

		if err != nil {
			return nil, err
		}

		var testResp []RFTest
		_, err = c.Do(req, &testResp)
		if err != nil {
			return nil, err
		}
		for i := range testResp {
			testResp[i].Execute = true
		}

		tests = append(tests, testResp...)

		totalPagesHeader := c.LastResponseHeaders.Get("X-Total-Pages")
		if totalPagesHeader == "" {
			return nil, fmt.Errorf("Rainforest API error: Total pages header missing from response")
		}

		totalPages, err := strconv.Atoi(totalPagesHeader)
		if err != nil {
			return nil, err
		}

		if page == totalPages {
			return tests, nil
		}

		page++
	}
}

// GetTest gets a test from RF specified by the given test ID
func (c *Client) GetTest(testID int) (*RFTest, error) {
	req, err := c.NewRequest("GET", "tests/"+strconv.Itoa(testID), nil)
	if err != nil {
		return nil, err
	}

	var testResp RFTest
	_, err = c.Do(req, &testResp)
	if err != nil {
		return nil, err
	}

	testResp.TestID = testID
	testResp.Execute = true
	return &testResp, nil
}

// DeleteTest deletes test with a specified ID from the RF test suite
func (c *Client) DeleteTest(testID int) error {
	// Prepare request
	req, err := c.NewRequest("DELETE", "tests/"+strconv.Itoa(testID), nil)
	if err != nil {
		return err
	}

	// Send request and process response
	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTestByRFMLID deletes test with a specified RFMLID from the RF test suite
func (c *Client) DeleteTestByRFMLID(testRFMLID string) error {
	testIDFinder, err := c.GetRFMLIDs()
	if err != nil {
		return err
	}

	testID, err := testIDFinder.GetTestID(testRFMLID)
	if err != nil {
		return err
	}
	return c.DeleteTest(testID)
}

// CreateTest creates new test on RF, requires RFTest struct to be prepared to upload using helpers
func (c *Client) CreateTest(test *RFTest) error {
	// Prepare request
	req, err := c.NewRequest("POST", "tests", test)
	if err != nil {
		return err
	}

	// Send request and process response
	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTest updates existing test on RF, requires RFTest struct to be prepared to upload using helpers
func (c *Client) UpdateTest(test *RFTest) error {
	if test.TestID == 0 {
		return errors.New("Couldn't update the test TestID not specified in RFTest")
	}

	// Prepare request
	req, err := c.NewRequest("PUT", "tests/"+strconv.Itoa(test.TestID), test)
	if err != nil {
		return err
	}

	// Send request and process response
	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}
