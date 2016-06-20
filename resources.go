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

func createResource(c *cli.Context, resourceType string) {
	params := makeBody(c)
	var IDs, Titles []string
	switch resourceType {
	case "Folders":
		IDs, Titles = getFolders(params)
	case "Sites":
		IDs, Titles = getSites(params)
	default:
		IDs, Titles = getBrowsers(params)
	}
	printResource(resourceType, IDs, Titles)
}

func makeBody(c *cli.Context) *resourceParams {
	return &resourceParams{
		Tags: c.StringSlice("tags"),
	}
}

func getFolders(params *resourceParams) (IDs []string, Titles []string) {
	var resBody *foldersResp
	data := get.getRequest("https://app.rainforestqa.com/api/1/folders.json")
	json.Unmarshal(data, &resBody)
	IDs, Titles = resBody.TableSlice()
	return IDs, Titles
}

func getSites(params *resourceParams) (IDs []string, Titles []string) {
	var resBody *sitesResp
	data := get.getRequest("https://app.rainforestqa.com/api/1/sites.json")
	json.Unmarshal(data, &resBody)
	IDs, Titles = resBody.TableSlice()
	return IDs, Titles
}

func getBrowsers(params *resourceParams) (IDs []string, Titles []string) {
	var resBody *browsersResp
	data := get.getRequest("https://app.rainforestqa.com/api/1/clients.json")
	json.Unmarshal(data, &resBody)
	IDs, Titles = resBody.TableSlice()
	return IDs, Titles
}

func printResource(resource string, IDs []string, Titles []string) {
	fmt.Printf("%v Id\t| %v Name\n", resource, resource)
	bar := strings.Repeat("-", 40)
	print("" + bar + "\n")
	for i := range IDs {
		fmt.Printf("\t%v\t| %v\n", IDs[i], Titles[i])
	}
}
