package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

type reporterCliContext interface {
	String(flag string) (val string)
	Args() (args cli.Args)
}

type reporterClient interface {
	GetRunDetails(runID int) (*rainforest.RunDetails, error)
}

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
		log.Fatalf("Error parsing file path `%v`: %v", filepath, err.Error())
		return err
	}

	var runDetails *rainforest.RunDetails
	runDetails, err = client.GetRunDetails(runID)

	if err != nil {
		log.Fatalf("Error fetching details for run #%v: %v", runID, err.Error())
		return err
	}

	fmt.Println(filepath)
	fmt.Println(runDetails)
	// output := fmt.Sprintf("Info for run #%v", runID)
	// ioutil.WriteFile(filepath, []byte(output), 0777)

	return nil
}
