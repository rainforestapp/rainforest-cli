package main

import (
	"testing"

	"github.com/urfave/cli"
)

type reporterFakeContext struct {
	mappings              map[string]interface{}
	args                  cli.Args
	createJunitReportStub func(filename string, runID int)
}

func (c reporterFakeContext) String(s string) string {
	val, ok := c.mappings[s].(string)

	if ok {
		return val
	}
	return ""
}

func (c reporterFakeContext) Args() cli.Args {
	return c.args
}

func TestReporterReportForRun(t *testing.T) {
	r := newReport()
	c := reporterFakeContext{}

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

	r.createJunitReport = func(filename string, runID int, client reporterClient) error {
		if filename != expectedFileName {
			t.Errorf("Unexpected filename given to createJunitReport.\nExpected: %v\nActual: %v", expectedFileName, filename)
		}

		if runID != expectedRunID {
			t.Errorf("Unexpected run ID given to createJunitReport.\nExpected: %v\nActual: %v", expectedRunID, runID)
		}

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
				"junit-file": "myfilename",
			},
			args:     []string{"112233"},
			runID:    112233,
			filename: "myfilename",
		},
		{
			mappings: map[string]interface{}{
				"junit-file": "myfilename",
				"run-id":     "112233",
			},
			args:     []string{},
			runID:    112233,
			filename: "myfilename",
		},
	}

	for _, testCase := range testCases {
		c.mappings = testCase.mappings
		c.args = testCase.args
		expectedRunID = testCase.runID
		expectedFileName = testCase.filename

		r.reportForRun(c)
	}
}
