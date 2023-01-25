package rainforest

import (
	"strconv"
)

const NO_BRANCH = 0

// Branch is a struct representing a Rainforest Branch
type Branch struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

// GetBranches returns all branches optionally filtered by name
func (c *Client) GetBranches(params ...string) ([]Branch, error) {
	if len(params) > 0 {
		params = []string{"name=" + params[0]}
	}

	var branches []Branch
	collect := func(coll interface{}) {
		newBranches := coll.(*[]Branch)
		for _, branch := range *newBranches {
			branches = append(branches, branch)
		}
	}

	err := c.getPaginatedResource("branches", &[]Branch{}, collect, params...)
	return branches, err
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

// MergeBranch merges an existing branch into the client's main branch
func (c *Client) MergeBranch(branchID int) error {
	// Prepare request
	req, err := c.NewRequest("PUT", "branches/"+strconv.Itoa(branchID)+"/merge", nil)
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

// DeleteBranch deletes an existing branch on RF specified by ID
func (c *Client) DeleteBranch(branchID int) error {
	// Prepare request
	req, err := c.NewRequest("DELETE", "branches/"+strconv.Itoa(branchID), nil)
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
