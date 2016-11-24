package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

func reportForRun(c *cli.Context) error {
	var runID int
	var err error

	if runIDArg := c.Args().Get(0); runIDArg != "" {
		runID, err = strconv.Atoi(runIDArg)
		if err != nil {
			return err
		}
	} else if deprecatedRunIDArg := c.String("run-id"); deprecatedRunIDArg != "" {
		runID, err = strconv.Atoi(deprecatedRunIDArg)
		if err != nil {
			return err
		}

		log.Println("Warning - `run-id` flag is deprecated. Please provide Run ID as an argument.")
	} else {
		return cli.NewExitError("No run found.", 1)
	}

	if junitFile := c.String("junit-file"); junitFile != "" {
		err = createJunitReport(junitFile, runID, api)
		if err != nil {
			return err
		}
	}

	return nil
}

func createJunitReport(filename string, runID int, api *rainforest.Client) error {
	filepath, err := filepath.Abs(filename)

	if err != nil {
		return err
	}

	output := fmt.Sprintf("Info for run #%v", runID)
	ioutil.WriteFile(filepath, []byte(output), 0777)

	return nil
}
