package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

func reportForRun(c *cli.Context) error {
	if junitFile := c.String("junit-file"); junitFile != "" {
		createJunitReport(junitFile, api)
	}

	return nil
}

func createJunitReport(filename string, api *rainforest.Client) {
	filepath, err := filepath.Abs(filename)

	if err != nil {
		log.Fatalf("Error parsing file path: %v", err)
	}

	ioutil.WriteFile(filepath, []byte("this is a test"), 0777)
}
