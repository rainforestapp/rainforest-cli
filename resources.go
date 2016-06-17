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
	var resBody *resourceResponse
	switch resourceType {
	case "folders":
		resBody = getFolders(params)
	case "sites":
		resBody = getSites(params)
	// case "browsers":
	// 	resBodi = getBrowsers(params)
	default:
		resBody = getFolders(params)
	}
	printTable(resBody)
}

func makeBody(c *cli.Context) *resourceParams {
	return &resourceParams{
		Tags: c.StringSlice("tags"),
	}
}

func getFolders(params *resourceParams) (resBody *resourceResponse) {
	//js, _ := json.Marshal(params)
	data := getRequest("https://app.rainforestqa.com/api/1/folders.json")
	json.Unmarshal(data, &resBody)
	return
}

func getSites(params *resourceParams) (resBody *resourceResponse) {
	//js, _ := json.Marshal(params)
	data := getRequest("https://app.rainforestqa.com/api/1/sites.json")
	json.Unmarshal(data, &resBody)
	return
}

// func getBrowsers(params *resourceParams) (resBody *resourceResponse) {
// 	data := getRequest("https://app.rainforestqa.com/api/1/clients.json")
// 	type Client map[string]interface{}
// 	var client Client
// 	json.Unmarshal(data, &client)
// 	//resBrowsers := (*resBody)
// 	fmt.Printf("\n\n\n\n\n%T\n\n\n\n\n\n\n", client["available_browsers"])
// 	resBody = resourceResponse(client["available_browsers"])
// 	return
// }

func printTable(resBody *resourceResponse) {
	// fmt.Printf("%v\n\n", *resBody)
	for _, item := range *resBody {
		fmt.Printf("\n%v\n\n", item)
	}
}
