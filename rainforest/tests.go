package rainforest

import (
	"errors"
	"fmt"
	"strconv"
)

// TestIDMap is a type representing RF tests that contain the test definitions.
type TestIDMap struct {
	ID     int    `json:"id"`
	RFMLID string `json:"rfml_id"`
}

// TestIDMappings is a slice of all the mapping pairs.
// And has a set of functions defined to get map of one to the other.
type TestIDMappings []TestIDMap

func (s TestIDMappings) mapIDtoRFMLID() map[int]string {
	resultMap := make(map[int]string)
	for _, mapping := range s {
		resultMap[mapping.ID] = mapping.RFMLID
	}
	return resultMap
}

func (s TestIDMappings) mapRFMLIDtoID() map[string]int {
	resultMap := make(map[string]int)
	for _, mapping := range s {
		resultMap[mapping.RFMLID] = mapping.ID
	}
	return resultMap
}

// RFTest is a struct representing the Rainforest Test with its settings and steps
type RFTest struct {
	RFMLID      string              `json:"rfml_id"`
	Title       string              `json:"title, omitempty"`
	StartURI    string              `json:"start_uri, omitempty"`
	SiteID      int                 `json:"site_id, omitempty"`
	Description string              `json:"description, omitempty"`
	Tags        []string            `json:"tags, omitempty"`
	BrowsersMap []map[string]string `json:"browsers, omitempty"`
	Source      string              `json:"source"`
	Elements    []testElement       `json:"elements, omitempty"`
	Browsers    []string
	Steps       []interface{}
}

// testElement is one of the helpers to construct the proper JSON test sturcture
type testElement struct {
	Redirect bool               `json:"redirection"`
	Type     string             `json:"type"`
	Details  testElementDetails `json:"element"`
}

// testElementDetails is one of the helpers to construct the proper JSON test sturcture
type testElementDetails struct {
	ID       int    `json:"id, omitempty"`
	Action   string `json:"action, omitempty"`
	Response string `json:"response, omitempty"`
}

func (t *RFTest) mapBrowsers() {
	for _, browser := range t.Browsers {
		mappedBrowser := map[string]string{
			"state": "enabled",
			"name":  browser,
		}
		t.BrowsersMap = append(t.BrowsersMap, mappedBrowser)
	}
}

func (t *RFTest) unmapBrowsers() {
	for _, browserMap := range t.BrowsersMap {
		if browserMap["state"] == "enabled" {
			t.Browsers = append(t.Browsers, browserMap["name"])
		}
	}
}

func (t *RFTest) marshallElements(mappings TestIDMappings) error {
	rfmlidToID := mappings.mapRFMLIDtoID()
	for _, step := range t.Steps {
		switch castStep := step.(type) {
		case RFTestStep:
			stepElementDetails := testElementDetails{Action: castStep.Action, Response: castStep.Response}
			stepElement := testElement{Redirect: castStep.Redirect, Type: "step", Details: stepElementDetails}
			t.Elements = append(t.Elements, stepElement)
		case RFEmbeddedTest:
			embeddedID, ok := rfmlidToID[castStep.RFMLID]
			if !ok {
				return errors.New("Couldn't convert RFML ID to test ID")
			}
			embeddedElementDetails := testElementDetails{ID: embeddedID}
			embeddedElement := testElement{Redirect: castStep.Redirect, Type: "test", Details: embeddedElementDetails}
			t.Elements = append(t.Elements, embeddedElement)
		}
	}
	return nil
}

func (t *RFTest) unmarshallElements(mappings TestIDMappings) error {
	idToRFMLID := mappings.mapIDtoRFMLID()
	for _, element := range t.Elements {
		switch element.Type {
		case "step":
			step := RFTestStep{Action: element.Details.Action, Response: element.Details.Response, Redirect: element.Redirect}
			t.Steps = append(t.Steps, step)
		case "test":
			rfmlID, ok := idToRFMLID[element.Details.ID]
			if !ok {
				return errors.New("Couldn't convert test ID to RFML ID")
			}
			embedd := RFEmbeddedTest{RFMLID: rfmlID, Redirect: element.Redirect}
			t.Steps = append(t.Steps, embedd)
		}
	}
	return nil
}

// RFTestStep contains single Rainforest step
type RFTestStep struct {
	Action   string
	Response string
	Redirect bool
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
	rfmlMap := testMappings.mapRFMLIDtoID()
	testID, ok := rfmlMap[testRFMLID]
	if !ok {
		return fmt.Errorf("RFML ID: %v doesn't exist in Rainforest", testRFMLID)
	}
	return c.DeleteTest(testID)
}

// CreateTest creates new test on RF
func (c *Client) CreateTest(test RFTest) error {
	// Get mappings needed to convert RFML IDs
	mappings, err := c.GetRFMLIDs()
	if err != nil {
		return err
	}
	// Prepare JSON struct
	test.Source = "rainforest-cli"
	test.mapBrowsers()
	err = test.marshallElements(mappings)
	if err != nil {
		return err
	}

	// Prepare request
	req, err := c.NewRequest("POST", "/tests", test)
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

// UpdateTest updates existing test on RF
func (c *Client) UpdateTest(test RFTest) error {
	// Get mappings needed to convert RFML IDs
	mappings, err := c.GetRFMLIDs()
	if err != nil {
		return err
	}

	testID, ok := mappings.mapRFMLIDtoID()[test.RFMLID]
	if !ok {
		return errors.New("Couldn't update the test - RFML ID doesn't exist on RF")
	}
	// Prepare JSON struct
	test.Source = "rainforest-cli"
	test.mapBrowsers()
	err = test.marshallElements(mappings)
	if err != nil {
		return err
	}

	// Prepare request
	req, err := c.NewRequest("PUT", "/tests/"+strconv.Itoa(testID), test)
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
