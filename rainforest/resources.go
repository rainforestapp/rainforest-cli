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

// getPaginatedResource gets all of a resource from endpoint (using multiple
// requests, if necessary) and returns the result.
//
// Because of the type system, usage is hairy. coll should be a pointer to an
// empty slice of the appropriate type, and collect is called every time
// resources are added to the collection. The caller should handle collecting
// the collection.
func (c *Client) getPaginatedResource(endpoint string, coll interface{}, collect func(interface{})) error {
	req, err := c.NewRequest("GET", endpoint+"?page_size=100", nil)
	if err != nil {
		return err
	}

	var res *http.Response
	res, err = c.Do(req, &coll)
	collect(coll)
	if err != nil {
		return err
	}

	totalPagesHeader := res.Header.Get("X-Total-Pages")
	if totalPagesHeader == "" {
		return fmt.Errorf("Missing X-Total-Pages header in HTTP Response")
	}

	totalPages, err := strconv.Atoi(totalPagesHeader)
	if err != nil {
		return err
	}

	for i := 1; i < totalPages; i++ {
		page := strconv.Itoa(i + 1)
		req, err = c.NewRequest("GET", endpoint+"?page_size=100&page="+page, nil)
		if err != nil {
			return err
		}

		res, err = c.Do(req, &coll)
		collect(coll)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetFolders returns a slice of Folders (their names and IDs) which are available
// for filtering RF tests.
func (c *Client) GetFolders() ([]Folder, error) {
	var folders []Folder

	collect := func(coll interface{}) {
		newFolders := coll.(*[]Folder)
		for _, folder := range *newFolders {
			folders = append(folders, folder)
		}
	}

	err := c.getPaginatedResource("folders", &[]Folder{}, collect)
	return folders, err
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

// Feature represents a single feature returned by the API.
type Feature struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// GetFeatures fetches available features.
func (c *Client) GetFeatures() ([]Feature, error) {
	var features []Feature
	collect := func(coll interface{}) {
		newFeatures := coll.(*[]Feature)
		for _, f := range *newFeatures {
			features = append(features, f)
		}
	}

	err := c.getPaginatedResource("features", &[]Feature{}, collect)
	return features, err
}

// RunGroup represents a single run group returned by the API.
type RunGroup struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// GetRunGroups fetches available run groups.
func (c *Client) GetRunGroups() ([]RunGroup, error) {
	var runGroups []RunGroup
	collect := func(coll interface{}) {
		newRunGroups := coll.(*[]RunGroup)
		for _, r := range *newRunGroups {
			runGroups = append(runGroups, r)
		}
	}

	err := c.getPaginatedResource("run_groups", &[]RunGroup{}, collect)
	return runGroups, err
}
