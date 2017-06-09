package rainforest

import "strconv"

// RunParams is a struct holding all potential parameters needed to start a new RF run.
type RunParams struct {
	// This can be eiter []int or string containing 'all'
	Tests         interface{} `json:"tests,omitempty"`
	Tags          []string    `json:"tags,omitempty"`
	SmartFolderID int         `json:"smart_folder_id,omitempty"`
	SiteID        int         `json:"site_id,omitempty"`
	Crowd         string      `json:"crowd,omitempty"`
	Conflict      string      `json:"conflict,omitempty"`
	Browsers      []string    `json:"browsers,omitempty"`
	Description   string      `json:"description,omitempty"`
	EnvironmentID int         `json:"environment_id,omitempty"`
	RunGroupID 	  int 		  `json:"run_group_id,omitempty"`
}

// RunStatus represents a status of a RF run in progress.
type RunStatus struct {
	ID           int    `json:"id"`
	State        string `json:"state"`
	StateDetails struct {
		Name         string `json:"name"`
		IsFinalState bool   `json:"is_final_state"`
	} `json:"state_details"`
	Result          string `json:"result"`
	CurrentProgress struct {
		Percent  int `json:"percent"`
		Total    int `json:"total"`
		Complete int `json:"complete"`
		NoResult int `json:"no_result"`
	} `json:"current_progress"`
	FrontendURL string `json:"frontend_url,omitempty"`
}

// CreateRun starts a new RF run with given params.
func (c *Client) CreateRun(params RunParams) (*RunStatus, error) {
	var runStatus RunStatus

	// Usual stuff - create a request and send it
	req, err := c.NewRequest("POST", "runs", params)
	if err != nil {
		return &runStatus, err
	}
	_, err = c.Do(req, &runStatus)
	if err != nil {
		return &runStatus, err
	}

	return &runStatus, nil
}

// CheckRunStatus returns the status of a specified run.
func (c *Client) CheckRunStatus(runID int) (*RunStatus, error) {
	var runStatus RunStatus
	// Get proper URL then prepare and send the request
	url := "runs/" + strconv.Itoa(runID)

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return &runStatus, err
	}
	_, err = c.Do(req, &runStatus)
	if err != nil {
		return &runStatus, err
	}

	return &runStatus, nil
}
