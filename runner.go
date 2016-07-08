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
}

type runResponse map[string]interface{}

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
	}
	return &runParams{
		Tags:          slicedTags,
		SmartFolderID: smartFolderID,
		SiteID:        siteID,
	}
}

func stringToIntSlice(s string) []int {
	s = strings.TrimSpace(s)
	slicedString := strings.Split(s, ",")
	var slicedInt []int
	for _, slice := range slicedString {
		newInt, err := strconv.Atoi(slice)
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
	return
}
