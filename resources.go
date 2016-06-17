package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

type resourceParams struct {
	Tags []string `json:"tags"`
}
type resourceResponse []map[string]interface{}
type test interface{}

func createResource(c *cli.Context, resourceType string) {
	params := makeBody(c)
	//var resBody *resourceResponse
	if resourceType == "browsers" {
		getBrowsers(params)
		//printBrowsers(resBody)
	} else {
		resBody := getResource(params, resourceType)
		printResource(resBody)
	}
}

func makeBody(c *cli.Context) *resourceParams {
	return &resourceParams{
		Tags: c.StringSlice("tags"),
	}
}

func getResource(params *resourceParams, resourceType string) (resBody *resourceResponse) {
	//js, _ := json.Marshal(params)
	url := "https://app.rainforestqa.com/api/1/" + resourceType + ".json"
	data := getRequest(url)
	json.Unmarshal(data, &resBody)
	return
}

func getBrowsers(params *resourceParams) (resBody *resourceResponse) {
	data := getRequest("https://app.rainforestqa.com/api/1/clients.json")
	type Client map[string]interface{}
	var client Client
	json.Unmarshal(data, &client)
	//resBrowsers := (*resBody)
	fmt.Printf("\n\n\n%T\n\n\n", client["available_browsers"])
	return
}

func printResource(resBody *resourceResponse) {
	// fmt.Printf("%v\n\n", *resBody)
	for _, item := range *resBody {
		fmt.Printf("\n%v\n\n", item)
	}
}
