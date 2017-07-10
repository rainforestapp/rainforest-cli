package main

import (
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

// printResourceTable uses olekukonko/tablewriter as a pretty printer
// for the tabular resources we get from the API and formatted using formatAsTable.
func printResourceTable(headers []string, rows [][]string) {
	// Init tablewriter with out global var as a target
	table := tablewriter.NewWriter(tablesOut)

	// Prepare the tablewriter
	table.SetHeader(headers)
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(rows) // Add Bulk Data

	// Render prints out the table to output specified during init.
	table.Render()
}

// resourceAPI is part of the API connected to available resources
type resourceAPI interface {
	GetFolders() ([]rainforest.Folder, error)
	GetBrowsers() ([]rainforest.Browser, error)
	GetSites() ([]rainforest.Site, error)
	GetRunGroups() ([]rainforest.RunGroup, error)
}

func printRunGroupDetails(c cliContext) error {
	runGroupIDStr := c.Args().First()
	if runGroupIDStr == "" {
		return cli.NewExitError("Please supply the Run Group ID. Use the `run-groups` command to retrieve a list of Run Group IDs.", 1)
	}

	runGroupID, err := strconv.Atoi(runGroupIDStr)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	runGroup, err := api.GetRunGroupDetails(runGroupID)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	runGroup.Print()
	return nil
}

func printRunGroups(api resourceAPI) error {

	runGroups, err := api.GetRunGroups()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	rows := make([][]string, len(runGroups))
	for i, rungrp := range runGroups {
		rows[i] = []string{strconv.Itoa(rungrp.ID), rungrp.Title}
	}

	printResourceTable([]string{"Run Group ID", "Run Group Name"}, rows)
	return nil
}

// printFolders fetches and prints out the available folders from the API
func printFolders(api resourceAPI) error {
	// Fetch the list of folders from the Rainforest
	folders, err := api.GetFolders()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	rows := make([][]string, len(folders))
	for i, folder := range folders {
		rows[i] = []string{strconv.Itoa(folder.ID), folder.Title}
	}

	printResourceTable([]string{"Folder ID", "Folder Name"}, rows)
	return nil
}

// printBrowsers fetches and prints out the browsers available to the client
func printBrowsers(api resourceAPI) error {
	// Fetch the list of browsers from the Rainforest
	browsers, err := api.GetBrowsers()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	rows := make([][]string, len(browsers))
	for i, browser := range browsers {
		rows[i] = []string{browser.Name, browser.Description}
	}

	printResourceTable([]string{"Browser ID", "Browser Name"}, rows)
	return nil
}

// printSites fetches and prints out the defined sites
func printSites(api resourceAPI) error {
	// Fetch the list of sites from the Rainforest
	sites, err := api.GetSites()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	humanizedSiteCategories := map[string]string{
		"device_farm": "Device Farm",
		"android":     "Android",
		"ios":         "iOS",
		"site":        "Site",
	}

	rows := make([][]string, len(sites))
	for i, site := range sites {
		category, ok := humanizedSiteCategories[site.Category]
		if !ok {
			category = site.Category
		}
		rows[i] = []string{strconv.Itoa(site.ID), site.Name, category}
	}

	printResourceTable([]string{"Site ID", "Site Name", "Category"}, rows)
	return nil
}
