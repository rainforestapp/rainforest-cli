package main

import (
	"encoding/json"
	"os"

	"github.com/olekukonko/tablewriter"
)

var get getResponse

type resourceParams struct {
	Tags []string `json:"tags"`
}

func fetchResource(resourceType string, web resourceGetter) (table [][]string) {
	switch resourceType {
	case "Folders":
		table = web.getFolders()
	case "Sites":
		table = web.getSites()
	case "Browsers":
		table = web.getBrowsers()
	default:
		panic("Not valid resource to fetch")
	}
	printResource(resourceType, table)
	return table
}

func (g webResGetter) getFolders() (tableData [][]string) {
	var resBody *foldersResp
	data := get.getRequest("https://app.rainforestqa.com/api/1/folders.json?page_size=100")
	json.Unmarshal(data, &resBody)
	tableData = resBody.TableSlice()
	return tableData
}

func (g webResGetter) getSites() (tableData [][]string) {
	var resBody *sitesResp
	data := get.getRequest("https://app.rainforestqa.com/api/1/sites.json")
	json.Unmarshal(data, &resBody)
	tableData = resBody.TableSlice()
	return tableData
}

func (g webResGetter) getBrowsers() (tableData [][]string) {
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
