package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

type runParams struct {
	Tags []string `json:"tags,omitempty"`
}

type runResponse map[string]interface{}

func createRun(c *cli.Context) {
	params := makeParams(c)
	resBody := postRun(params)

	fmt.Println(resBody)
}

type flagParser interface {
	StringSlice(string) []string
}

func makeParams(c flagParser) *runParams {
	return &runParams{
		Tags: c.StringSlice("tags"),
	}
}

func postRun(params *runParams) (resBody *runResponse) {
	js, _ := json.Marshal(params)
	data := postRequest("https://app.rainforestqa.com/api/1/runs", js)
	json.Unmarshal(data, &resBody)
	return
}
