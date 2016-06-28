package main

import (
	"fmt"

	"github.com/olekukonko/tablewriter"
)

func printFolders() {
	var table [][]string
	var resBody foldersResp
	getFolders("folders.json?page_size=100", &resBody)
	fmt.Printf("\n%v\n", resBody)
	table = resBody.displayTable()
	printResource("Folders", table)
}

// func getSites() [][]string {
// 	var tableData [][]string
// 	var resBody sitesResp
// 	getResource("sites.json", resBody)
// 	tableData = resBody.displayTable()
// 	return tableData
// }
// func getBrowsers() [][]string {
// 	var tableData [][]string
// 	var resBody *browsersResp
// 	getResource("clients.json", resBody)
// 	tableData = resBody.displayTable()
// 	return tableData
// }

type resourceParams struct {
	Tags []string `json:"tags"`
}

// func fetchResource(resourceType string) error {
// 	var table [][]string
// 	switch resourceType {
// 	case "Folders":
// 		table = getFolders()
// 	case "Sites":
// 		table = getSites()
// 	case "Browsers":
// 		table = getBrowsers()
// 	default:
// 		return errors.New("Not valid resource to fetch")
// 	}
// 	printResource(resourceType, table)
// 	return nil
// }

func printResource(resource string, data [][]string) {
	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{resource + " ID", resource + " Description"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
