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

type reporterClient interface {
	GetRunDetails(runID int) (*rainforest.RunDetails, error)
}

type reporter struct {
	getRunDetails               func(runID int, client *rainforest.Client) (*rainforest.RunDetails, error)
	createOutputFile            func(filepath string) (*os.File, error)
	createJunitReportSchema     func(*rainforest.RunDetails) (*jUnitReportSchema, error)
	createJunitTestReportSchema func(tests *[]rainforest.RunTestDetails) (*[]jUnitTestReportSchema, error)
	writeJUnitReport            func(*jUnitReportSchema, *os.File) error
}

func createReport(c *cli.Context) error {
	r := newReporter()
	return r.createReport(c)
}

func newReporter() *reporter {
	return &reporter{
		getRunDetails:               getRunDetails,
		createOutputFile:            createOutputFile,
		createJunitReportSchema:     createJunitReportSchema,
		createJunitTestReportSchema: createJunitTestReportSchema,
		writeJUnitReport:            writeJUnitReport,
	}
}

func (r *reporter) createReport(c cliContext) error {
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
		err = r.createJUnitReport(runID, junitFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *reporter) createJUnitReport(runID int, junitFile string) error {
	if filepath.Ext(junitFile) != ".xml" {
		errMessage := "JUnit file extension must be .xml"
		log.Fatal(errMessage)
		return fmt.Errorf(errMessage)
	}

	filepath, err := filepath.Abs(junitFile)
	if err != nil {
		log.Fatalf("Error parsing file path `%v`: %v", filepath, err.Error())
		return err
	}

	var runDetails *rainforest.RunDetails
	runDetails, err = r.getRunDetails(runID, api)
	if err != nil {
		return err
	}

	var outputFile *os.File
	outputFile, err = r.createOutputFile(filepath)
	if err != nil {
		return err
	}

	var reportSchema *jUnitReportSchema
	reportSchema, err = r.createJunitReportSchema(runDetails)
	if err != nil {
		return err
	}

	err = r.writeJUnitReport(reportSchema, outputFile)
	if err != nil {
		return err
	}

	return nil
}

func getRunDetails(runID int, client *rainforest.Client) (*rainforest.RunDetails, error) {
	var runDetails *rainforest.RunDetails
	var err error
	if runDetails, err = client.GetRunDetails(runID); err != nil {
		log.Fatalf("Error fetching details for run #%v: %v", runID, err.Error())
		return runDetails, err
	}

	if !runDetails.StateDetails.IsFinalState {
		errMessage := "Report cannot be created for an incomplete run"
		log.Fatalf(errMessage)
		err = fmt.Errorf(errMessage)
	}

	return runDetails, err
}

type jUnitTestReportSchema struct {
	XMLName xml.Name `xml:"testcase"`
	Name    string   `xml:"name,attr"`
	Time    float64  `xml:"time,attr"`
}

type jUnitReportSchema struct {
	XMLName   xml.Name `xml:"testsuite"`
	Name      string   `xml:"name,attr"`
	Tests     int      `xml:"tests,attr"`
	Errors    int      `xml:"errors,attr"`
	Failures  int      `xml:"failures,attr"`
	Time      float64  `xml:"time,attr"`
	TestCases []jUnitTestReportSchema
}

func createJunitReportSchema(runDetails *rainforest.RunDetails) (*jUnitReportSchema, error) {
	var err error
	var duration float64

	finalStateName := runDetails.StateDetails.Name
	duration, err = timeStringDifferenceSecs(runDetails.Timestamps["created_at"], runDetails.Timestamps[finalStateName])
	if err != nil {
		return &jUnitReportSchema{}, err
	}

	var testCases *[]jUnitTestReportSchema
	testCases, err = createJunitTestReportSchema(&runDetails.Tests)
	if err != nil {
		return &jUnitReportSchema{}, err
	}

	report := &jUnitReportSchema{
		Name:      runDetails.Description,
		Errors:    runDetails.TotalNoResultTests,
		Failures:  runDetails.TotalFailedTests,
		Tests:     runDetails.TotalTests,
		TestCases: *testCases,
		Time:      duration,
	}

	return report, nil
}

func createJunitTestReportSchema(tests *[]rainforest.RunTestDetails) (*[]jUnitTestReportSchema, error) {
	testCases := []jUnitTestReportSchema{}

	for _, test := range *tests {
		duration, err := timeStringDifferenceSecs(test.CreatedAt, test.UpdatedAt)

		if err != nil {
			return &[]jUnitTestReportSchema{}, err
		}

		testCase := jUnitTestReportSchema{
			Name: test.Title,
			Time: duration,
		}

		testCases = append(testCases, testCase)
	}

	return &testCases, nil
}

func writeJUnitReport(reportSchema *jUnitReportSchema, file *os.File) error {
	enc := xml.NewEncoder(file)

	file.Write([]byte(xml.Header))

	enc.Indent("", "  ")
	err := enc.Encode(reportSchema)
	if err != nil {
		log.Fatalf("Error encoding XML report: %v", err.Error())
		return err
	}

	return nil
}

/*
	Utility Functions
*/

func createOutputFile(filepath string) (*os.File, error) {
	var file *os.File
	var err error
	if file, err = os.Create(filepath); err != nil {
		log.Fatalf("Error creating file at %v: %v", filepath, err.Error())
	}
	return file, err
}

func timeStringDifferenceSecs(start string, end string) (float64, error) {
	var startTime time.Time
	var endTime time.Time
	var err error
	var diff float64

	if startTime, err = time.Parse(time.RFC3339Nano, start); err != nil {
		log.Fatalf("Error in parsing time string %v: %v", start, err)
		return diff, err
	}

	if endTime, err = time.Parse(time.RFC3339Nano, end); err != nil {
		log.Fatalf("Error in parsing time string %v: %v", start, end)
		return diff, err
	}

	diff = endTime.Sub(startTime).Seconds()

	return diff, nil
}
