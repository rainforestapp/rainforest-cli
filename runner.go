package main

import (
	"encoding/json"
	"fmt"
)

type runParams struct {
	Tags          []string `json:"tags,omitempty"`
	SmartFolderID int      `json:"smart_folder_id,omitempty"`
}

type runResponse map[string]interface{}

type flagParser interface {
	StringSlice(string) []string
	Int(string) int
}

func createRun(f flagParser) {
	params := makeParams(f)
	resBody := postRun(params)

	fmt.Println(resBody)
}

func makeParams(c flagParser) *runParams {
	return &runParams{
		Tags:          c.StringSlice("tags"),
		SmartFolderID: c.Int("smart-folder-id"),
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
