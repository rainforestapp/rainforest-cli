package main

import (
	"flag"
	"io"
	"os"
)

var (
	apiToken        string
	smartFolderID   int
	siteID          int
	tags            string
	testIDs         string
	baseURL                   = "https://app.rainforestqa.com/api/1"
	out             io.Writer = os.Stdout
	allowForeGround           = true

	crowd               string
	conflict            string
	browsers            string //[]string
	description         string
	environmentID       int
	testingNoForeGround bool
)

func parseArgs(arguments []string) ([]string, []string) {
	var commands []string
	var flags []string
	for i := 1; i < len(arguments); i++ {
		if arguments[i][0] != '-' {
			commands = append(commands, arguments[i])
		} else {
			flags = append(flags, arguments[i])
		}
	}
	return commands, flags
}

func main() {
	commands, flags := parseArgs(os.Args)
	command := commands[0]

	flag.StringVar(&apiToken, "token", "", "API token. You can find your account token at https://app.rainforestqa.com/settings/integrations")
	flag.IntVar(&smartFolderID, "smart_folder_id", 0, "Smart Folder Id. use the `folders` command to find the ID's of your smart folders")
	flag.IntVar(&siteID, "site_id", 0, "Site ID. use the `sites` command to find the ID's of your sites")
	flag.StringVar(&tags, "tags", "", "Test tags. enter in a comma separated list")
	flag.StringVar(&testIDs, "tests", "", "Run test by id. Enter in a comma separated list")

	flag.StringVar(&crowd, "crowd", "", "Crowd to run test with. Enter `default` or `on_premise_crowd`")
	flag.StringVar(&conflict, "conflict", "", "Handling of runs in progress. (A) Abort: only abort runs in the same environment as your new run. (B) Abort All.")
	flag.StringVar(&browsers, "browsers", "", "Browsers to test against. enter in a comma separated list")
	flag.StringVar(&description, "description", "", "An arbitrary string to associate with the run")
	flag.IntVar(&environmentID, "environment_id", 0, "Use a specific environment for this run")
	flag.CommandLine.Parse(flags)

	if len(apiToken) == 0 {
		envToken, present := os.LookupEnv("RAINFOREST_API_TOKEN")

		if present {
			apiToken = envToken
		}
	}

	switch command {
	case "run":
		createRun()
	case "sites":
		printSites()
	case "folders":
		printFolders()
	case "browsers":
		printBrowsers()
	default:
		// TODO: Print out usage
	}
}
