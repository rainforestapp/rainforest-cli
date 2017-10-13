package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"log"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

func newFakeReporter() *reporter {
	r := newReporter()

	r.createJUnitReportSchema = func(*rainforest.RunDetails, reporterAPI) (*jUnitReportSchema, error) {
		return &jUnitReportSchema{}, nil
	}

	r.writeJUnitReport = func(*jUnitReportSchema, *os.File) error {
		return nil
	}

	r.getRunDetails = func(int, reporterAPI) (*rainforest.RunDetails, error) {
		return &rainforest.RunDetails{}, nil
	}

	r.createOutputFile = func(path string) (*os.File, error) {
		return os.NewFile(1, "test"), nil
	}

	return r
}

func TestReporterCreateReport_WithoutFlags(t *testing.T) {
	// No Flags
	r := newReporter()
	c := newFakeContext(make(map[string]interface{}), cli.Args{})

	err := r.createReport(c)

	if err == nil {
		t.Error("No error produced in reporter.reportForRun when Run ID is omitted")
	} else {
		if err.Error() != "No run ID argument found." {
			t.Errorf("Unexpected error in reporter.reportForRun when Run ID is omitted: %v", err.Error())
		}
	}
}

func TestReporterCreateReport(t *testing.T) {
	var expectedFileName string
	var expectedRunID int

	r := newFakeReporter()

	r.getRunDetails = func(runID int, client reporterAPI) (*rainforest.RunDetails, error) {
		runDetails := rainforest.RunDetails{}
		if runID != expectedRunID {
			t.Errorf("Unexpected run ID given to createJunitReport.\nExpected: %v\nActual: %v", expectedRunID, runID)
			return &runDetails, fmt.Errorf("Test failed")
		}

		return &runDetails, nil
	}

	r.createOutputFile = func(path string) (*os.File, error) {
		filename := filepath.Base(path)

		if filename != expectedFileName {
			t.Errorf("Unexpected filename given to createJunitReport.\nExpected: %v\nActual: %v", expectedFileName, filename)
		}

		return os.Create("myfilename.xml")
	}

	defer os.Remove("myfilename.xml")

	testCases := []struct {
		mappings map[string]interface{}
		args     []string
		runID    int
		filename string
	}{
		{
			mappings: map[string]interface{}{
				"junit-file": "myfilename.xml",
			},
			args:     []string{"112233"},
			runID:    112233,
			filename: "myfilename.xml",
		},
		{
			mappings: map[string]interface{}{
				"junit-file": "myfilename.xml",
				"run-id":     "112233",
			},
			args:     []string{},
			runID:    112233,
			filename: "myfilename.xml",
		},
	}

	for _, testCase := range testCases {
		c := newFakeContext(testCase.mappings, testCase.args)
		expectedRunID = testCase.runID
		expectedFileName = testCase.filename

		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stdout)
		err := r.createReport(c)
		if err != nil {
			t.Errorf("Unexpected error in reporter.createReport: %v", err.Error())
		}
	}
}

type fakeReporterAPI struct {
	RunMappings map[int][]rainforest.RunTestDetails
}

func (api fakeReporterAPI) GetRunTestDetails(runID int, testID int) (*rainforest.RunTestDetails, error) {
	runTests, ok := api.RunMappings[runID]
	if !ok {
		return nil, fmt.Errorf("No Run found with ID %v", runID)
	}

	for _, runTestDetails := range runTests {
		if runTestDetails.ID == testID {
			return &runTestDetails, nil
		}
	}

	return nil, fmt.Errorf("No RunTest found with ID %v", testID)
}

func (api fakeReporterAPI) GetRunDetails(int) (*rainforest.RunDetails, error) {
	// implement when needed
	return nil, errStub
}

func newFakeReporterAPI(runID int, runTestDetails []rainforest.RunTestDetails) *fakeReporterAPI {
	return &fakeReporterAPI{
		RunMappings: map[int][]rainforest.RunTestDetails{
			runID: runTestDetails,
		},
	}
}

func TestCreateJUnitReportSchema(t *testing.T) {
	// Without failures
	now := time.Now()
	runDesc := "very descriptive description"
	totalTests := 1
	totalNoResultTests := 0
	totalFailedTests := 0
	stateName := "complete"

	runDetails := rainforest.RunDetails{
		ID:                 123,
		Description:        runDesc,
		TotalTests:         totalTests,
		TotalNoResultTests: totalNoResultTests,
		TotalFailedTests:   totalFailedTests,
		StateDetails: rainforest.RunStateDetails{
			Name:         stateName,
			IsFinalState: true,
		},
		Timestamps: map[string]time.Time{
			"created_at":  now.Add(-30 * time.Minute),
			"in_progress": now.Add(-25 * time.Minute),
			stateName:     now,
		},
		Tests: []rainforest.RunTestDetails{
			{
				Title:     "My test title",
				CreatedAt: now.Add(-25 * time.Minute),
				UpdatedAt: now,
				Result:    "passed",
			},
		},
	}

	// api should not be used so input shouldn't matter
	api := newFakeReporterAPI(-1, []rainforest.RunTestDetails{})

	schema, err := createJUnitReportSchema(&runDetails, api)
	if err != nil {
		t.Errorf("Unexpected error returned by createJunitTestReportSchema: %v", err)
	}

	expectedSchema := jUnitReportSchema{
		Name:     runDesc,
		Errors:   totalNoResultTests,
		Failures: totalFailedTests,
		Tests:    totalTests,
		Time:     30 * time.Minute.Seconds(),
		TestCases: []jUnitTestReportSchema{
			{
				Name: runDetails.Tests[0].Title,
				Time: 25 * time.Minute.Seconds(),
			},
		},
	}

	if !reflect.DeepEqual(expectedSchema, *schema) {
		t.Error("Incorrect JUnitTestReportSchema returned by createJunitTestReportSchema")
		t.Errorf("Expected: %#v", expectedSchema)
		t.Errorf("Actual: %#v", schema)
	}

	// With failures
	failedBrowser := "chrome"
	failedNote := "This note should appear"

	runDetails.TotalFailedTests = 1

	failedTest := rainforest.RunTestDetails{
		ID:        999888,
		Title:     "My failed test",
		CreatedAt: now.Add(-25 * time.Minute),
		UpdatedAt: now,
		Result:    "failed",
	}
	runDetails.Tests = []rainforest.RunTestDetails{failedTest}

	apiTests := []rainforest.RunTestDetails{
		{
			ID:        failedTest.ID,
			Title:     failedTest.Title,
			CreatedAt: failedTest.CreatedAt,
			UpdatedAt: failedTest.UpdatedAt,
			Result:    failedTest.Result,
			Steps: []rainforest.RunStepDetails{
				{
					Browsers: []rainforest.RunBrowserDetails{
						{
							Name: failedBrowser,
							Feedback: []rainforest.RunFeedback{
								{
									AnswerGiven: "no",
									JobState:    "approved",
									Note:        failedNote,
								},
								{
									AnswerGiven: "yes",
									JobState:    "approved",
									Note:        "This note should not appear",
								},
								{
									AnswerGiven: "no",
									JobState:    "rejected",
									Note:        "This note should not appear either",
								},
							},
						},
					},
				},
			},
		},
	}

	api = newFakeReporterAPI(runDetails.ID, apiTests)

	var out bytes.Buffer
	log.SetOutput(&out)
	schema, err = createJUnitReportSchema(&runDetails, api)
	log.SetOutput(os.Stdout)

	if err != nil {
		t.Errorf("Unexpected error returned by createJunitTestReportSchema: %v", err)
	}

	expectedSchema.Failures = 1
	expectedSchema.TestCases = []jUnitTestReportSchema{
		{
			Name: failedTest.Title,
			Time: 25 * time.Minute.Seconds(),
			Failures: []jUnitTestReportFailure{
				{
					Type:    failedBrowser,
					Message: failedNote,
				},
			},
		},
	}

	if !reflect.DeepEqual(expectedSchema, *schema) {
		t.Error("Incorrect JUnitTestReportSchema returned by createJunitTestReportSchema")
		t.Errorf("Expected: %#v", expectedSchema)
		t.Errorf("Actual: %#v", schema)
	}
}
