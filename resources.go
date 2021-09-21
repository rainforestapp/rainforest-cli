package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

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
	GetEnvironments() ([]rainforest.Environment, error)
	GetFeatures() ([]rainforest.Feature, error)
	GetRunGroups() ([]rainforest.RunGroup, error)
	GetRunJunit(int) (*string, error)
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

// printEnvironments fetches and prints out the defined enviroments
func printEnvironments(api resourceAPI) error {
	// Fetch the list of enviroments from the Rainforest
	environments, err := api.GetEnvironments()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	rows := make([][]string, len(environments))
	for i, environment := range environments {
		rows[i] = []string{strconv.Itoa(environment.ID), environment.Name}
	}

	printResourceTable([]string{"Environment ID", "Environment Name"}, rows)
	return nil
}

// printFeatures fetches and prints features
func printFeatures(api resourceAPI) error {
	// Fetch the list of features from the Rainforest
	features, err := api.GetFeatures()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	rows := make([][]string, len(features))
	for i, feature := range features {
		rows[i] = []string{strconv.Itoa(feature.ID), feature.Title}
	}

	printResourceTable([]string{"Feature ID", "Feature Title"}, rows)
	return nil
}

// printRunGroups fetches and prints runGroups
func printRunGroups(api resourceAPI) error {
	// Fetch the list of runGroups from the Rainforest
	runGroups, err := api.GetRunGroups()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	rows := make([][]string, len(runGroups))
	for i, runGroup := range runGroups {
		rows[i] = []string{strconv.Itoa(runGroup.ID), runGroup.Title}
	}

	printResourceTable([]string{"Run Group ID", "Run Group Title"}, rows)
	return nil
}

func postRunJUnitReport(c cliContext, runID int) error {
	// Get the csv file path either and skip uploading if it's not present
	fileName := c.String("junit-file")
	if fileName == "" {
		return nil
	}
	api = rainforest.NewClient(c.String("token"), c.Bool("debug"))

	cmd := []string{
		"rainforest-cli",
		"report", strconv.Itoa(runID),
		"--skip-update", // skip auto-updates for reports inside a run
	}

	if token := c.GlobalString("token"); len(token) > 0 {
		cmd = append(cmd, "--token", token)
	}
	cmd = append(cmd, "--junit-file", fileName)

	path, err := os.Executable()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	exec_err := syscall.Exec(path, cmd, os.Environ())
	if exec_err != nil {
		return cli.NewExitError(exec_err.Error(), 1)
	}
	return nil
}

func augmentJunitFileName(junitFile string, rerunAttempt uint) string {
	if rerunAttempt > 0 {
		junitFile = fmt.Sprintf("%v.%v", junitFile, rerunAttempt)
	}

	return junitFile
}

// write writeJunit fetches and writes a junit.xml file
func writeJunit(c cliContext, api resourceAPI) error {
	var runID int
	var err error

	if runIDArg := c.Args().Get(0); runIDArg != "" {
		runID, err = strconv.Atoi(runIDArg)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	} else {
		return cli.NewExitError("No run ID argument found.", 1)
	}

	junitFile := c.String("junit-file")
	if junitFile == "" {
		return cli.NewExitError("JUnit output file not specified", 1)
	}

	rerunAttempt := c.Uint("rerun-attempt")
	if rerunAttempt > 0 {
		junitFile = augmentJunitFileName(junitFile, rerunAttempt)
	}

	xml, err := api.GetRunJunit(runID)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	file, err := os.Create(junitFile)
	defer file.Close()

	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	} else {
		file.WriteString(*xml)
	}

	return nil
}
