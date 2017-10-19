package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"errors"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

// Maximum concurrency for multithreaded HTTP requests
const reporterConcurrency = 4

// resourceAPI is part of the API connected to available resources
type reporterAPI interface {
	GetRunTestDetails(int, int) (*rainforest.RunTestDetails, error)
	GetRunDetails(int) (*rainforest.RunDetails, error)
}

type reporter struct {
	getRunDetails           func(int, reporterAPI) (*rainforest.RunDetails, error)
	createOutputFile        func(string) (*os.File, error)
	createJUnitReportSchema func(*rainforest.RunDetails, reporterAPI) (*jUnitReportSchema, error)
	writeJUnitReport        func(*jUnitReportSchema, *os.File) error
}

func createReport(c cliContext) error {
	r := newReporter()
	return r.createReport(c)
}

func postRunJUnitReport(c cliContext, runID int) error {
	// Get the csv file path either and skip uploading if it's not present
	fileName := c.String("junit-file")
	if fileName == "" {
		return nil
	}

	r := newReporter()
	return r.createJUnitReport(runID, fileName)
}

func newReporter() *reporter {
	return &reporter{
		getRunDetails:           getRunDetails,
		createOutputFile:        os.Create,
		createJUnitReportSchema: createJUnitReportSchema,
		writeJUnitReport:        writeJUnitReport,
	}
}

func (r *reporter) createReport(c cliContext) error {
	var runID int
	var err error

	if runIDArg := c.Args().Get(0); runIDArg != "" {
		runID, err = strconv.Atoi(runIDArg)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	} else if deprecatedRunIDArg := c.String("run-id"); deprecatedRunIDArg != "" {
		runID, err = strconv.Atoi(deprecatedRunIDArg)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		log.Println("Warning: `run-id` flag is deprecated. Please provide Run ID as an argument.")
	} else {
		return cli.NewExitError("No run ID argument found.", 1)
	}

	if junitFile := c.String("junit-file"); junitFile != "" {
		err = r.createJUnitReport(runID, junitFile)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	} else {
		return cli.NewExitError("Output file not specified", 1)
	}

	return nil
}

func (r *reporter) createJUnitReport(runID int, junitFile string) error {
	log.Print("Creating JUnit report for run #" + strconv.Itoa(runID) + ": " + junitFile)

	if filepath.Ext(junitFile) != ".xml" {
		return errors.New("JUnit file extension must be .xml")
	}

	filepath, err := filepath.Abs(junitFile)
	if err != nil {
		return err
	}

	var runDetails *rainforest.RunDetails
	runDetails, err = r.getRunDetails(runID, api)
	if err != nil {
		return err
	}

	var outputFile *os.File
	outputFile, err = r.createOutputFile(filepath)
	defer outputFile.Close()
	if err != nil {
		return err
	}

	var reportSchema *jUnitReportSchema
	reportSchema, err = r.createJUnitReportSchema(runDetails, api)
	if err != nil {
		return err
	}

	err = r.writeJUnitReport(reportSchema, outputFile)
	if err != nil {
		return err
	}

	return nil
}

func getRunDetails(runID int, client reporterAPI) (*rainforest.RunDetails, error) {
	var runDetails *rainforest.RunDetails
	var err error

	log.Printf("Fetching details for run #" + strconv.Itoa(runID))
	if runDetails, err = client.GetRunDetails(runID); err != nil {
		return runDetails, err
	}

	if !runDetails.StateDetails.IsFinalState {
		err = errors.New("Report cannot be created for an incomplete run")
	}

	return runDetails, err
}

type jUnitTestReportFailure struct {
	XMLName xml.Name `xml:"failure"`
	Type    string   `xml:"type,attr"`
	Message string   `xml:"message,attr"`
}

type jUnitTestReportSchema struct {
	XMLName  xml.Name `xml:"testcase"`
	Name     string   `xml:"name,attr"`
	Time     float64  `xml:"time,attr"`
	Failures []jUnitTestReportFailure
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

func createJUnitReportSchema(runDetails *rainforest.RunDetails, api reporterAPI) (*jUnitReportSchema, error) {
	type processedTestCase struct {
		TestCase jUnitTestReportSchema
		Error    error
	}

	// Create channels for work to be done and results
	testCount := len(runDetails.Tests)
	processedTestChan := make(chan processedTestCase, testCount)
	testsChan := make(chan rainforest.RunTestDetails, testCount)

	processTestWorker := func(testsChan <-chan rainforest.RunTestDetails) {
		for test := range testsChan {
			testCase := jUnitTestReportSchema{}

			testDuration := test.UpdatedAt.Sub(test.CreatedAt).Seconds()

			testCase.Name = test.Title
			testCase.Time = testDuration

			if test.Result == "failed" {
				log.Printf("Fetching information for failed test #" + strconv.Itoa(test.ID))
				testDetails, err := api.GetRunTestDetails(runDetails.ID, test.ID)

				if err != nil {
					processedTestChan <- processedTestCase{TestCase: jUnitTestReportSchema{}, Error: err}
					return
				}

				for _, step := range testDetails.Steps {
					for _, browser := range step.Browsers {
						for _, feedback := range browser.Feedback {
							if feedback.Result != "failed" || feedback.JobState != "approved" {
								continue
							}

							if feedback.FailureNote != "" {
								reportFailure := jUnitTestReportFailure{Type: browser.Name, Message: feedback.FailureNote}
								testCase.Failures = append(testCase.Failures, reportFailure)
							} else if feedback.CommentReason != "" {
								// The step failed due to a special comment type being selected
								message := feedback.CommentReason

								if feedback.Comment != "" {
									message += ": " + feedback.Comment
								}
								reportFailure := jUnitTestReportFailure{Type: browser.Name, Message: message}
								testCase.Failures = append(testCase.Failures, reportFailure)
							}
						}
					}
				}
			}

			processedTestChan <- processedTestCase{TestCase: testCase}
		}
	}

	// spawn workers
	for i := 0; i < reporterConcurrency; i++ {
		go processTestWorker(testsChan)
	}

	// give them work
	for _, test := range runDetails.Tests {
		testsChan <- test
	}
	close(testsChan)

	// and collect the results
	testCases := make([]jUnitTestReportSchema, testCount)
	for i := 0; i < testCount; i++ {
		processed := <-processedTestChan

		if processed.Error != nil {
			return &jUnitReportSchema{}, processed.Error
		}

		testCases[i] = processed.TestCase
	}

	finalStateName := runDetails.StateDetails.Name
	runDuration := runDetails.Timestamps[finalStateName].Sub(runDetails.Timestamps["created_at"]).Seconds()
	report := &jUnitReportSchema{
		Errors:    runDetails.TotalNoResultTests,
		Failures:  runDetails.TotalFailedTests,
		Tests:     runDetails.TotalTests,
		TestCases: testCases,
		Time:      runDuration,
	}

	if runDesc := runDetails.Description; runDesc != "" {
		report.Name = runDesc
	} else {
		report.Name = fmt.Sprintf("Run #%v", runDetails.ID)
	}

	return report, nil
}

func writeJUnitReport(reportSchema *jUnitReportSchema, file *os.File) error {
	enc := xml.NewEncoder(file)

	file.Write([]byte(xml.Header))

	enc.Indent("", "  ")
	err := enc.Encode(reportSchema)
	if err != nil {
		return err
	}

	log.Printf("JUnit report successfully written to %v", file.Name())
	return nil
}
