package main

import (
	"encoding/json"
	"fmt"
)

type runParams struct {
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
	return &runParams{
		Tags:          tags,
		SmartFolderID: smartFolderID,
		SiteID:        siteID,
	}
}

func postRun(params *runParams) (resBody *runResponse) {
	js, err := json.Marshal(params)

	if err != nil {
		panic(fmt.Sprintf("Unable to format JSON for run. Params: %v", params))
	}

	data := postRequest("https://app.rainforestqa.com/api/1/runs", js)
	json.Unmarshal(data, &resBody)
	return
}
