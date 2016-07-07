package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type runParams struct {
	Tests         string   `json:"tests,omitempty"`
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
		testIDs = strings.TrimSpace(testIDs)
		return &runParams{Tests: testIDs}
	}
	tags = strings.TrimSpace(tags)
	slicedTags := strings.Split(tags, ",")
	return &runParams{
		Tags:          slicedTags,
		SmartFolderID: smartFolderID,
		SiteID:        siteID,
	}
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
