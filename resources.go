package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/urfave/cli.v2"
)

var get getResponse

type resourceParams struct {
	Tags []string `json:"tags"`
}
type resourceResponse []map[string]interface{}
type test interface{}

func createResource(c *cli.Context, resourceType string) {
	params := makeBody(c)
	if resourceType == "browsers" {
		getBrowsers(params)
	} else {
		resBody := getResource(params, resourceType)
		printResource(resBody, resourceType)
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
	data := get.getRequest(url)
	json.Unmarshal(data, &resBody)
	return
}

func getBrowsers(params *resourceParams) (resBody *resourceResponse) {
	data := get.getRequest("https://app.rainforestqa.com/api/1/clients.json")
	var client browsersResp
	json.Unmarshal(data, &client)
	for _, item := range client.AvailableBrowsers {
		fmt.Printf("\t%v\t| %v\n", item.Name, item.Description)
	}

	return
}

func printResource(resBody *resourceResponse, resourceType string) {
	resource := resourceType[0 : len(resourceType)-1]
	fmt.Printf("%vwhy Id\t| %v Name\n", resource, resource)
	bar := strings.Repeat("-", 40)
	print("" + bar + "\n")
	for _, item := range *resBody {
		fmt.Printf("\t%v\t| %v\n", item["id"], item["title"])
	}
}
