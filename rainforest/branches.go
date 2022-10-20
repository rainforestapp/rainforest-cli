package rainforest

import (
	"fmt"
	"strconv"
)

// Branch is a struct representing a Rainforest Branch
type Branch struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

// GetBranches returns all branches optionally filtered by name
func (c *Client) GetBranches(params ...string) ([]Branch, error) {
	name := ""
	branches := []Branch{}
	page := 1

	if len(params) > 0 {
		name = params[0]
	}

	for {
		branchesURL := "branches?page_size=50&page=" + strconv.Itoa(page)

		if name != "" {
			branchesURL = branchesURL + "&name=" + name
		}

		req, err := c.NewRequest("GET", branchesURL, nil)

		if err != nil {
			return nil, err
		}

		var branchResp []Branch
		_, err = c.Do(req, &branchResp)
		if err != nil {
			return nil, err
		}

		branches = append(branches, branchResp...)

		totalPagesHeader := c.LastResponseHeaders.Get("X-Total-Pages")
		if totalPagesHeader == "" {
			return nil, fmt.Errorf("Rainforest API error: Total pages header missing from response")
		}

		totalPages, err := strconv.Atoi(totalPagesHeader)
		if err != nil {
			return nil, err
		}

		if page == totalPages {
			return branches, nil
		}

		page++
	}
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
