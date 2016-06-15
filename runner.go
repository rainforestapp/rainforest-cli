package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/urfave/cli.v2"
)

type runParams struct {
	Tags []string `json:"tags"`
}

type runResponse map[string]interface{}

func createRun(c *cli.Context) {
	params := makeParams(c)
	js, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", "https://app.rainforestqa.com/api/1/runs", bytes.NewBuffer(js))

	addAuthHeaders(req)
	res, _ := client.Do(req)

	var resBody runResponse
	data, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(data, &resBody)

	fmt.Println(resBody)
}

func makeParams(c *cli.Context) *runParams {
	return &runParams{
		Tags: c.StringSlice("tags"),
	}
}
