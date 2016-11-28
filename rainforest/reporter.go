package rainforest

import "strconv"

// RunStateDetails contains details about the state of a Run
type RunStateDetails struct {
	Name         string `json:"name"`
	IsFinalState bool   `json:"is_final_state"`
}

// RunTestDetails contains details about a specific Run Test
type RunTestDetails struct {
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// RunDetails contains top level details of a Run
type RunDetails struct {
	Description        string            `json:"description"`
	TotalTests         int               `json:"total_tests"`
	TotalFailedTests   int               `json:"total_failed_tests"`
	TotalNoResultTests int               `json:"total_no_result_tests"`
	StateDetails       RunStateDetails   `json:"state_details"`
	Timestamps         map[string]string `json:"timestamps"`
	Tests              []RunTestDetails
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
	if err != nil {
		return &runDetails, err
	}

	// NOTE: This extra request is only necessary because `update_at` is not
	// currently exposed in the `/runs/:id` endpoint. This may change in the future:
	// https://github.com/rainforestapp/rainforest-cli/issues/216
	var runTests []RunTestDetails
	url = "runs/" + strconv.Itoa(runID) + "/tests?page_size=" + strconv.Itoa(runDetails.TotalTests)

	req, err = c.NewRequest("GET", url, nil)
	if err != nil {
		return &runDetails, err
	}

	_, err = c.Do(req, &runTests)
	if err != nil {
		return &runDetails, err
	}

	runDetails.Tests = runTests

	return &runDetails, err
}
