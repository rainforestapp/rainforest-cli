package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

type foldersParams struct {
	Tags []string `json:"tags"`
}
type foldersResponse map[string]interface{}

func getFolders(c *cli.Context) {
	params := makeBody(c)
	resBody := postFolders(params)

	fmt.Println(resBody)
}

func makeBody(c *cli.Context) *foldersParams {
	return &foldersParams{
		Tags: c.StringSlice("tags"),
	}
}

func postFolders(params *foldersParams) (resBody *foldersResponse) {
	js, _ := json.Marshal(params)
	data := postRequest("https://app.rainforestqa.com/api/1/runs", js)
	json.Unmarshal(data, &resBody)
	return
}
