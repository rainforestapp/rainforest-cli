package rainforest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetRunDetails(t *testing.T) {
	setup()
	defer cleanup()

	runID := 1337
	reqMethod := "GET"
	runsURL := fmt.Sprintf("/runs/%d", runID)

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

	runDetails.Tests = runTests

	mux.HandleFunc(runsURL, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Unexpected request method in GetRunTestDetails. Expected: %v, Actual: %v", reqMethod, r.Method)
		}

		enc := json.NewEncoder(w)
		enc.Encode(runDetails)
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
		ID:        testID,
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
								Result:      "failed",
								JobState:    "approved",
								FailureNote: "did not work",
							},
							{
								Result:   "passed",
								JobState: "rejected",
							},
						},
					},
				},
			},
		},
	}

	url := fmt.Sprintf("/runs/%d/tests/%d", runID, testID)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Unexpected request method in GetRunTestDetails. Expected: %v, Actual: %v", reqMethod, r.Method)
		}

		var createdAtData, updatedAtData []byte
		var err error
		createdAtData, err = runTest.CreatedAt.MarshalJSON()
		if err != nil {
			t.Fatal(err.Error())
		}

		updatedAtData, err = runTest.UpdatedAt.MarshalJSON()
		if err != nil {
			t.Fatal(err.Error())
		}

		jsonResponse := fmt.Sprintf(`{
"id": %v,
"title": "%v",
"created_at": %v,
"updated_at": %v,
"result": "passed",
"steps":[
	{
		"browsers": [
			{
				"name": "chrome",
				"feedback": [
				{
					"result": "failed",
					"job_state": "approved",
					"note": "did not work"
				},
				{
					"result": "passed",
					"job_state": "rejected"
				}
				]
			}
		]
	}
]
}`,
			runTest.ID,
			runTest.Title,
			string(createdAtData),
			string(updatedAtData),
		)

		_, err = w.Write([]byte(jsonResponse))
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	out, err := client.GetRunTestDetails(runID, testID)

	if err != nil {
		t.Errorf("Unexpected error in GetRunTestDetails: %v", err)
	} else if !reflect.DeepEqual(runTest, *out) {
		t.Errorf("Unexpected return value from GetRunTestDetails.\nExpected: %#v\nGot: %#v", runTest, *out)
	}
}
