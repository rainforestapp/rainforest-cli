package rainforest

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestGetRunDetails(t *testing.T) {
	setup()
	defer cleanup()

	runID := 1337
	reqMethod := "GET"
	runsURL := "/runs/" + strconv.Itoa(runID)

	completeTime, _ := time.Parse(time.RFC3339Nano, "2016-07-13T22:21:31.492Z")
	inProgressTime, _ := time.Parse(time.RFC3339Nano, "2016-07-13T22:06:18.279Z")
	validatingTime, _ := time.Parse(time.RFC3339Nano, "2016-07-13T22:06:12.411Z")
	createdAtTime, _ := time.Parse(time.RFC3339Nano, "2016-07-13T22:06:10.034Z")

	runDetails := RunDetails{
		ID:                 runID,
		Description:        "run description",
		TotalTests:         10,
		TotalFailedTests:   2,
		TotalNoResultTests: 1,
		StateDetails: RunStateDetails{
			Name:         "aborted",
			IsFinalState: true,
		},
		Timestamps: map[string]time.Time{
			"complete":    completeTime,
			"in_progress": inProgressTime,
			"validating":  validatingTime,
			"created_at":  createdAtTime,
		},
	}

	mux.HandleFunc(runsURL, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Unexpected request method in GetRunTestDetails. Expected: %v, Actual: %v", reqMethod, r.Method)
		}

		enc := json.NewEncoder(w)
		enc.Encode(runDetails)
	})

	updatedAt, _ := time.Parse(time.RFC3339Nano, "2016-07-13T22:21:31.492Z")
	createdAt := updatedAt.Add(-10 * time.Minute)
	runTests := []RunTestDetails{
		{
			ID:        999,
			Title:     "Run test title",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Result:    "failed",
		},
	}

	testsURL := "/runs/" + strconv.Itoa(runID) + "/tests"
	mux.HandleFunc(testsURL, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Unexpected request method in GetRunTestDetails. Expected: %v, Actual: %v", reqMethod, r.Method)
		}

		enc := json.NewEncoder(w)
		enc.Encode(runTests)
	})

	out, err := client.GetRunDetails(runID)

	if err != nil {
		t.Errorf("Unexpected error in GetRunTestDetails: %v", err)
	}

	expectedRunDetails := RunDetails{
		ID:                 runDetails.ID,
		Description:        runDetails.Description,
		TotalTests:         runDetails.TotalTests,
		TotalFailedTests:   runDetails.TotalFailedTests,
		TotalNoResultTests: runDetails.TotalNoResultTests,
		StateDetails:       runDetails.StateDetails,
		Timestamps:         runDetails.Timestamps,
		Tests:              runTests,
	}

	if !reflect.DeepEqual(expectedRunDetails, *out) {
		t.Errorf("Unexpected return value from GetRunDetails.\nExpected: %#v\nGot: %#v", expectedRunDetails, *out)
	}
}

func TestGetRunTestDetails(t *testing.T) {
	setup()
	defer cleanup()

	runID := 123
	testID := 456
	reqMethod := "GET"

	updatedAt, _ := time.Parse(time.RFC3339Nano, "2016-07-13T22:21:31.492Z")
	createdAt := updatedAt.Add(-10 * time.Minute)
	runTest := RunTestDetails{
		ID:        runID,
		Title:     "my test title",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Result:    "passed",
		Steps: []RunStepDetails{
			{
				Browsers: []RunBrowserDetails{
					{
						Name: "chrome",
						Feedback: []RunFeedback{
							{
								AnswerGiven: "no",
								JobState:    "approved",
								Note:        "did not work",
							},
							{
								AnswerGiven: "yes",
								JobState:    "rejected",
								Note:        "it worked",
							},
						},
					},
				},
			},
		},
	}

	// TODO: Find the correct pattern for this
	url := "/runs/" + strconv.Itoa(runID) + "/tests/" + strconv.Itoa(testID)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Unexpected request method in GetRunTestDetails. Expected: %v, Actual: %v", reqMethod, r.Method)
		}

		enc := json.NewEncoder(w)
		enc.Encode(runTest)
	})

	out, err := client.GetRunTestDetails(runID, testID)

	if err != nil {
		t.Errorf("Unexpected error in GetRunTestDetails: %v", err)
	}

	if !reflect.DeepEqual(runTest, *out) {
		t.Errorf("Unexpected return value from GetRunTestDetails.\nExpected: %#v\nGot: %#v", runTest, *out)
	}
}
