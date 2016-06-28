package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/olekukonko/tablewriter"
)

func getFolders() [][]string {
	var tableData [][]string
	var resBody *foldersResp
	data := getRequest("folders.json?page_size=100")
	json.Unmarshal(data, &resBody)
	tableData = resBody.displayTable()
	return tableData
}

func getSites() [][]string {
	var tableData [][]string
	var resBody *sitesResp
	data := getRequest("sites.json")
	err := json.Unmarshal(data, &resBody)
	if err != nil {
		log.Fatalf("Invalid response from API: %v\n", err)
	}
	tableData = resBody.displayTable()
	return tableData
}
func getBrowsers() [][]string {
	var tableData [][]string
	var resBody *browsersResp
	data := getRequest("clients.json")
	json.Unmarshal(data, &resBody)
	tableData = resBody.displayTable()
	return tableData
}

type resourceParams struct {
	Tags []string `json:"tags"`
}

func fetchResource(resourceType string) error {
	var table [][]string
	switch resourceType {
	case "Folders":
		table = getFolders()
	case "Sites":
		table = getSites()
	case "Browsers":
		table = getBrowsers()
	default:
		return errors.New("Not valid resource to fetch")
	}
	printResource(resourceType, table)
	return nil
}

func printResource(resource string, data [][]string) {
	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{resource + " ID", resource + " Description"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
