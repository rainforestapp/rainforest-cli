package rainforest

import "strconv"

// RunStateDetails contains details about the state of a Run
type RunStateDetails struct {
	Name         string `json:"name"`
	IsFinalState bool   `json:"is_final_state"`
}

// RunDetails contains top level details of a Run
type RunDetails struct {
	Description        string            `json:"description"`
	TotalTests         int               `json:"total_tests"`
	TotalFailedTests   int               `json:"total_failed_tests"`
	TotalNoResultTests int               `json:"total_no_result_tests"`
	StateDetails       RunStateDetails   `json:"state_details"`
	Timestamps         map[string]string `json:"timestamps"`
}

// GetRunDetails returns the top level details of a Run
func (c *Client) GetRunDetails(runID int) (*RunDetails, error) {
	var runDetails RunDetails
	url := "runs/" + strconv.Itoa(runID)

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return &runDetails, err
	}

	_, err = c.Do(req, &runDetails)

	return &runDetails, err
}
