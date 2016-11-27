package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

func TestReporterReportForRun(t *testing.T) {
	r := newReporter()
	c := newFakeContext(make(map[string]interface{}), cli.Args{})

	err := r.reportForRun(c)

	if err == nil {
		t.Error("No error produced in reporter.reportForRun when Run ID is omitted")
	} else {
		if err.Error() != "No run ID argument found." {
			t.Errorf("Unexpected error in reporter.reportForRun when Run ID is omitted: %v", err.Error())
		}
	}

	var expectedFileName string
	var expectedRunID int

	r.getRunDetails = func(runID int, client *rainforest.Client) (*rainforest.RunDetails, error) {
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

		return os.NewFile(1, "test"), nil
	}

	r.writeJUnitReport = func(*rainforest.RunDetails, *os.File) error {
		return nil
	}

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

		r.reportForRun(c)
	}
}
