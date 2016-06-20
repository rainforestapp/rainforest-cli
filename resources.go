package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/urfave/cli.v2"
)

type clientResp struct {
	AvailableBrowsers []struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		Category       string `json:"category"`
		BrowserVersion string `json:"browser_version"`
		OsVersion      string `json:"os_version"`
		Default        bool   `json:"default"`
	} `json:"available_browsers"`
}

type resourceParams struct {
	Tags []string `json:"tags"`
}
type resourceResponse []map[string]interface{}
type test interface{}

func createResource(c *cli.Context, resourceType string) {
	params := makeBody(c)
	//var resBody *resourceResponse
	if resourceType == "Browsers" {
		getBrowsers(params)
		//printBrowsers(resBody)
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
	data := getRequest(url)
	json.Unmarshal(data, &resBody)
	return
}

func getBrowsers(params *resourceParams) (resBody *resourceResponse) {
	data := getRequest("https://app.rainforestqa.com/api/1/clients.json")
	var client clientResp
	json.Unmarshal(data, &client)
	for _, item := range client.AvailableBrowsers {
		fmt.Printf("\t%v\t| %v\n", item.Name, item.Description)
	}

	return
}

func printResource(resBody *resourceResponse, resourceType string) {
	resource := resourceType[0 : len(resourceType)-1]
	fmt.Printf("%v Id\t| %v Name\n", resource, resource)
	bar := strings.Repeat("-", 40)
	print("" + bar + "\n")
	for _, item := range *resBody {
		fmt.Printf("\t%v\t| %v\n", item["id"], item["title"])
	}
}
