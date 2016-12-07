package rainforest

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
	Title       string              `json:"title"`
	StartURI    string              `json:"start_uri"`
	SiteID      int                 `json:"site_id"`
	Description string              `json:"description"`
	Tags        []string            `json:"tags"`
	BrowsersMap []map[string]string `json:"browser_json"`
	Browsers    []string
	Steps       []interface{}
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
