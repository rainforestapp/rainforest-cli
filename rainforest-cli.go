package main

import (
	"flag"
	"io"
	"os"
)

var apiToken string

var baseURL = "https://app.rainforestqa.com/api/1"
var out io.Writer = os.Stdout

func parseCommands() []string {
	commands := make([]string, 0, 5)
	for i := len(os.Args) - 1; i > 0; i-- {
		if os.Args[i][0] != '-' {
			commands = append(commands, os.Args[i])
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
		}
	}
	return commands
}

func main() {
	commands := parseCommands()
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
