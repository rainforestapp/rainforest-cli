package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
}

func createRun() {
	params := makeParams()
	resBody := postRun(params)

	fmt.Println(resBody)
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

func postRun(params *runParams) (resBody *runResponse) {
	js, err := json.Marshal(params)

	if err != nil {
		panic(fmt.Sprintf("Unable to format JSON for run. Params: %v", params))
	}

	data := postRequest(baseURL+"/runs", js)
	json.Unmarshal(data, &resBody)
	// runID := string(resBody.ID)
	// if fg == "" {
	// 	checkRunProgress(runID)
	// }
	return
}

// func checkRunProgress(runID string) {
// 	running := true
// 	for running {
// 		var response runResponse
// 		data := getRun(runID, response)
//
//
// 		if !response.StateDetails.IsFinalState {
//
// 		}
//
// 	}
//
// }
