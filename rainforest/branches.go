package rainforest

// Branch is a struct representing a Rainforest Branch
type Branch struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

// CreateBranch creates new branch on RF, requires Branch struct to be prepared to upload using helpers
func (c *Client) CreateBranch(branch *Branch) error {
	// Prepare request
	req, err := c.NewRequest("POST", "branches", branch)

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
