package rainforest

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// RunGroupDetails shows the details for a particular run group.
type RunGroupDetails struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Environment struct {
		Name string `json:"name"`
	} `json:"environment"`
	Crowd      string `json:"crowd"`
	RerouteGeo string `json:"reroute_geo"`
	Schedule   struct {
		RepeatRules []struct {
			Day  string `json:"day"`
			Time string `json:"time"`
		} `json:"repeat_rules"`
	} `json:"schedule"`
}

// Print prints out details for a particular run group.
func (rgd *RunGroupDetails) Print() {
	fmt.Printf(`Details for Run Group #%v:
Name: %v
Environment: %v
Tester Crowd: %v
Location: %v
`,
		rgd.ID, rgd.Title, rgd.Environment.Name, rgd.Crowd, rgd.RerouteGeo)
	sched := rgd.Schedule

	if daysQuantity := len(sched.RepeatRules); daysQuantity > 0 {
		scheduledTime := sched.RepeatRules[0].Time

		days := make([]string, daysQuantity)
		for i, sr := range sched.RepeatRules {
			days[i] = sr.Day
		}

		daysStr := strings.Join(days, ", ")
		fmt.Printf("Schedule: %v @ %v\n", daysStr, scheduledTime)
	}
}

// GetRunGroupDetails gets details for a run group from the API.
func (c *Client) GetRunGroupDetails(runGroupID int) (*RunGroupDetails, error) {
	req, err := c.NewRequest("GET", "run_groups/"+strconv.Itoa(runGroupID), nil)
	if err != nil {
		return nil, err
	}

	var details RunGroupDetails
	_, err = c.Do(req, &details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

// GetRunJunit gets a run JUnit from the API.
func (c *Client) GetRunJunit(runID int) (*string, error) {
	req, err := c.NewRequest("GET", "runs/"+strconv.Itoa(runID)+"/junit.xml", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req, nil)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()

	return &newStr, nil
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

// GetEnvironments fetches environments available to use during the RF runs.
func (c *Client) GetEnvironments() ([]Environment, error) {
	// Prepare request
	req, err := c.NewRequest("GET", "environments", nil)
	if err != nil {
		return nil, err
	}

	// Send request and process response
	var envsResp []Environment
	_, err = c.Do(req, &envsResp)
	if err != nil {
		return nil, err
	}
	return envsResp, nil
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
