package rainforest

import "strconv"

// RunDetails contains top level details of a Run
type RunDetails struct {
	Description        string `json:"description"`
	TotalTests         int    `json:"total_tests"`
	TotalFailedTests   int    `json:"total_failed_tests"`
	TotalNoResultTests int    `json:"total_no_result_tests"`
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
