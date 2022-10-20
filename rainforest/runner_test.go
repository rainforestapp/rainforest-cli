package rainforest

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestCreateRun(t *testing.T) {
	testCases := []struct {
		runParams RunParams
		wantBody  string
	}{
		{
			runParams: RunParams{Tags: []string{"foo", "bar"}, SiteID: 125},
			wantBody:  `{"tags":["foo","bar"],"site_id":125}`,
		},
		{
			runParams: RunParams{FeatureID: 25, Tags: []string{"baz"}},
			wantBody:  `{"tags":["baz"],"feature_id":25}`,
		},
	}

	for _, tc := range testCases {
		setup()
		defer cleanup()

		const reqMethod = "POST"
		mux.HandleFunc("/runs", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != reqMethod {
				t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			s := strings.TrimSpace(buf.String())
			if s != tc.wantBody {
				t.Errorf("Request body = %v, want %v", s, tc.wantBody)
			}
			fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
		})

		wantStatus := &RunStatus{ID: 123, State: "in_progress"}
		gotStatus, _ := client.CreateRun(tc.runParams)
		if !reflect.DeepEqual(gotStatus, wantStatus) {
			t.Errorf("Response out = %v, want %v", gotStatus, wantStatus)
		}
	}
}

func TestCreateRunFromRerun(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "POST"
	runID := 42
	var completed = false

	runParams := RunParams{
		RunID:    runID,
		Conflict: "cancel",
	}
	wantBody := `{"conflict":"cancel"}`

	mux.HandleFunc("/runs", func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("/runs endpoint hit, expected /runs/:id/rerun_failed")
	})

	wantPath := fmt.Sprintf("/runs/%v/rerun_failed", runID)
	mux.HandleFunc(wantPath, func(w http.ResponseWriter, r *http.Request) {
		completed = true
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		s := strings.TrimSpace(buf.String())
		if s != wantBody {
			t.Errorf("Request body = %v, want %v", s, wantBody)
		}
		fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
	})
	out, err := client.CreateRun(runParams)
	if err != nil {
		t.Error("Error creating run:", err)
	}

	wantStatus := &RunStatus{ID: 123, State: "in_progress"}
	if !reflect.DeepEqual(out, wantStatus) {
		t.Errorf("Response out = %v, want %v", out, wantStatus)
	}

	if !completed {
		t.Error("Run API endpoint not hit")
	}

	runID = 117
	// Check that you can't combine filters
	badParams := []RunParams{
		{
			RunID: runID,
			Tags:  []string{"foo"},
		},
		{
			RunID:    runID,
			Browsers: []string{"chrome"},
		},
		{
			RunID: runID,
			Tests: "all",
		},
		{
			RunID:   runID,
			RFMLIDs: []string{"rmfl1"},
		},
		{
			RunID:           runID,
			ExecutionMethod: "automation",
		},
		{
			RunID:       runID,
			Description: "My fancy rerun",
		},
		{
			RunID:   runID,
			Release: "somerandomsha",
		},
		{
			RunID:         runID,
			EnvironmentID: 35,
		},
		{
			RunID:  runID,
			SiteID: 35,
		},
		{
			RunID:     runID,
			FeatureID: 35,
		},
		{
			RunID:         runID,
			SmartFolderID: 123,
		},
		{
			RunID:      runID,
			RunGroupID: 123,
		},
		{
			RunID:    runID,
			BranchID: 123,
		},
	}
	mux.HandleFunc(fmt.Sprintf("/runs/%v/rerun_failed", runID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
	})
	for _, params := range badParams {
		_, err := client.CreateRun(params)
		if err == nil {
			t.Errorf("Expected error for params %v but there was no error", params)
		}
	}
}

func TestCreateRunFromRunGroup(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "POST"
	runGroupID := 25
	var completed = false

	runParams := RunParams{
		EnvironmentID: 23,
		Conflict:      "cancel",
		RunGroupID:    runGroupID,
	}
	wantBody := `{"conflict":"cancel","environment_id":23}`

	mux.HandleFunc("/runs", func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("/runs endpoint hit, expected /run_groups/:id/runs")
	})

	wantPath := fmt.Sprintf("/run_groups/%v/runs", runGroupID)
	mux.HandleFunc(wantPath, func(w http.ResponseWriter, r *http.Request) {
		completed = true
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		s := strings.TrimSpace(buf.String())
		if s != wantBody {
			t.Errorf("Request body = %v, want %v", s, wantBody)
		}
		fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
	})

	out, err := client.CreateRun(runParams)
	if err != nil {
		t.Error("Error creating run:", err)
	}

	wantStatus := &RunStatus{ID: 123, State: "in_progress"}
	if !reflect.DeepEqual(out, wantStatus) {
		t.Errorf("Response out = %v, want %v", out, wantStatus)
	}

	if !completed {
		t.Error("Run API endpoint not hit")
	}

	runGroupID = 42
	// Check that you can't combine filters
	badParams := []RunParams{
		{
			RunGroupID: runGroupID,
			Tags:       []string{"foo"},
		},
		{
			RunGroupID: runGroupID,
			SiteID:     17,
		},
		{
			RunGroupID: runGroupID,
			Tests:      "all",
		},
		{
			RunGroupID:    runGroupID,
			SmartFolderID: 123,
		},
		{
			RunGroupID: runGroupID,
			Browsers:   []string{"chrome"},
		},
		{
			RunGroupID: runGroupID,
			FeatureID:  35,
		},
	}
	mux.HandleFunc(fmt.Sprintf("/run_groups/%v/runs", runGroupID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 123, "state":"in_progress"}`)
	})
	for _, params := range badParams {
		_, err := client.CreateRun(params)
		if err == nil {
			t.Errorf("Expected error for params %v but there was no error", params)
		}
	}
}

func TestCheckRunStatus(t *testing.T) {
	setup()
	defer cleanup()

	const reqMethod = "GET"

	mux.HandleFunc("/runs/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != reqMethod {
			t.Errorf("Request method = %v, want %v", r.Method, reqMethod)
		}

		fmt.Fprint(w, `{"id": 123, "state":"in_progress", "result":"passed"}`)
	})

	out, _ := client.CheckRunStatus(123)

	want := &RunStatus{ID: 123, State: "in_progress", Result: "passed"}

	if !reflect.DeepEqual(out, want) {
		t.Errorf("Response out = %v, want %v", out, want)
	}
}
