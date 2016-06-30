package main

import (
	"flag"
	"io"
	"os"
)

var apiToken string

var baseURL = "https://app.rainforestqa.com/api/1"
var out io.Writer = os.Stdout

func parseCommands(arguments []string) []string {
	commands := make([]string, 0, 5)
	for i := len(arguments) - 1; i > 0; i-- {
		if arguments[i][0] != '-' {
			commands = append(commands, arguments[i])
			os.Args = append(arguments[:i], arguments[i+1:]...)
		}
	}
	return commands
}

func main() {
	commands := parseCommands(os.Args)
	command := commands[0]

	flag.StringVar(&apiToken, "token", "", "API token. You can find your account token at https://app.rainforestqa.com/settings/integrations")
	flag.Parse()

	if len(apiToken) == 0 {
		envToken, present := os.LookupEnv("RAINFOREST_API_TOKEN")

		if present {
			apiToken = envToken
		}
	}

	switch command {
	case "runs":

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
