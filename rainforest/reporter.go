package rainforest

import (
	"strconv"
	"time"
)

// RunFeedback contains details about the feedback of a Run Step for a browser
type RunFeedback struct {
	AnswerGiven string `json:"answer_given"`
	JobState    string `json:"job_state"`
	Note        string `json:"note"`
}

// RunBrowserDetails contains details about a Browser of a Run Step
type RunBrowserDetails struct {
	Name     string        `json:"name"`
	Feedback []RunFeedback `json:"feedback"`
}

// RunStepDetails contains details about a Run Step
type RunStepDetails struct {
	Browsers []RunBrowserDetails `json:"browsers"`
}

// RunTestDetails contains details about a Run Test
type RunTestDetails struct {
	ID        int              `json:"id"`
	Title     string           `json:"title"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Result    string           `json:"result"`
	Steps     []RunStepDetails `json:"steps"`
}

// RunStateDetails contains details about the state of a Run
type RunStateDetails struct {
	Name         string `json:"name"`
	IsFinalState bool   `json:"is_final_state"`
}

// RunDetails contains top level details of a Run
type RunDetails struct {
	ID                 int                  `json:"id"`
	Description        string               `json:"description"`
	TotalTests         int                  `json:"total_tests"`
	TotalFailedTests   int                  `json:"total_failed_tests"`
	TotalNoResultTests int                  `json:"total_no_result_tests"`
	StateDetails       RunStateDetails      `json:"state_details"`
	Timestamps         map[string]time.Time `json:"timestamps"`
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

	// NOTE: This extra request is only necessary because `updated_at` is not
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

// GetRunTestDetails returns the detailed information for a RunTest
func (c *Client) GetRunTestDetails(runID int, testID int) (*RunTestDetails, error) {
	var runTestDetails RunTestDetails
	url := "runs/" + strconv.Itoa(runID) + "/tests/" + strconv.Itoa(testID)

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return &runTestDetails, err
	}

	_, err = c.Do(req, &runTestDetails)
	if err != nil {
		return &runTestDetails, err
	}

	return &runTestDetails, nil
}
