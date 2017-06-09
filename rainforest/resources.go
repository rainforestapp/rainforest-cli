package rainforest

import (
	"fmt"
	"net/http"
	"strconv"
)

// Folder type represents a single folder returned by the API call for a list of folders
type Folder struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// GetFolders returns a slice of Folders (their names and IDs) which are available
// for filtering RF tests.
func (c *Client) GetFolders() ([]Folder, error) {
	// Prepare request
	req, err := c.NewRequest("GET", "folders?page_size=100", nil)
	if err != nil {
		return nil, err
	}

	// Send request and process response
	var folderResp []Folder
	var res *http.Response
	res, err = c.Do(req, &folderResp)
	if err != nil {
		return nil, err
	}

	totalPagesHeader := res.Header.Get("X-Total-Pages")
	if totalPagesHeader == "" {
		return nil, fmt.Errorf("Missing X-Total-Pages header in HTTP Response")
	}

	totalPages, err := strconv.Atoi(totalPagesHeader)
	if err != nil {
		return nil, err
	}

	for i := 1; i < totalPages; i++ {
		page := strconv.Itoa(i + 1)
		req, err = c.NewRequest("GET", "folders?page_size=100&page="+page, nil)
		if err != nil {
			return nil, err
		}

		var extraFoldersResp []Folder
		res, err = c.Do(req, &extraFoldersResp)
		if err != nil {
			return nil, err
		}

		folderResp = append(folderResp, extraFoldersResp...)
	}

	return folderResp, nil
}

// Browser type represents a single browser returned by the API call for a list of browsers
type Browser struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetBrowsers returns a slice of Browsers which are available for the client to run
// RF tests against.
func (c *Client) GetBrowsers() ([]Browser, error) {
	// Prepare request
	req, err := c.NewRequest("GET", "clients", nil)
	if err != nil {
		return nil, err
	}

	// browsersResp is a type returned by the API call for list of browsers
	type browsersResponse struct {
		AvailableBrowsers []Browser `json:"available_browsers"`
	}
	var browsersResp browsersResponse
	// Send request and process response
	_, err = c.Do(req, &browsersResp)
	if err != nil {
		return nil, err
	}
	return browsersResp.AvailableBrowsers, nil
}

type RunGroup struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (c *Client) GetRunGroups() ([]RunGroup, error) {
	// Prepare request
	req, err := c.NewRequest("GET", "run_groups", nil)
	if err != nil {
		return nil, err
	}

	// Send request and process response
	var runGroupResp []RunGroup
	_, err = c.Do(req, &runGroupResp)
	if err != nil {
		return nil, err
	}
	return runGroupResp, nil
}

// Site type represents a single site returned by the API call for a list of sites
type Site struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

// GetSites fetches sites available to use during the RF runs.
func (c *Client) GetSites() ([]Site, error) {
	// Prepare request
	req, err := c.NewRequest("GET", "sites", nil)
	if err != nil {
		return nil, err
	}

	// Send request and process response
	var sitesResp []Site
	_, err = c.Do(req, &sitesResp)
	if err != nil {
		return nil, err
	}
	return sitesResp, nil
}
