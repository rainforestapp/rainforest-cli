package rainforest

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestGetRunTestDetails(t *testing.T) {
	setup()
	defer cleanup()

	runID := 123
	testID := 456
	reqMethod := "GET"

	runTest := RunTestDetails{
		ID:        runID,
		Title:     "my test title",
		CreatedAt: "2016-07-13T22:00:00Z",
		UpdatedAt: "2016-07-13T22:10:00Z",
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

	out, _ := client.GetRunTestDetails(runID, testID)

	if !reflect.DeepEqual(runTest, *out) {
		t.Errorf("Unexpected return value from GetRunTestDetails.\nExpected: %#v\nGot: %#v", runTest, *out)
	}
}
