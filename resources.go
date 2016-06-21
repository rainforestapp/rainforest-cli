package main

import (
	"encoding/json"
	"os"

	"github.com/olekukonko/tablewriter"

	"gopkg.in/urfave/cli.v2"
)

var get getResponse

type resourceParams struct {
	Tags []string `json:"tags"`
}

func createResource(c *cli.Context, resourceType string) {
	params := makeBody(c)
	var table [][]string
	switch resourceType {
	case "Folders":
		table = getFolders(params)
	case "Sites":
		table = getSites(params)
	default: //"Browsers"
		table = getBrowsers(params)
	}
	printResource(resourceType, table)
}

func makeBody(c *cli.Context) *resourceParams {
	return &resourceParams{
		Tags: c.StringSlice("tags"),
	}
}

func getFolders(params *resourceParams) (tableData [][]string) {
	var resBody *foldersResp
	data := get.getRequest("https://app.rainforestqa.com/api/1/folders.json?page_size=100")
	json.Unmarshal(data, &resBody)
	tableData = resBody.TableSlice()
	return tableData
}

func getSites(params *resourceParams) (tableData [][]string) {
	var resBody *sitesResp
	data := get.getRequest("https://app.rainforestqa.com/api/1/sites.json")
	json.Unmarshal(data, &resBody)
	tableData = resBody.TableSlice()
	return tableData
}

func getBrowsers(params *resourceParams) (tableData [][]string) {
	var resBody *browsersResp
	data := get.getRequest("https://app.rainforestqa.com/api/1/clients.json")
	json.Unmarshal(data, &resBody)
	tableData = resBody.TableSlice()
	return tableData
}

func printResource(resource string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{resource + " ID", resource + " Description"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
