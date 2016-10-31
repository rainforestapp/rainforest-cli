package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type runParams struct {
	Tests         []int    `json:"tests,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	SmartFolderID int      `json:"smart_folder_id,omitempty"`
	SiteID        int      `json:"site_id,omitempty"`
	Crowd         string   `json:"crowd,omitempty"`
	Conflict      string   `json:"conflict,omitempty"`
	Browsers      string   `json:"browsers,omitempty"`
	Description   string   `json:"description,omitempty"`
	EnvironmentID int      `json:"environment_id,omitempty"`
}

type runResponse struct {
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

func createRun() {
	params := makeParams()
	response := postRun(params)
	fmt.Println(string(response))
}

func makeParams() *runParams {
	if testIDs != "" {
		testIDsSlice := stringToIntSlice(testIDs)
		return &runParams{Tests: testIDsSlice}
	}
	tags = strings.TrimSpace(tags)
	var slicedTags []string
	if tags != "" {
		slicedTags = strings.Split(tags, ",")
		for i, slice := range slicedTags {
			slicedTags[i] = strings.TrimSpace(slice)
		}
	}
	return &runParams{
		Tags:          slicedTags,
		SmartFolderID: smartFolderID,
		SiteID:        siteID,
		Crowd:         crowd,
	}
}

func stringToIntSlice(s string) []int {
	slicedString := strings.Split(s, ",")
	var slicedInt []int
	for _, slice := range slicedString {
		newInt, err := strconv.Atoi(strings.TrimSpace(slice))
		if err != nil {
			panic(err)
		}
		slicedInt = append(slicedInt, newInt)
	}
	return slicedInt
}

func postRun(params *runParams) []byte {
	var resBody runResponse
	js, err := json.Marshal(params)

	if err != nil {
		panic(fmt.Sprintf("Unable to format JSON for run. Params: %v", params))
	}

	data := postRequest(baseURL+"/runs", js)
	json.Unmarshal(data, &resBody)
	if !runTestInBackground && resBody.ID != 0 {
		runID := resBody.ID
		checkRunProgress(runID)
	}
	return data
}

func checkRunProgress(runID int) {
	running := true
	var response runResponse
	for running {

		getRun(strconv.Itoa(runID), &response)

		isFinalState := response.StateDetails.IsFinalState
		state := response.State
		currentPercent := response.CurrentProgress.Percent

		if !isFinalState {
			fmt.Printf("Run %v is %v and is %v%% complete\n", runID, state, currentPercent)
			if response.Result == "failed" && failFast {
				running = false
			}
		} else {
			fmt.Printf("Run %v is now %v and has %v\n", runID, state, response.Result)
			running = false
		}
		time.Sleep(waitTime)
	}
	if response.FrontendURL != "" {
		fmt.Printf("The detailed results are available at %v\n", response.FrontendURL)
	}

	if response.Result != "passed" {
		os.Exit(1)
	}

}
