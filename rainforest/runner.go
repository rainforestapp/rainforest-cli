package rainforest

import (
	"errors"
	"fmt"
	"strconv"
)

// RunParams is a struct holding all potential parameters needed to start a new RF run.
type RunParams struct {
	// This can be eiter []int or string containing 'all'
	Tests                interface{} `json:"tests,omitempty"`
	RFMLIDs              []string    `json:"rfml_ids,omitempty"`
	Tags                 []string    `json:"tags,omitempty"`
	SmartFolderID        int         `json:"smart_folder_id,omitempty"`
	SiteID               int         `json:"site_id,omitempty"`
	Crowd                string      `json:"crowd,omitempty"`
	Conflict             string      `json:"conflict,omitempty"`
	Browsers             []string    `json:"browsers,omitempty"`
	Description          string      `json:"description,omitempty"`
	Release              string      `json:"release,omitempty"`
	EnvironmentID        int         `json:"environment_id,omitempty"`
	FeatureID            int         `json:"feature_id,omitempty"`
	RunGroupID           int         `json:"-"`
	RunID                int         `json:"-"`
	AutomationMaxRetries int         `json:"automation_max_retries,omitempty"`
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
		Passed   int `json:"passed"`
		Failed   int `json:"failed"`
	} `json:"current_progress"`
	FrontendURL string `json:"frontend_url,omitempty"`
}

// CreateRun starts a new RF run with given params.
func (c *Client) CreateRun(params RunParams) (*RunStatus, error) {
	var runStatus RunStatus

	endpoint := "runs"
	if params.RunID > 0 {
		err := validateRerunParams(params)
		if err != nil {
			return &runStatus, err
		}
		endpoint = fmt.Sprintf("runs/%v/rerun_failed", params.RunID)
	} else if params.RunGroupID > 0 {
		err := validateRunGroupParams(params)
		if err != nil {
			return &runStatus, err
		}
		endpoint = fmt.Sprintf("run_groups/%v/runs", params.RunGroupID)
	}

	// Usual stuff - create a request and send it
	req, err := c.NewRequest("POST", endpoint, params)
	if err != nil {
		return &runStatus, err
	}
	_, err = c.Do(req, &runStatus)
	if err != nil {
		return &runStatus, err
	}

	return &runStatus, nil
}

func validateRerunParams(params RunParams) error {
	if params.Tests != nil {
		return errors.New("Tests cannot be specified for rerun")
	}
	if params.RFMLIDs != nil {
		return errors.New("RFML tests cannot be specified for rerun")
	}
	if params.Tags != nil {
		return errors.New("Tags cannot be specified for rerun")
	}
	if params.SmartFolderID != 0 {
		return errors.New("Folder cannot be specified for rerun")
	}
	if params.SiteID != 0 {
		return errors.New("Site cannot be specified for rerun")
	}
	if params.Crowd != "" {
		return errors.New("Crowd cannot be specified for rerun")
	}
	if params.Browsers != nil {
		return errors.New("Browsers cannot be specified for rerun")
	}
	if params.Description != "" {
		return errors.New("Description cannot be specified for rerun")
	}
	if params.Release != "" {
		return errors.New("Release cannot be specified for rerun")
	}
	if params.EnvironmentID != 0 {
		return errors.New("Environment cannot be specified for rerun")
	}
	if params.FeatureID != 0 {
		return errors.New("Feature cannot be specified for rerun")
	}
	if params.RunGroupID != 0 {
		return errors.New("Run Group cannot be specified for rerun")
	}

	return nil
}

func validateRunGroupParams(params RunParams) error {
	if params.Tests != nil {
		return errors.New("Tests cannot be specified alongside run group")
	}
	if params.Tags != nil {
		return errors.New("Tags cannot be specified alongside run group")
	}
	if params.SmartFolderID != 0 {
		return errors.New("Folder cannot be specified alongside run group")
	}
	if params.SiteID != 0 {
		return errors.New("Site cannot be specified alongside run group")
	}
	if params.Browsers != nil {
		return errors.New("Browsers cannot be specified alongside run group")
	}
	if params.FeatureID != 0 {
		return errors.New("Feature cannot be specified alongside run group")
	}

	return nil
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
