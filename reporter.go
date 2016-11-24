package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"

	"github.com/urfave/cli"
)

type reporterCliContext interface {
	String(flag string) (val string)
	Args() (args cli.Args)
}

type reporterClient interface{}

type reporter struct {
	createJunitReport func(filename string, runID int, client reporterClient) error
}

func createReport(c *cli.Context) error {
	r := newReport()
	return r.reportForRun(c)
}

func newReport() *reporter {
	return &reporter{
		createJunitReport: createJunitReport,
	}
}

func (r *reporter) reportForRun(c reporterCliContext) error {
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

		log.Println("Warning: `run-id` flag is deprecated. Please provide Run ID as an argument.")
	} else {
		return cli.NewExitError("No run ID argument found.", 1)
	}

	if junitFile := c.String("junit-file"); junitFile != "" {
		err = r.createJunitReport(junitFile, runID, api)
		if err != nil {
			return err
		}
	}

	return nil
}

func createJunitReport(filename string, runID int, client reporterClient) error {
	filepath, err := filepath.Abs(filename)

	if err != nil {
		return err
	}

	output := fmt.Sprintf("Info for run #%v", runID)
	ioutil.WriteFile(filepath, []byte(output), 0777)

	return nil
}
