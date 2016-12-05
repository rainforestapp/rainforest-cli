package rainforest

// Test is a type representing RF tests that contain the test definitions.
type Test struct {
	ID     int    `json:"id"`
	RFMLID string `json:"rfml_id"`
}

// GetRFMLIDs returns all tests IDs and RFML IDs to properly map tests to their IDs
// for uploading and deleting.
func (c *Client) GetRFMLIDs() ([]Test, error) {
	// Prepare request
	req, err := c.NewRequest("GET", "tests/rfml_ids", nil)
	if err != nil {
		return nil, err
	}

	// Send request and process response
	var testResp []Test
	_, err = c.Do(req, &testResp)
	if err != nil {
		return nil, err
	}
	return testResp, nil
}
