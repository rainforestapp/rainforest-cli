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

func fetchResource(c *cli.Context, resourceType string) {
	var table [][]string
	switch resourceType {
	case "Folders":
		table = getFolders()
	case "Sites":
		table = getSites()
	case "Browsers":
		table = getBrowsers()
	default:
		panic("Not valid resource to fetch")
	}
	printResource(resourceType, table)
}

func getFolders() (tableData [][]string) {
	var resBody *foldersResp
	data := get.getRequest("https://app.rainforestqa.com/api/1/folders.json?page_size=100")
	json.Unmarshal(data, &resBody)
	tableData = resBody.TableSlice()
	return tableData
}

func getSites() (tableData [][]string) {
	var resBody *sitesResp
	data := get.getRequest("https://app.rainforestqa.com/api/1/sites.json")
	json.Unmarshal(data, &resBody)
	tableData = resBody.TableSlice()
	return tableData
}

func getBrowsers() (tableData [][]string) {
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
