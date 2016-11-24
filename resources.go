package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

// formatAsTable takes slice of a pointers of a type which represents one of the Resources
// returned from API and returns it formated as a slice of [id, description] rows converted to string.
func formatAsTable(resources []rainforest.Resource) [][]string {
	// Create slice which will hold the resulting table
	table := make([][]string, 0, len(resources))

	// Iterate over slice of resources and convert them to tabular format
	for _, resource := range resources {
		tableRow := []string{resource.GetID(), resource.GetDescription()}
		table = append(table, tableRow)
	}
	return table
}

// printResourceTable uses olekukonko/tablewriter as a pretty printer
// for the tabular resources we get from the API and formatted using formatAsTable.
func printResourceTable(resourceName string, data [][]string) {
	// Init tablewriter with out global var as a target
	table := tablewriter.NewWriter(out)

	// Prepare the tablewriter
	table.SetHeader([]string{resourceName + " ID", resourceName + " Description"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data

	// Render prints out the table to output specified during init.
	table.Render()
}

// resourceAPI is part of the API connected to available resources
type resourceAPI interface {
	GetFolders() ([]rainforest.Folder, error)
	GetBrowsers() ([]rainforest.Browser, error)
	GetSites() ([]rainforest.Site, error)
}

// printFolders fetches and prints out the available folders from the API
func printFolders(c *cli.Context, api resourceAPI) error {
	// Fetch the list of folders from the Rainforest
	folders, err := api.GetFolders()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// We need to manually cast to Resource type as we have a slice of the folders not resources
	// and go will refuse to do it behind the scenes as it's O(n)
	resources := make([]rainforest.Resource, len(folders))
	for i, v := range folders {
		resources[i] = rainforest.Resource(v)
	}

	// Now we can use it as an argument for our cool printing methods
	tabular := formatAsTable(resources)
	printResourceTable("Folder", tabular)
	return nil
}

// printBrowsers fetches and prints out the browsers available to the client
func printBrowsers(c *cli.Context, api resourceAPI) error {
	// Fetch the list of browsers from the Rainforest
	browsers, err := api.GetBrowsers()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// We need to manually cast to Resource type as we have a slice of the browsers not resources
	// and go will refuse to do it behind the scenes as it's O(n)
	resources := make([]rainforest.Resource, len(browsers))
	for i, v := range browsers {
		resources[i] = rainforest.Resource(v)
	}

	// Now we can use it as an argument for our cool printing methods
	tabular := formatAsTable(resources)
	printResourceTable("Browser", tabular)
	return nil
}

// printSites fetches and prints out the defined sites
func printSites(c *cli.Context, api resourceAPI) error {
	// Fetch the list of sites from the Rainforest
	sites, err := api.GetSites()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// We need to manually cast to Resource type as we have a slice of the sites not resources
	// and go will refuse to do it behind the scenes as it's O(n)
	resources := make([]rainforest.Resource, len(sites))
	for i, v := range sites {
		resources[i] = rainforest.Resource(v)
	}

	// Now we can use it as an argument for our cool printing methods
	tabular := formatAsTable(resources)
	printResourceTable("Site", tabular)
	return nil
}
