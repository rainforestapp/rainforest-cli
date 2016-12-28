package rainforest

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const uploadableRegex = `{{ *file\.(download|screenshot)\(([^\)]+)\) *}}`

// TestIDMap is a type representing RF tests that contain the test definitions.
type TestIDMap struct {
	ID     int    `json:"id"`
	RFMLID string `json:"rfml_id"`
}

// TestIDMappings is a slice of all the mapping pairs.
// And has a set of functions defined to get map of one to the other.
type TestIDMappings []TestIDMap

// MapIDtoRFMLID creates a map from test IDs to RFML IDs
func (s TestIDMappings) MapIDtoRFMLID() map[int]string {
	resultMap := make(map[int]string)
	for _, mapping := range s {
		resultMap[mapping.ID] = mapping.RFMLID
	}
	return resultMap
}

// MapRFMLIDtoID creates a map from RFML IDs to IDs
func (s TestIDMappings) MapRFMLIDtoID() map[string]int {
	resultMap := make(map[string]int)
	for _, mapping := range s {
		resultMap[mapping.RFMLID] = mapping.ID
	}
	return resultMap
}

// RFTest is a struct representing the Rainforest Test with its settings and steps
type RFTest struct {
	RFMLID      string              `json:"rfml_id"`
	Source      string              `json:"source"`
	Title       string              `json:"title,omitempty"`
	StartURI    string              `json:"start_uri,omitempty"`
	SiteID      int                 `json:"site_id,omitempty"`
	Description string              `json:"description,omitempty"`
	Tags        []string            `json:"tags,omitempty"`
	BrowsersMap []map[string]string `json:"browsers,omitempty"`
	Elements    []testElement       `json:"elements,omitempty"`

	// Browsers, Steps and TestID are helper fields
	Browsers []string      `json:"-"`
	Steps    []interface{} `json:"-"`
	TestID   int           `json:"-"`
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
	// if there are no browsers skip mapping
	if len(t.Browsers) == 0 {
		return
	}
	t.BrowsersMap = make([]map[string]string, len(t.Browsers))
	for i, browser := range t.Browsers {
		mappedBrowser := map[string]string{
			"state": "enabled",
			"name":  browser,
		}
		t.BrowsersMap[i] = mappedBrowser
	}
}

// unmapBrowsers parses browsers from the API format to internal go one
func (t *RFTest) unmapBrowsers() {
	// if there are no browsers skip unmapping
	if len(t.BrowsersMap) == 0 {
		return
	}

	for _, browserMap := range t.BrowsersMap {
		if browserMap["state"] == "enabled" {
			t.Browsers = append(t.Browsers, browserMap["name"])
		}
	}
}

// marshallElements converts go rfml structs into format understood by the API
func (t *RFTest) marshallElements(mappings TestIDMappings) error {
	// if there are no steps skip marshalling
	if len(t.Steps) == 0 {
		return nil
	}
	t.Elements = make([]testElement, len(t.Steps))
	rfmlidToID := mappings.MapRFMLIDtoID()
	for i, step := range t.Steps {
		switch castStep := step.(type) {
		case RFTestStep:
			stepElementDetails := testElementDetails{Action: castStep.Action, Response: castStep.Response}
			stepElement := testElement{Redirect: castStep.Redirect, Type: "step", Details: stepElementDetails}
			t.Elements[i] = stepElement
		case RFEmbeddedTest:
			embeddedID, ok := rfmlidToID[castStep.RFMLID]
			if !ok {
				return errors.New("Couldn't convert RFML ID to test ID")
			}
			embeddedElementDetails := testElementDetails{ID: embeddedID}
			embeddedElement := testElement{Redirect: castStep.Redirect, Type: "test", Details: embeddedElementDetails}
			t.Elements[i] = embeddedElement
		}
	}
	return nil
}

// unmarshalElements converts API elements format into RFML go structs
func (t *RFTest) unmarshalElements(mappings TestIDMappings) error {
	if len(t.Elements) == 0 {
		return nil
	}
	t.Steps = make([]interface{}, len(t.Elements))
	idToRFMLID := mappings.MapIDtoRFMLID()

	for i, element := range t.Elements {
		switch element.Type {
		case "step":
			step := RFTestStep{Action: element.Details.Action, Response: element.Details.Response, Redirect: element.Redirect}
			t.Steps[i] = step
		case "test":
			rfmlID, ok := idToRFMLID[element.Details.ID]
			if !ok {
				return errors.New("Couldn't convert test ID to RFML ID")
			}
			embedd := RFEmbeddedTest{RFMLID: rfmlID, Redirect: element.Redirect}
			t.Steps[i] = embedd
		}
	}
	return nil
}

// PrepareToUploadFromRFML uses different helper methods to prepare struct for API upload
func (t *RFTest) PrepareToUploadFromRFML(mappings TestIDMappings) error {
	t.Source = "rainforest-cli"
	if t.StartURI == "" {
		t.StartURI = "/"
	}
	t.mapBrowsers()
	err := t.marshallElements(mappings)
	if err != nil {
		return err
	}
	testID, ok := mappings.MapRFMLIDtoID()[t.RFMLID]
	if ok {
		t.TestID = testID
	}
	return nil
}

// PrepareToWriteAsRFML uses different helper methods to prepare struct for translation to RFML
func (t *RFTest) PrepareToWriteAsRFML(mappings TestIDMappings) error {
	err := t.unmarshalElements(mappings)
	if err != nil {
		return err
	}
	t.unmapBrowsers()
	return nil
}

func (t *RFTest) hasUploadableFiles() bool {
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
	return len(s.uploadablesInAction()) > 0 || len(s.uploadablesInResponse()) > 0
}

func (s *RFTestStep) uploadablesInAction() [][]string {
	return findUploadables(s.Action)
}

func (s *RFTestStep) uploadablesInResponse() [][]string {
	return findUploadables(s.Response)
}

func findUploadables(s string) [][]string {
	// Shouldn't fail compilation unless uploadableRegex is incorrect
	reg := regexp.MustCompile(uploadableRegex)
	return reg.FindAllStringSubmatch(s, -1)
}

// RFEmbeddedTest contains an embedded test details
type RFEmbeddedTest struct {
	RFMLID   string
	Redirect bool
}

// GetRFMLIDs returns all tests IDs and RFML IDs to properly map tests to their IDs
// for uploading and deleting.
func (c *Client) GetRFMLIDs() (TestIDMappings, error) {
	// Prepare request
	req, err := c.NewRequest("GET", "tests/rfml_ids", nil)
	if err != nil {
		return nil, err
	}

	// Send request and process response
	var testResp TestIDMappings
	_, err = c.Do(req, &testResp)
	if err != nil {
		return nil, err
	}
	return testResp, nil
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
	testMappings, err := c.GetRFMLIDs()
	if err != nil {
		return err
	}
	rfmlMap := testMappings.MapRFMLIDtoID()
	testID, ok := rfmlMap[testRFMLID]
	if !ok {
		return fmt.Errorf("RFML ID: %v doesn't exist in Rainforest", testRFMLID)
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

func isFileUploaded(data []byte, uploadedFiles []UploadedFile) bool {
	sum := md5.Sum(data)
	digest := string(sum[:])
	for _, uploadedFile := range uploadedFiles {
		if uploadedFile.Digest == digest {
			return true
		}
	}
	return false
}

// UpdateTest updates existing test on RF, requires RFTest struct to be prepared to upload using helpers
func (c *Client) UpdateTest(test *RFTest) error {
	if test.TestID == 0 {
		return errors.New("Couldn't update the test TestID not specified in RFTest")
	}

	if test.hasUploadableFiles() {
		uploadedFiles, err := c.getUploadedFiles(test.TestID)
		if err != nil {
			return err
		}

		digestToFileMap := map[string]UploadedFile{}
		for _, uploadedFile := range uploadedFiles {
			digestToFileMap[uploadedFile.Digest] = uploadedFile
		}

		for _, step := range test.Steps {
			s, ok := step.(RFTestStep)
			if ok && s.hasUploadableFiles() {
				if matches := s.uploadablesInAction(); len(matches) > 0 {
					for _, match := range matches {
						var filePath string
						filePath, err = filepath.Abs(match[2])
						if err != nil {
							return err
						}

						var file *os.File
						file, err = os.Open(filePath)
						if err != nil {
							return err
						}

						data, err := ioutil.ReadAll(file)
						if err != nil {
							return err
						}

						checksum := md5.Sum(data)
						fileDigest := string(checksum[:])
						uploadedFile, ok := digestToFileMap[fileDigest]
						if !ok {
							// File has not been uploaded before
							// Upload to RF
							awsFileInfo, err := c.createTestFile(test.TestID, file, data)
							if err != nil {
								return err
							}
							// Upload to AWS
							err = c.uploadTestFile(filepath.Base(filePath), data, awsFileInfo)
							if err != nil {
								return err
							}
							uploadedFile = UploadedFile{
								ID:        awsFileInfo.FileID,
								Signature: awsFileInfo.FileSignature,
								Digest:    fileDigest,
							}
							// Add to the mappings for future reference
							digestToFileMap[fileDigest] = uploadedFile
						}

						sig := uploadedFile.Digest[0:6]
						var replacement string
						stepVar := match[1]
						if stepVar == "screenshot" {
							replacement = fmt.Sprintf("{{ file.screenshot(%v, %v) }}", uploadedFile.ID, sig)
						} else if stepVar == "download" {
							replacement = fmt.Sprintf("{{ file.download(%v, %v, %v) }}", uploadedFile.ID, sig, filepath.Base(filePath))
						}

						s.Action = strings.Replace(s.Action, match[0], replacement, 1)
					}
				}
			}
		}
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
