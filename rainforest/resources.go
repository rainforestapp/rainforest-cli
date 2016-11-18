package rainforest

import "strconv"

// Resource interface is a standardized way of looking at the types returned from the resources
// list API. It treats every resource as a pair of ID and description of entity under this ID.
type Resource interface {
	GetID() string
	GetDescription() string
}

// Folder type represents a single folder returned by the API call for a list of folders
type Folder struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// GetID returns the Folder ID
func (f Folder) GetID() string {
	return strconv.Itoa(f.ID)
}

// GetDescription returns the Folder name
func (f Folder) GetDescription() string {
	return f.Title
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
	_, err = c.Do(req, &folderResp)
	if err != nil {
		return nil, err
	}
	return folderResp, nil
}

// Browser type represents a single browser returned by the API call for a list of browsers
type Browser struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetID returns the Browser name
func (b Browser) GetID() string {
	return b.Name
}

// GetDescription returns the Browser description
func (b Browser) GetDescription() string {
	return b.Description
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
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetID returns the Site ID
func (s Site) GetID() string {
	return strconv.Itoa(s.ID)
}

// GetDescription returns the Site name
func (s Site) GetDescription() string {
	return s.Name
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
