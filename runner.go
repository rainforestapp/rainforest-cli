package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

type runParams struct {
	Tags []string `json:"tags"`
}

func createRun(c *cli.Context) {
	params := makeParams(c)
	js, _ := json.Marshal(params)
	fmt.Println(string(js))
}

func makeParams(c *cli.Context) *runParams {
	return &runParams{
		Tags: c.StringSlice("tags"),
	}
}
