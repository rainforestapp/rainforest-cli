package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

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
	createJUnitReport func(filename string, runID int, client reporterClient) error
}

func createReport(c *cli.Context) error {
	r := newReport()
	return r.reportForRun(c)
}

func newReport() *reporter {
	return &reporter{
		createJUnitReport: createJUnitReport,
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
		err = r.createJUnitReport(junitFile, runID, api)
		if err != nil {
			return err
		}
	}

	return nil
}

// JUnitReport defines the format of the JUnit XML report.
type JUnitReport struct {
	XMLName  xml.Name `xml:"testsuite"`
	Name     string   `xml:"name,attr"`
	Tests    int      `xml:"tests,attr"`
	Errors   int      `xml:"errors,attr"`
	Failures int      `xml:"failures,attr"`
	Time     float64  `xml:"time,attr"`
}

func createJUnitReport(filename string, runID int, client reporterClient) error {
	filepath, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalf("Error parsing file path `%v`: %v", filepath, err.Error())
		return err
	}

	var runDetails *rainforest.RunDetails
	if runDetails, err = client.GetRunDetails(runID); err != nil {
		log.Fatalf("Error fetching details for run #%v: %v", runID, err.Error())
		return err
	}

	if !runDetails.StateDetails.IsFinalState {
		log.Fatalf("Report cannot be created for an incomplete run")
		return fmt.Errorf("Report cannot be created for an incomplete run")
	}

	var file *os.File
	if file, err = os.Create(filepath); err != nil {
		log.Fatalf("Error creating file at %v: %v", filepath, err.Error())
		return err
	}

	file.Write([]byte(xml.Header))

	enc := xml.NewEncoder(file)
	var createdAt time.Time
	var completedAt time.Time

	createdAt, err = time.Parse(time.RFC3339Nano, runDetails.Timestamps["created_at"])
	if err != nil {
		log.Fatalf("Error parsing Run timestamp %v: %v", runDetails.Timestamps["created_at"], err.Error())
		return err
	}

	finalStateName := runDetails.StateDetails.Name
	completedAt, err = time.Parse(time.RFC3339Nano, runDetails.Timestamps[finalStateName])
	if err != nil {
		log.Fatalf("Error parsing Run timestamp %v: %v", runDetails.Timestamps[finalStateName], err.Error())
		return err
	}

	v := &JUnitReport{
		Name:     runDetails.Description,
		Errors:   runDetails.TotalNoResultTests,
		Failures: runDetails.TotalFailedTests,
		Tests:    runDetails.TotalTests,
		Time:     completedAt.Sub(createdAt).Seconds(),
	}

	err = enc.Encode(v)
	if err != nil {
		log.Fatalf("Error encoding XML report: %v", err.Error())
		return err
	}

	return nil
}
